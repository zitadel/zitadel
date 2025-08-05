select private_key
from projections.web_keys1
where instance_id = $1
and state = $2
limit 1;
