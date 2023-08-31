UPDATE projections.current_states SET last_updated = now() WHERE projection_name = $1 AND instance_id = $2;
    