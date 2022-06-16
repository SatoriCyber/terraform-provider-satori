locals {
  dataaccess_controller_id = "<assigned dataaccess_controller_id>"
  some_identity_provider_id = "<assigned some_identity_provider_id>"
}
resource "satori_datastore" "datastore0" {
  name                     = "exampleDatastore"
  hostname                 = "data.source.target.hostname"
  dataaccess_controller_id = local.dataaccess_controller_id
  type                     = "SNOWFLAKE"
  origin_port              = 8081
  baseline_security_policy {
    unassociated_queries_category {
      query_action = "PASS"
    }
    unsupported_queries_category {
      query_action = "PASS"
    }
    exclusions {
    }
  }
  network_policy {}
  identity_provider_id     = local.some_identity_provider_id
  satori_auth_settings {
    enabled = false
    credentials {
      password = ""
      username = ""
    }
  }
}

# output of generated id for newly created datastore
output "datastore_created_id" {
  value = satori_datastore.datastore0.id
}
