CREATE OR REPLACE FUNCTION projections.current_states_keep_in_tx_order_opt_in()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  PERFORM set_config('zitadel.keep_in_tx_order', 'true', true);
  RETURN NULL;
END;
$$;

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

-- UPDATE OF in_tx_order is the opt-in; see handler/v2/state_set.sql.
DROP TRIGGER IF EXISTS current_states_keep_in_tx_order_opt_in ON projections.current_states;
CREATE TRIGGER current_states_keep_in_tx_order_opt_in
  BEFORE UPDATE OF in_tx_order ON projections.current_states
  FOR EACH STATEMENT
  EXECUTE FUNCTION projections.current_states_keep_in_tx_order_opt_in();

DROP TRIGGER IF EXISTS current_states_keep_in_tx_order ON projections.current_states;
CREATE TRIGGER current_states_keep_in_tx_order
  BEFORE UPDATE ON projections.current_states
  FOR EACH ROW
  EXECUTE FUNCTION projections.current_states_keep_in_tx_order();
