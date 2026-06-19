-- +goose Up
-- +goose StatementBegin
create type document_status as enum (
    'UPLOADED', 'IN_REVIEW', 'APPROVED', 'REJECTED', 'ARCHIVED'
);
create type process_status as enum (
    'PENDING', 'IN_PROGRESS', 'APPROVED', 'REJECTED', 'CANCELLED'
);
create type step_status as enum (
    'PENDING', 'ACTIVE', 'APPROVED', 'REJECTED', 'SKIPPED'
);
create type step_policy as enum (
    'ANY', 'ALL', 'QUORUM'
);
create type decision_type as enum (
    'APPROVE', 'REJECT'
);
create type signature_algorithm as enum (
    'ED25519', 'ECDSA_P256_SHA256', 'RSA_PSS_SHA256'
);

create table folders (
    id         uuid primary key default gen_random_uuid(),
    name       text not null,
    parent_id  uuid references folders (id),
    created_by uuid not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
create index idx_folders_parent_id on folders (parent_id);

create table documents (
    id           uuid primary key default gen_random_uuid(),
    name         text not null,
    content_type text not null,
    size_bytes   bigint not null,
    sha256       text not null,
    folder_id    uuid references folders (id),
    status       document_status not null default 'UPLOADED',
    storage_key  text not null,
    created_by   uuid not null,
    created_at   timestamptz not null default now(),
    updated_at   timestamptz not null default now()
);
create index idx_documents_folder_id on documents (folder_id);
create index idx_documents_status on documents (status);

create table approval_processes (
    id                 uuid primary key default gen_random_uuid(),
    document_id        uuid not null references documents (id),
    status             process_status not null default 'PENDING',
    current_step_index int,
    created_by         uuid not null,
    created_at         timestamptz not null default now(),
    updated_at         timestamptz not null default now()
);
create index idx_approval_processes_document_id on approval_processes (document_id);
create index idx_approval_processes_status on approval_processes (status);

create table process_steps (
    id                uuid primary key default gen_random_uuid(),
    process_id        uuid not null references approval_processes (id) on delete cascade,
    index             int not null,
    name              text not null,
    policy            step_policy not null,
    required_approvals int,
    status            step_status not null default 'PENDING',
    approver_user_ids uuid[] not null default '{}',
    approver_roles    text[] not null default '{}',
    created_at        timestamptz not null default now(),
    updated_at        timestamptz not null default now(),
    unique (process_id, index)
);
create index idx_process_steps_process_id on process_steps (process_id);

create table decisions (
    id         uuid primary key default gen_random_uuid(),
    step_id    uuid not null references process_steps (id) on delete cascade,
    actor_id   uuid not null,
    type       decision_type not null,
    comment    text,
    created_at timestamptz not null default now()
);
create index idx_decisions_step_id on decisions (step_id);

create table signatures (
    id                    uuid primary key default gen_random_uuid(),
    decision_id           uuid not null unique references decisions (id) on delete cascade,
    algorithm             signature_algorithm not null,
    value                 text not null,
    key_id                text not null,
    certificate           text,
    signed_payload_sha256 text not null,
    signed_at             timestamptz not null
);

create table review_comments (
    id         uuid primary key default gen_random_uuid(),
    process_id uuid not null references approval_processes (id) on delete cascade,
    step_id    uuid references process_steps (id),
    author_id  uuid not null,
    body       text not null,
    created_at timestamptz not null default now()
);
create index idx_review_comments_process_id on review_comments (process_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists review_comments;
drop table if exists signatures;
drop table if exists decisions;
drop table if exists process_steps;
drop table if exists approval_processes;
drop table if exists documents;
drop table if exists folders;

drop type if exists signature_algorithm;
drop type if exists decision_type;
drop type if exists step_policy;
drop type if exists step_status;
drop type if exists process_status;
drop type if exists document_status;
-- +goose StatementEnd
