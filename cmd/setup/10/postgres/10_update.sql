UPDATE eventstore.events e SET creation_date = we.next_cd FROM wrong_events we WHERE e.event_sequence = we.event_sequence and e.instance_id = we.instance_id;
