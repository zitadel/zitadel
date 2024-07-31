select public_key
from projections.web_keys
where instance_id = $1;
