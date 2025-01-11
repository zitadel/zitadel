DO
$$BEGIN
    CREATE TRIGGER write_event_queue
    AFTER INSERT ON eventstore.events2
    FOR EACH ROW
    EXECUTE FUNCTION subscriptions.queue_events();
EXCEPTION
   WHEN duplicate_object THEN
      NULL;
END;
$$;
