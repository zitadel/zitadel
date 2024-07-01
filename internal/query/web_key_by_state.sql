select private_key
from projections.web_keys
where instance_id = $1
and state = $2
limit 1;
