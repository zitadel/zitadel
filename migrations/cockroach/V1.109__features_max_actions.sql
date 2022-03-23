alter table zitadel.projections.features
    ADD COLUMN  actions_allowed INT2 default 0,
    ADD COLUMN  max_actions INT8 default 0;

update zitadel.projections.features set actions_allowed = 2 where actions = true;
