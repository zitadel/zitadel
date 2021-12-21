truncate zitadel.projections.users cascade;
update zitadel.projections.current_sequences set current_sequence = 0 where projection_name = 'zitadel.projections.users';
