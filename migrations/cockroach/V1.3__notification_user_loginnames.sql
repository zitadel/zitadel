BEGIN;

ALTER TABLE notification.notify_users ADD COLUMN login_names TEXT ARRAY;
ALTER TABLE notification.notify_users ADD COLUMN preferred_login_name TEXT;

COMMIT;