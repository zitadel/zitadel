alter table zitadel.projections.features
    alter column is_default set default false,
    alter column login_policy_factors set default false,
    alter column login_policy_idp set default false,
    alter column login_policy_passwordless set default false,
    alter column login_policy_registration set default false,
    alter column login_policy_username_login set default false,
    alter column login_policy_password_reset set default false,
    alter column password_complexity_policy set default false,
    alter column label_policy_private_label set default false,
    alter column label_policy_watermark set default false,
    alter column custom_domain set default false,
    alter column privacy_policy set default false,
    alter column metadata_user set default false,
    alter column custom_text_message set default false,
    alter column custom_text_login set default false,
    alter column lockout_policy set default false,
    alter column actions set default false;

update zitadel.projections.features set is_default = default where is_default is null;
update zitadel.projections.features set state = 0 where state is null;
update zitadel.projections.features set audit_log_retention = 0 where audit_log_retention is null;
update zitadel.projections.features set login_policy_factors = default where login_policy_factors is null;
update zitadel.projections.features set login_policy_idp = default where login_policy_idp is null;
update zitadel.projections.features set login_policy_passwordless = default where login_policy_passwordless is null;
update zitadel.projections.features set login_policy_registration = default where login_policy_registration is null;
update zitadel.projections.features set login_policy_username_login = default where login_policy_username_login is null;
update zitadel.projections.features set login_policy_password_reset = default where login_policy_password_reset is null;
update zitadel.projections.features set password_complexity_policy = default where password_complexity_policy is null;
update zitadel.projections.features set label_policy_private_label = default where label_policy_private_label is null;
update zitadel.projections.features set label_policy_watermark = default where label_policy_watermark is null;
update zitadel.projections.features set custom_domain = default where custom_domain is null;
update zitadel.projections.features set privacy_policy = default where privacy_policy is null;
update zitadel.projections.features set metadata_user = default where metadata_user is null;
update zitadel.projections.features set custom_text_message = default where custom_text_message is null;
update zitadel.projections.features set custom_text_login = default where custom_text_login is null;
update zitadel.projections.features set lockout_policy = default where lockout_policy is null;
update zitadel.projections.features set actions = default where actions is null;

alter table zitadel.projections.features
    alter column change_date set not null,
    alter column sequence set not null,
    alter column is_default set not null,
    alter column state set not null,
    alter column audit_log_retention set not null,
    alter column login_policy_factors set not null,
    alter column login_policy_idp set not null,
    alter column login_policy_passwordless set not null,
    alter column login_policy_registration set not null,
    alter column login_policy_username_login set not null,
    alter column login_policy_password_reset set not null,
    alter column password_complexity_policy set not null,
    alter column label_policy_private_label set not null,
    alter column label_policy_watermark set not null,
    alter column custom_domain set not null,
    alter column privacy_policy set not null,
    alter column metadata_user set not null,
    alter column custom_text_message set not null,
    alter column custom_text_login set not null,
    alter column lockout_policy set not null,
    alter column actions set not null;
