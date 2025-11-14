update cache.objects
set last_used_at = now()
where cache_name = $1
	and	(
		select object_id
		from cache.string_keys k
		where cache_name = $1
			and index_id = $2
			and index_key = $3
		) = id
	and case when $4::interval > '0s'
		then created_at > now()-$4::interval -- max age
		else true
	end
	and case when $5::interval > '0s'
		then last_used_at > now()-$5::interval -- last use
		else true
	end
returning payload;
