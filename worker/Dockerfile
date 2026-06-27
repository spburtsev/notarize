FROM python:3.14-slim

COPY --from=ghcr.io/astral-sh/uv:0.10 /uv /uvx /bin/

WORKDIR /app
ENV PYTHONUNBUFFERED=1 \
    UV_COMPILE_BYTECODE=1 \
    UV_LINK_MODE=copy \
    UV_PYTHON=python3.14 \
    UV_PYTHON_DOWNLOADS=never \
    HF_HOME=/root/.cache/huggingface \
    PATH=/app/.venv/bin:$PATH

COPY pyproject.toml uv.lock ./
RUN --mount=type=cache,target=/root/.cache/uv \
    uv sync --frozen --no-dev

COPY worker.py ./

CMD ["python", "worker.py"]
