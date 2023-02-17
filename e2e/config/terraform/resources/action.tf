resource "zitadel_action" "sleep_five_seconds" {
  org_id          = zitadel_org.actions.id
  name            = "sleepFiveSeconds"
  script          = data.local_file.sleep_five_seconds.content
  timeout         = "10s"
  allowed_to_fail = false
}
