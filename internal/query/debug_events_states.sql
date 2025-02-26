select id, creation_date, change_date, resource_owner, sequence, blob
from projections.debug_events
where instance_id = $1
order by creation_date asc;
