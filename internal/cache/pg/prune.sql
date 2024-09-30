delete from cache.objects o
where o.cache_name = $1
	and (
		case when $2 > '0s'
			then created_at < now()-$2 -- max age
			else false
		end
		or case when $3 > '0s'
			then last_used_at < now()-$3 -- last use
			else false
		end
		or o.id not in (
			select object_id
			from cache.string_keys
			where cache_name = $1
		)
	)
;
