terraform {
  required_providers {
    satori = {
      version = "~>1.0.0"
      source  = "satoricyber.com/terraform/satori"
    }
  }
}

provider "satori" {
  #can be provided via environment variable: SATORI_SA
  service_account     = "57beb9b4-2d08-47ee-981a-5b33dfa18cd8"
  #can be provided via environment variable: SATORI_SA_KEY
  service_account_key = "UI18V8XwTVBRVSLG6D2CYwBmmj7vTRtnRO1fJt0tsnmJI9lk3vxABPOQWHq9urQ8"
  #satori account id for all resources in this terraform
  satori_account      = "450338ff-5fb6-44a9-8662-2f3a9f4d2a76"

  #for local env only:
  url        = "http://127.0.0.1:8014"
  verify_tls = false
}
locals {
  dataaccesscontroller_id = "e0293f3a-e527-47bf-8885-404d43c2608b"
}
resource "satori_datasource_definition" "datastore3" {
  definition {
    name         = "datastore"
    hostname     = join(name, "terraform.com")
    port         = 8080
    projectids   = {a:121,b:233}
    customIngressPort= 3939


  }
}
