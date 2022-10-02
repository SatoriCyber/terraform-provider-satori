data "satori_data_access_controller" "dac0" {
  type = "<assigned type>"
  region = "<assigned region>"
  cloud_provider = "<assigned cloud provider>"
  unique_name = "<assigned unique name>"
}

resource "satori_datastore" "datastore0" {
  name                     = "exampleDatastore"
  hostname                 = "data.source.target.hostname"
  dataaccess_controller_id = data.satori_data_access_controller.dac0.id
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
}

resource "satori_datastore" "datastoreWithIgnorePasswordUpdate" {
  // lifecycle.ignore_changes should be used after first time creation in order to ignore password update as API does not return it.
  name                     = "exampleDatastore"
  hostname                 = "data.source.target.hostname"
  dataaccess_controller_id = data.satori_data_access_controller.dac0.id
  type                     = "SNOWFLAKE"
  origin_port              = 8081
  satori_auth_settings {
    enabled = true
    credentials {
      password = "*********"
      username = "adminuser"
    }
  }
  lifecycle {
    ignore_changes = [
      satori_auth_settings.0.credentials.0.password
    ]
  }
  network_policy {}
}

# output of generated id for newly created datastore
output "datastore_created_id" {
  value = satori_datastore.datastore0.id
}
