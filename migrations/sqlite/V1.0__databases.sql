ATTACH DATABASE '${GOPATH}/src/github.com/caos/zitadel/.local/management.db' AS 'management';
ATTACH DATABASE '${GOPATH}/src/github.com/caos/zitadel/.local/auth.db' AS 'auth';
ATTACH DATABASE '${GOPATH}/src/github.com/caos/zitadel/.local/notification.db' AS 'notification';
ATTACH DATABASE '${GOPATH}/src/github.com/caos/zitadel/.local/adminapi.db' AS 'adminapi';
ATTACH DATABASE '${GOPATH}/src/github.com/caos/zitadel/.local/authz.db' AS 'authz';
ATTACH DATABASE '${GOPATH}/src/github.com/caos/zitadel/.local/eventstore.db' AS 'eventstore';



-- CREATE USER eventstore;
-- GRANT SELECT, INSERT ON DATABASE eventstore TO eventstore;

-- CREATE USER management;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON DATABASE management TO management;
-- GRANT SELECT, INSERT ON DATABASE eventstore TO management;

-- CREATE USER adminapi;
-- GRANT SELECT, INSERT, UPDATE, DELETE, DROP ON DATABASE adminapi TO adminapi;
-- GRANT SELECT, INSERT ON DATABASE eventstore TO adminapi;
-- GRANT SELECT, INSERT, UPDATE, DROP, DELETE  ON DATABASE auth TO adminapi;
-- GRANT SELECT, INSERT, UPDATE, DROP, DELETE ON DATABASE authz TO adminapi;
-- GRANT SELECT, INSERT, UPDATE, DROP, DELETE ON DATABASE management TO adminapi;
-- GRANT SELECT, INSERT, UPDATE, DROP, DELETE ON DATABASE notification TO adminapi;

-- CREATE USER auth;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON DATABASE auth TO auth;
-- GRANT SELECT, INSERT ON DATABASE eventstore TO auth;

-- CREATE USER notification;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON DATABASE notification TO notification;
-- GRANT SELECT, INSERT ON DATABASE eventstore TO notification;

-- CREATE USER authz;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON DATABASE authz TO authz;
-- GRANT SELECT, INSERT ON DATABASE eventstore TO authz;
-- GRANT SELECT, INSERT, UPDATE ON DATABASE auth TO authz;
