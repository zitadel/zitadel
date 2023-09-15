DO $$
BEGIN

    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_schema = 'eventstore' AND table_name = 'events' AND column_name = 'position') THEN
        ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS "position" xid8;
        ALTER TABLE eventstore.events ALTER COLUMN "position" SET NOT NULL,
                                      ALTER COLUMN "position" TYPE xid8 USING pg_current_xact_id();
    END IF;

    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_schema = 'eventstore' AND table_name = 'events' AND column_name = 'in_tx_order') THEN
            ALTER TABLE eventstore.events ADD COLUMN IF NOT EXISTS in_tx_order INTEGER;
            ALTER TABLE eventstore.events ALTER COLUMN in_tx_order SET NOT NULL,
                                          ALTER COLUMN in_tx_order TYPE INTEGER USING event_sequence;
    END IF;

END $$;