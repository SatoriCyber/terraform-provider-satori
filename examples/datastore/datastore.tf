locals {
  dataaccess_controller_id = "<assigned dataaccess_controller_id>"
}

data "satori_data_access_controller" "public_dac" {
  type = "PUBLIC"
  region = "<assigned region>"
  cloud_provider = "<assigned cloud provider>"
}

data "satori_data_access_controller" "private_dac" {
  type = "<assigned type - PRIVATE | PRIVATE_MANAGED>"
  id = "<assigned id>"
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
}

resource "satori_datastore" "datastoreWithIgnorePasswordUpdate" {
  // lifecycle.ignore_changes should be used after first time creation in order to ignore password update as API does not return it.
  name                     = "exampleDatastore"
  hostname                 = "data.source.target.hostname"
  dataaccess_controller_id = data.satori_data_access_controller.public_dac.id
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

resource "satori_datastore" "datastoreWithPrivateDac" {
  // lifecycle.ignore_changes should be used after first time creation in order to ignore password update as API does not return it.
  name                     = "exampleDatastore"
  hostname                 = "data.source.target.hostname"
  dataaccess_controller_id = data.satori_data_access_controller.private_dac.id
  type                     = "SNOWFLAKE"
  origin_port              = 8081
  lifecycle {
    ignore_changes = [
      satori_auth_settings.0.credentials.0.password
    ]
  }
  network_policy {}
}

resource "satori_datastore" "mongodbDatastore" {
  name                     = "mongoExample"
  hostname                 = "mongo.example.mongodb.net"
  dataaccess_controller_id = data.satori_data_access_controller.public_dac.id
  type                     = "MONGO"
  datastore_settings {
    deployment_type = "MONGODB_SRV"
  }
  network_policy {}
}

# output of generated id for newly created datastore
output "datastore_created_id" {
  value = satori_datastore.datastore0.id
}
