UPDATE eventstore.events e SET (creation_date, "position") = (we.next_cd, we.next_cd::DECIMAL) FROM wrong_events we WHERE e.event_sequence = we.event_sequence AND e.instance_id = we.instance_id;
