UPDATE projections.current_states SET last_updated = statement_timestamp() WHERE projection_name = $1 AND instance_id = $2;
    