with object as (
	insert into cache.objects (cache_name, payload)
	values ($1, $3)
	returning id
)
insert into cache.string_keys (
	cache_name,
	index_id,
	index_key,
	object_id
)
select $1, keys.index_id, keys.index_key, id as object_id
from object, jsonb_to_recordset($2) keys (
	index_id bigint,
	index_key text
)
on conflict (cache_name, index_id, index_key) do
	update set object_id = EXCLUDED.object_id
;
