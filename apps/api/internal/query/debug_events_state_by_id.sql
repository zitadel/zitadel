select id, creation_date, change_date, resource_owner, sequence, blob
from projections.debug_events
where instance_id = $1
and id = $2
limit 1;
