create schema if not exists cache;

create unlogged table if not exists cache.objects (
	cache_name varchar not null,
    id uuid not null default gen_random_uuid(),
    created_at timestamptz not null default now(),
	last_used_at timestamptz not null default now(),
    payload jsonb not null,

	primary key(cache_name, id)
)
partition by list (cache_name);

create unlogged table if not exists cache.string_keys(
    cache_name varchar not null check (cache_name <> ''),
    index_id integer not null check (index_id > 0),
    index_key varchar not null check (index_key <> ''),
    object_id uuid not null,

    primary key (cache_name, index_id, index_key),
    constraint fk_object
        foreign key(cache_name, object_id)
        references cache.objects(cache_name, id)
        on delete cascade
)
partition by list (cache_name);

create index if not exists string_keys_object_id_idx
    on cache.string_keys (cache_name, object_id); -- for delete cascade
