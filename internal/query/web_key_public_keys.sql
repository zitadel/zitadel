select public_key
from projections.web_keys1
where instance_id = $1;
