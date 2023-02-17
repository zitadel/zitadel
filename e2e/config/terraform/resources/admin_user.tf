

resource "zitadel_human_user" "human_admin_user" {
  org_id             = data.zitadel_org.default.id
  user_name          = "zitadel-admin@zitadel.localhost"
  first_name         = "firstname"
  last_name          = "lastname"
  nick_name          = "nickname"
  display_name       = "displayname"
  preferred_language = "de"
  gender             = "GENDER_MALE"
  phone              = "+41799999999"
  is_phone_verified  = true
  email              = "test@zitadel.com"
  is_email_verified  = true
  initial_password   = "Password1!"
}

resource "zitadel_instance_member" "human_admin_user_member" {
  user_id = zitadel_human_user.human_admin_user.id
  roles   = ["IAM_OWNER"]
}
