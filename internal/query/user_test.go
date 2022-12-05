package query

var (
	loginNamesQuery = `SELECT login_names.user_id, ARRAY_AGG(login_names.login_name)::TEXT[] AS loginnames, login_names.instance_id, login_names.user_owner_removed, login_names.policy_owner_removed, login_names.domain_owner_removed` +
		` FROM projections.login_names2 AS login_names` +
		` GROUP BY login_names.user_id, login_names.instance_id, login_names.user_owner_removed, login_names.policy_owner_removed, login_names.domain_owner_removed`
	preferredLoginNameQuery = `SELECT preferred_login_name.user_id, preferred_login_name.login_name, preferred_login_name.instance_id, preferred_login_name.user_owner_removed, preferred_login_name.policy_owner_removed, preferred_login_name.domain_owner_removed` +
		` FROM projections.login_names2 AS preferred_login_name` +
		` WHERE  preferred_login_name.is_primary = $1`
	userQuery = `SELECT projections.users6.id,` +
		` projections.users6.creation_date,` +
		` projections.users6.change_date,` +
		` projections.users6.resource_owner,` +
		` projections.users6.sequence,` +
		` projections.users6.state,` +
		` projections.users6.type,` +
		` projections.users6.username,` +
		` login_names.loginnames,` +
		` preferred_login_name.login_name,` +
		` projections.users6_humans.user_id,` +
		` projections.users6_humans.first_name,` +
		` projections.users6_humans.last_name,` +
		` projections.users6_humans.nick_name,` +
		` projections.users6_humans.display_name,` +
		` projections.users6_humans.preferred_language,` +
		` projections.users6_humans.gender,` +
		` projections.users6_humans.avatar_key,` +
		` projections.users6_humans.email,` +
		` projections.users6_humans.is_email_verified,` +
		` projections.users6_humans.phone,` +
		` projections.users6_humans.is_phone_verified,` +
		` projections.users6_machines.user_id,` +
		` projections.users6_machines.name,` +
		` projections.users6_machines.description,` +
		` COUNT(*) OVER ()` +
		` FROM projections.users6` +
		` LEFT JOIN projections.users6_humans ON projections.users6.id = projections.users6_humans.user_id AND projections.users6.instance_id = projections.users6_humans.instance_id` +
		` LEFT JOIN projections.users6_machines ON projections.users6.id = projections.users6_machines.user_id AND projections.users6.instance_id = projections.users6_machines.instance_id` +
		` LEFT JOIN` +
		` (` + loginNamesQuery + `) AS login_names` +
		` ON login_names.user_id = projections.users6.id AND login_names.instance_id = projections.users6.instance_id` +
		` LEFT JOIN` +
		` (` + preferredLoginNameQuery + `) AS preferred_login_name` +
		` ON preferred_login_name.user_id = projections.users6.id AND preferred_login_name.instance_id = projections.users6.instance_id`
	userCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"state",
		"type",
		"username",
		"loginnames",
		"login_name",
		//human
		"user_id",
		"first_name",
		"last_name",
		"nick_name",
		"display_name",
		"preferred_language",
		"gender",
		"avatar_key",
		"email",
		"is_email_verified",
		"phone",
		"is_phone_verified",
		//machine
		"user_id",
		"name",
		"description",
		"count",
	}
	profileQuery = `SELECT projections.users6.id,` +
		` projections.users6.creation_date,` +
		` projections.users6.change_date,` +
		` projections.users6.resource_owner,` +
		` projections.users6.sequence,` +
		` projections.users6_humans.user_id,` +
		` projections.users6_humans.first_name,` +
		` projections.users6_humans.last_name,` +
		` projections.users6_humans.nick_name,` +
		` projections.users6_humans.display_name,` +
		` projections.users6_humans.preferred_language,` +
		` projections.users6_humans.gender,` +
		` projections.users6_humans.avatar_key` +
		` FROM projections.users6` +
		` LEFT JOIN projections.users6_humans ON projections.users6.id = projections.users6_humans.user_id AND projections.users6.instance_id = projections.users6_humans.instance_id`
	profileCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"user_id",
		"first_name",
		"last_name",
		"nick_name",
		"display_name",
		"preferred_language",
		"gender",
		"avatar_key",
	}
	emailQuery = `SELECT projections.users6.id,` +
		` projections.users6.creation_date,` +
		` projections.users6.change_date,` +
		` projections.users6.resource_owner,` +
		` projections.users6.sequence,` +
		` projections.users6_humans.user_id,` +
		` projections.users6_humans.email,` +
		` projections.users6_humans.is_email_verified` +
		` FROM projections.users6` +
		` LEFT JOIN projections.users6_humans ON projections.users6.id = projections.users6_humans.user_id AND projections.users6.instance_id = projections.users6_humans.instance_id`
	emailCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"user_id",
		"email",
		"is_email_verified",
	}
	phoneQuery = `SELECT projections.users6.id,` +
		` projections.users6.creation_date,` +
		` projections.users6.change_date,` +
		` projections.users6.resource_owner,` +
		` projections.users6.sequence,` +
		` projections.users6_humans.user_id,` +
		` projections.users6_humans.phone,` +
		` projections.users6_humans.is_phone_verified` +
		` FROM projections.users6` +
		` LEFT JOIN projections.users6_humans ON projections.users6.id = projections.users6_humans.user_id AND projections.users6.instance_id = projections.users6_humans.instance_id`
	phoneCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"user_id",
		"phone",
		"is_phone_verified",
	}
	userUniqueQuery = `SELECT projections.users6.id,` +
		` projections.users6.state,` +
		` projections.users6.username,` +
		` projections.users6_humans.user_id,` +
		` projections.users6_humans.email,` +
		` projections.users6_humans.is_email_verified` +
		` FROM projections.users6` +
		` LEFT JOIN projections.users6_humans ON projections.users6.id = projections.users6_humans.user_id AND projections.users6.instance_id = projections.users6_humans.instance_id`
	userUniqueCols = []string{
		"id",
		"state",
		"username",
		"user_id",
		"email",
		"is_email_verified",
	}
	notifyUserQuery = `SELECT projections.users6.id,` +
		` projections.users6.creation_date,` +
		` projections.users6.change_date,` +
		` projections.users6.resource_owner,` +
		` projections.users6.sequence,` +
		` projections.users6.state,` +
		` projections.users6.type,` +
		` projections.users6.username,` +
		` login_names.loginnames,` +
		` preferred_login_name.login_name,` +
		` projections.users6_humans.user_id,` +
		` projections.users6_humans.first_name,` +
		` projections.users6_humans.last_name,` +
		` projections.users6_humans.nick_name,` +
		` projections.users6_humans.display_name,` +
		` projections.users6_humans.preferred_language,` +
		` projections.users6_humans.gender,` +
		` projections.users6_humans.avatar_key,` +
		` projections.users6_notifications.user_id,` +
		` projections.users6_notifications.last_email,` +
		` projections.users6_notifications.verified_email,` +
		` projections.users6_notifications.last_phone,` +
		` projections.users6_notifications.verified_phone,` +
		` projections.users6_notifications.password_set,` +
		` COUNT(*) OVER ()` +
		` FROM projections.users6` +
		` LEFT JOIN projections.users6_humans ON projections.users6.id = projections.users6_humans.user_id AND projections.users6.instance_id = projections.users6_humans.instance_id` +
		` LEFT JOIN projections.users6_notifications ON projections.users6.id = projections.users6_notifications.user_id AND projections.users6.instance_id = projections.users6_notifications.instance_id` +
		` LEFT JOIN` +
		` (` + loginNamesQuery + `) AS login_names` +
		` ON login_names.user_id = projections.users6.id AND login_names.instance_id = projections.users6.instance_id` +
		` LEFT JOIN` +
		` (` + preferredLoginNameQuery + `) AS preferred_login_name` +
		` ON preferred_login_name.user_id = projections.users6.id AND preferred_login_name.instance_id = projections.users6.instance_id`
	notifyUserCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"state",
		"type",
		"username",
		"loginnames",
		"login_name",
		//human
		"user_id",
		"first_name",
		"last_name",
		"nick_name",
		"display_name",
		"preferred_language",
		"gender",
		"avatar_key",
		//machine
		"user_id",
		"last_email",
		"verified_email",
		"last_phone",
		"verified_phone",
		"password_set",
		"count",
	}
	usersQuery = `SELECT projections.users6.id,` +
		` projections.users6.creation_date,` +
		` projections.users6.change_date,` +
		` projections.users6.resource_owner,` +
		` projections.users6.sequence,` +
		` projections.users6.state,` +
		` projections.users6.type,` +
		` projections.users6.username,` +
		` login_names.loginnames,` +
		` preferred_login_name.login_name,` +
		` projections.users6_humans.user_id,` +
		` projections.users6_humans.first_name,` +
		` projections.users6_humans.last_name,` +
		` projections.users6_humans.nick_name,` +
		` projections.users6_humans.display_name,` +
		` projections.users6_humans.preferred_language,` +
		` projections.users6_humans.gender,` +
		` projections.users6_humans.avatar_key,` +
		` projections.users6_humans.email,` +
		` projections.users6_humans.is_email_verified,` +
		` projections.users6_humans.phone,` +
		` projections.users6_humans.is_phone_verified,` +
		` projections.users6_machines.user_id,` +
		` projections.users6_machines.name,` +
		` projections.users6_machines.description,` +
		` COUNT(*) OVER ()` +
		` FROM projections.users6` +
		` LEFT JOIN projections.users6_humans ON projections.users6.id = projections.users6_humans.user_id AND projections.users6.instance_id = projections.users6_humans.instance_id` +
		` LEFT JOIN projections.users6_machines ON projections.users6.id = projections.users6_machines.user_id AND projections.users6.instance_id = projections.users6_machines.instance_id` +
		` LEFT JOIN` +
		` (` + loginNamesQuery + `) AS login_names` +
		` ON login_names.user_id = projections.users6.id AND login_names.instance_id = projections.users6.instance_id` +
		` LEFT JOIN` +
		` (` + preferredLoginNameQuery + `) AS preferred_login_name` +
		` ON preferred_login_name.user_id = projections.users6.id AND preferred_login_name.instance_id = projections.users6.instance_id`
	usersCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"state",
		"type",
		"username",
		"loginnames",
		"login_name",
		//human
		"user_id",
		"first_name",
		"last_name",
		"nick_name",
		"display_name",
		"preferred_language",
		"gender",
		"avatar_key",
		"email",
		"is_email_verified",
		"phone",
		"is_phone_verified",
		//machine
		"user_id",
		"name",
		"description",
		"count",
	}
)
