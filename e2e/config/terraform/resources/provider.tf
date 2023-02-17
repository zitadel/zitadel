terraform {
  required_providers {
    zitadel = {
      source  = "zitadel/zitadel"
      version = "1.0.0-alpha.11"
    }
  }
}

provider "zitadel" {
  domain   = "localhost"
  insecure = "true"
  port     = "8080"
  token    = "../machinekey/zitadel-admin-sa.json"
}
