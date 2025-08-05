select key_id, creation_date, change_date, sequence, state, config, config_type
from projections.web_keys1
where instance_id = $1
order by creation_date asc;
