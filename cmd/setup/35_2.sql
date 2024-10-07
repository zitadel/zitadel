DROP INDEX IF EXISTS eventstore.es_wm;

ALTER INDEX eventstore.es_wm_temp RENAME TO es_wm;
