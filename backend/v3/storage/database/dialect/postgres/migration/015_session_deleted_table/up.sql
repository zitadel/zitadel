CREATE TABLE zitadel.sessions_deleted (
    instance_id TEXT NOT NULL
    , id TEXT NOT NULL CHECK (id <> '')
    , token_id TEXT
    , user_agent_id TEXT
    , expiration TIMESTAMPTZ
    , user_id TEXT
    , creator_id TEXT
    , deleted_at TIMESTAMPTZ DEFAULT NOW() NOT NULL

    , PRIMARY KEY (instance_id, id)
    , FOREIGN KEY (instance_id) REFERENCES zitadel.instances(id)
    , FOREIGN KEY (instance_id, user_id) REFERENCES zitadel.users(instance_id, id) ON DELETE CASCADE
    , FOREIGN KEY (instance_id, user_agent_id) REFERENCES zitadel.session_user_agents(instance_id, fingerprint_id) ON DELETE SET NULL (user_agent_id)
);

CREATE OR REPLACE FUNCTION zitadel.move_to_deleted_sessions() RETURNS trigger AS $$
BEGIN
    INSERT INTO zitadel.sessions_deleted (instance_id, id, token_id, user_agent_id, expiration, user_id, creator_id, deleted_at)
    VALUES (OLD.instance_id, OLD.id, OLD.token_id, OLD.user_agent_id, OLD.expiration, OLD.user_id, OLD.creator_id, now());
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_move_to_deleted_sessions
AFTER DELETE ON zitadel.sessions
FOR EACH ROW
EXECUTE FUNCTION zitadel.move_to_deleted_sessions();

CREATE OR REPLACE FUNCTION zitadel.throw_not_permitted() returns boolean AS $$
BEGIN
    RAISE EXCEPTION 'Permission denied: User does not have permission'
        USING ERRCODE = 'insufficient_privilege';
END;
$$ LANGUAGE plpgsql;