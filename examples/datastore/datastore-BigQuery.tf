provider "satori" {
  #can be provided via environment variable: SATORI_SA
  service_account     = "aaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaaa"
  #can be provided via environment variable: SATORI_SA_KEY
  service_account_key = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
  #satori account id for all resources in this terraform
  satori_account      = "cccccc-cccc-cccc-cccc-cccccccccccc"
  url        = "https://<satori-mgmt-server-address>:8014"
  verify_tls = false
}
locals {
  dataaccess_controller_id = "<assigned dataaccess_controller_id>"
}

resource "satori_datastore" "datastore0" {
  name = "exampleDatastore"

  dataaccess_controller_id = local.dataaccess_controller_id
  # data source specific connection settings
  type                     = "BIGQUERY"
  project_ids              = [111, 222] #  BigQuery affected project ids
  hostname                 = "data source target hostname"
  origin_port              = 8081 # data source server's ip
  identity_provider_id     = "aaaaaaaaaaaaa-ddddd-ddddddddd-dddddddd"
  ####### BASELINE_POLICY SETTINGS #########
  baseline_security_policy {

    unassociated_queries_category {
      query_action = "REDACT" #Allowed: PASS┃REDACT┃BLOCK
    }
    unsupported_queries_category {
      query_action = "REDACT"
    }
    exclusions {

      excluded_identities {
        identity_type = "USER"
        identity = "user1"
      }
      excluded_identities {
        identity_type = "USER"
        identity = "user2"
      }
      excluded_query_patterns {
        pattern = ".*a.*"
      }
      excluded_query_patterns {
        pattern = ".*b.*"
      }
    }
  }
}


output "datastore_created_id" {
  value = satori_datastore.datastore0.id
}
