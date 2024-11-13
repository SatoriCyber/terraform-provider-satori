locals {
  dataaccess_controller_id = "<assigned dataaccess_controller_id>"
}

resource "satori_datastore" "databricks_example" {
  // lifecycle.ignore_changes should be used after first time creation in order to ignore password update as API does not return it.
  name = "Databricks Datastore Example"

  hostname                 = "<databricks_instance>" // Databricks Instance (workspace host)
  dataaccess_controller_id = local.dataaccess_controller_id
  type                     = "DATABRICKS"

  databricks_settings {
    account_id = "account_id"
    warehouse_id = "sql_warehouse_id"
    credentials {
      type         = "AWS_SERVICE_PRINCIPAL_TOKEN"
      client_id     = "application_client_id"
      client_secret = "*********"
    }
  }
  lifecycle {
    ignore_changes = [
      databricks_settings.0.credentials.0.client_secret
    ]
  }
}