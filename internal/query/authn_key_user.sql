select u.id as user_id, u.resource_owner, u.username, m.access_token_type, k.public_key
from projections.authn_keys2 k
join projections.users14 u
	on k.instance_id = u.instance_id
	and k.identifier = u.id
join projections.users14_machines m
	on u.instance_id = m.instance_id
	and u.id = m.user_id
where k.instance_id = $1
	and k.id = $2
	and u.id = $3
    and k.expiration > current_timestamp;
