provider "satori" {
  #can be provided via environment variable: SATORI_SA
  service_account     = "9bd8dd14-50c9-4d52-b148-bc0dece8b964"
  #can be provided via environment variable: SATORI_SA_KEY
  service_account_key = "HBRwevC9uDOEYv5LIhNvtzdrgc7WTErqWJE3iqmTFxFbMFGkcEFvjiGqm3BlWACd"
  #satori account id for all resources in this terraform
  satori_account      = "fdd00136-69f2-471a-9b9e-b8ccb9658b81"

  url        = "https://<satori-mgmt-server-address>:8014"
  verify_tls = false
}
locals {
  dataaccess_controller_id = "<assigned dataaccess_controller_id>"
}
locals {
  dataaccess_controller_id = "509bf523-e5d9-436b-9b8c-7e9ee18c95fc"
  #  dataaccesscontroller_id = "e0293f3ae52747bf8885404d43c2608b"
}
resource "satori_datastore" "datastore0" {
  name = "exampleDatastore"

  dataaccess_controller_id = local.dataaccess_controller_id
  # data source specific connection settings
  type                     = "BIGQUERY"
  project_ids              = [111, 222] #  port to connect to
  hostname                 = "data source target hostname"
  port                     = 8081 # data source server's ip
  custom_ingress_port      = 8083
  identity_provider_id     = "saml string"
  ####### BASELINE_POLICY SETTINGS #########
  baseline_security_policy {

    unassociated_queries_category {
      query_action = "REDACT"
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
