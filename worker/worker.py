from __future__ import annotations

import io
import multiprocessing as mp
import os
import signal
import sys
import tempfile
import traceback
import uuid
from multiprocessing.synchronize import Event
from typing import TYPE_CHECKING

import psycopg
from minio import Minio
from psycopg.rows import DictRow, dict_row

if TYPE_CHECKING:
    from types import FrameType
    from docling.document_converter import DocumentConverter


def env(name: str, default: str | None = None) -> str:
    val = os.environ.get(name, default)
    if val is None:
        sys.exit(f"missing required env var {name}")
    return val


DATABASE_URL = env("DATABASE_URL")
S3_ENDPOINT = env("S3_ENDPOINT")
S3_ACCESS_KEY = env("S3_ACCESS_KEY")
S3_SECRET_KEY = env("S3_SECRET_KEY")
S3_BUCKET = env("S3_BUCKET")
S3_USE_SSL = env("S3_USE_SSL", "false").lower() == "true"
POLL_INTERVAL = float(env("POLL_INTERVAL_SECONDS", "2"))
WORKER_CONCURRENCY = int(env("WORKER_CONCURRENCY", "2"))


def claim_job(conn: psycopg.Connection) -> DictRow | None:
    with conn.transaction():
        with conn.cursor(row_factory=dict_row) as cur:
            cur.execute(
                """
                SELECT id, document_id FROM document_processing_jobs
                WHERE status = 'QUEUED'
                ORDER BY created_at
                FOR UPDATE SKIP LOCKED
                LIMIT 1
                """
            )
            job = cur.fetchone()
            if job is None:
                return None
            cur.execute(
                """
                UPDATE document_processing_jobs
                SET status='PROCESSING', started_at=now(), attempts=attempts+1, updated_at=now()
                WHERE id=%s
                """,
                (job["id"],),
            )
    return job


def fetch_document(conn: psycopg.Connection, document_id: uuid.UUID) -> DictRow:
    with conn.cursor(row_factory=dict_row) as cur:
        cur.execute(
            "SELECT name, content_type, storage_key FROM documents WHERE id=%s",
            (document_id,),
        )
        fetched = cur.fetchone()
        assert fetched is not None, f"document {document_id} not found"
        return fetched


def extract_markdown(s3: Minio, converter: DocumentConverter, doc: DictRow) -> str:
    suffix = os.path.splitext(doc["name"])[1]
    with tempfile.NamedTemporaryFile(suffix=suffix) as tmp:
        s3.fget_object(S3_BUCKET, doc["storage_key"], tmp.name)
        result = converter.convert(tmp.name)
        return result.document.export_to_markdown()


def finish(conn: psycopg.Connection, job_id: uuid.UUID, status: str, *, error: str | None = None):
    with conn.transaction():
        with conn.cursor() as cur:
            cur.execute(
                """
                UPDATE document_processing_jobs
                SET status=%s, error=%s, finished_at=now(), updated_at=now()
                WHERE id=%s
                """,
                (status, error, job_id),
            )


def process(conn: psycopg.Connection, s3: Minio, converter: DocumentConverter, job: DictRow):
    document_id = job["document_id"]
    name = mp.current_process().name
    try:
        doc = fetch_document(conn, document_id)
        md = extract_markdown(s3, converter, doc)
        out_key = f"derived/{document_id}.md"
        data = md.encode("utf-8")
        s3.put_object(S3_BUCKET, out_key, io.BytesIO(data), len(data), content_type="text/markdown")
        finish(conn, job["id"], "DONE")
        print(f"[{name}] job {job['id']} done -> {out_key} ({len(data)} bytes)", flush=True)
    except Exception:
        err = traceback.format_exc()
        print(f"[{name}] job {job['id']} failed:\n{err}", flush=True)
        finish(conn, job["id"], "FAILED", error=err)


def worker_loop(stop_event: Event):
    name = mp.current_process().name
    s3 = Minio(S3_ENDPOINT, access_key=S3_ACCESS_KEY, secret_key=S3_SECRET_KEY, secure=S3_USE_SSL)
    converter = None  # built on first job (imports torch post-fork, in this process)
    print(f"[{name}] started", flush=True)

    while not stop_event.is_set():
        try:
            with psycopg.connect(DATABASE_URL, autocommit=True) as conn:
                while not stop_event.is_set():
                    job = claim_job(conn)
                    if job is None:
                        stop_event.wait(POLL_INTERVAL)
                        continue
                    if converter is None:
                        from docling.datamodel.base_models import InputFormat
                        from docling.datamodel.pipeline_options import PdfPipelineOptions
                        from docling.document_converter import (
                            DocumentConverter,
                            PdfFormatOption,
                        )

                        pdf_opts = PdfPipelineOptions()
                        pdf_opts.do_ocr = False
                        converter = DocumentConverter(
                            format_options={InputFormat.PDF: PdfFormatOption(pipeline_options=pdf_opts)}
                        )
                    process(conn, s3, converter, job)
        except Exception as e:
            print(f"[{name}] loop error, retrying in {POLL_INTERVAL}s: {e}", flush=True)
            stop_event.wait(POLL_INTERVAL)

    print(f"[{name}] stopped", flush=True)


def main():
    n = max(1, WORKER_CONCURRENCY)
    stop_event = mp.Event()

    def request_stop(signum: int, frame: FrameType | None) -> None:
        stop_event.set()

    signal.signal(signal.SIGTERM, request_stop)
    signal.signal(signal.SIGINT, request_stop)

    print(f"starting {n} worker processes", flush=True)
    procs = [
        mp.Process(target=worker_loop, args=(stop_event,), name=f"worker-{i}") for i in range(n)
    ]
    for p in procs:
        p.start()
    for p in procs:
        p.join()


if __name__ == "__main__":
    main()
