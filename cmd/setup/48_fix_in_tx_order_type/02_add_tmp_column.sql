ALTER TABLE IF EXISTS eventstore.events2
    ADD COLUMN IF NOT EXISTS in_tx_order_tmp INTEGER;

CREATE OR REPLACE FUNCTION eventstore.sync_in_tx_order()
RETURNS trigger
AS $$
    BEGIN
        NEW.in_tx_order_tmp := NEW.in_tx_order::INTEGER;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER sync_in_tx_order
BEFORE INSERT ON eventstore.events2
FOR EACH ROW EXECUTE FUNCTION eventstore.sync_in_tx_order();
