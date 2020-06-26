BEGIN;

ALTER TABLE management.failed_event RENAME TO management.failed_events;
ALTER TABLE auth.failed_event RENAME TO auth.failed_events;
ALTER TABLE notification.failed_event RENAME TO notification.failed_events;
ALTER TABLE authz.failed_event RENAME TO authz.failed_events;
ALTER TABLE admin.failed_event RENAME TO admin.failed_events;

COMMIT;
