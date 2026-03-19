DROP TRIGGER IF EXISTS trg_move_to_archived_sessions ON zitadel.sessions;
DROP FUNCTION IF EXISTS zitadel.move_to_archived_sessions() CASCADE;
DROP TABLE IF EXISTS zitadel.archived_sessions CASCADE;
DROP FUNCTION IF EXISTS zitadel.throw_not_permitted();