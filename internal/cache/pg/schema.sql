create schema if not exists cache;

create unlogged table if not exists cache.objects (
    id bigint primary key generated always as identity,
    created_at timestamptz not null default now(),
    payload jsonb not null
);

create unlogged table if not exists cache.string_keys(
    cache_name varchar not null check (cache_name <> ''),
    index_id integer not null check (index_id <> 0),
    index_key varchar not null check (index_key <> ''),
    object_id bigint not null,

    primary key (cache_name, index_id, index_key),
    constraint fk_object
        foreign key(object_id)
        references cache.objects(id)
        on delete cascade
);

create index if not exists object_id_idx
    on cache.string_keys (object_id); -- for delete cascade
