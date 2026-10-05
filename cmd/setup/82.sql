ALTER TABLE projections.current_states ADD COLUMN IF NOT EXISTS in_tx_order INTEGER;

UPDATE projections.current_states cs
SET in_tx_order = e.in_tx_order
FROM eventstore.events2 e
WHERE cs.instance_id = e.instance_id
  AND cs.aggregate_id = e.aggregate_id
  AND cs.aggregate_type = e.aggregate_type
  AND cs."sequence" = e.sequence
  AND e.in_tx_order <> 0;

CREATE OR REPLACE FUNCTION projections.current_states_keep_in_tx_order()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  IF current_setting('zitadel.keep_in_tx_order', true) IS DISTINCT FROM 'true' THEN
    NEW.in_tx_order := NULL;
  END IF;
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS current_states_keep_in_tx_order ON projections.current_states;
CREATE TRIGGER current_states_keep_in_tx_order
  BEFORE INSERT OR UPDATE ON projections.current_states
  FOR EACH ROW
  EXECUTE FUNCTION projections.current_states_keep_in_tx_order();
