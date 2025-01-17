DO
$$BEGIN
    CREATE TRIGGER reduce_queued_events
    BEFORE INSERT ON subscriptions.queue
    FOR EACH ROW
    WHEN (NEW.allow_reduce)
    EXECUTE FUNCTION subscriptions.reduce_queued_events();
EXCEPTION
   WHEN duplicate_object THEN
      NULL;
END;
$$;