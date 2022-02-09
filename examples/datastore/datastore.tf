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
resource "satori_datastore" "datastore0" {
  name                     = "exampleDatastore"
  hostname                 = "data source target hostname"
  dataaccess_controller_id = local.dataaccess_controller_id
  type                     = "SNOWFLAKE"
  port                     = 8081
  custom_ingress_port      = 8083
  identity_provider_id     = "saml authXXXX"

}
# output of generated id for newly created datastore
output "datastore_created_id" {
  value = satori_datastore.datastore0.id
}
