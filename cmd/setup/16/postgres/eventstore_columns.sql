DO $$
BEGIN

    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_schema = 'eventstore' AND table_name = 'events' AND column_name = 'position') THEN
        ALTER TABLE eventstore.events ALTER COLUMN "position" SET NOT NULL,
                                      ALTER COLUMN "position" TYPE DECIMAL USING COALESCE("position", EXTRACT(EPOCH FROM creation_date));
    END IF;

    IF NOT EXISTS(SELECT 1 FROM information_schema.columns WHERE table_schema = 'eventstore' AND table_name = 'events' AND column_name = 'in_tx_order') THEN
            ALTER TABLE eventstore.events ALTER COLUMN in_tx_order SET NOT NULL,
                                          ALTER COLUMN in_tx_order TYPE INTEGER USING COALESCE(in_tx_order, event_sequence);
    END IF;

END $$;