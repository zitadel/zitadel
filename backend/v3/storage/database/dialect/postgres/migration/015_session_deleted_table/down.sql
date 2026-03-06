DROP TRIGGER IF EXISTS trg_move_to_deleted_sessions ON zitadel.sessions;
DROP FUNCTION IF EXISTS zitadel.move_to_deleted_sessions() CASCADE;
DROP TABLE IF EXISTS zitadel.sessions_deleted CASCADE;
DROP FUNCTION IF EXISTS zitadel.throw_not_permitted();