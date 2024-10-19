delete from cache.objects o
where o.cache_name = $1
	and (
		case when $2::interval > '0s'
			then created_at < now()-$2::interval -- max age
			else false
		end
		or case when $3::interval > '0s'
			then last_used_at < now()-$3::interval -- last use
			else false
		end
		or o.id not in (
			select object_id
			from cache.string_keys
			where cache_name = $1
		)
	)
;
