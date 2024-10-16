delete from cache.string_keys k
where k.cache_name = $1
	and k.index_id = $2
	and k.index_key = any($3)
;
