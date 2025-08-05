delete from cache.objects o
using cache.string_keys k
where k.cache_name = $1
	and k.index_id = $2
	and k.index_key = any($3)
	and o.cache_name = k.cache_name
	and o.id = k.object_id
;

