provider "satori" {
  #can be provided via environment variable: SATORI_SA
  service_account     = "aaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaaa"
  #can be provided via environment variable: SATORI_SA_KEY
  service_account_key = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
  #satori account id for all resources in this terraform
  satori_account      = "cccccc-cccc-cccc-cccc-cccccccccccc"
  url        = "https://<satori-mgmt-server-address>:8014"
  verify_tls = true
}
locals {
  dataaccess_controller_id = "<assigned dataaccess_controller_id>"
}
resource "satori_datastore" "datastore0" {
  name                     = "exampleDatastore"
  hostname                 = "data.source.target.hostname"
  dataaccess_controller_id = local.dataaccess_controller_id
  type                     = "SNOWFLAKE"
  originPort                     = 8081
  identity_provider_id     = "aaaaaaaaaaaaa-ddddd-ddddddddd-dddddddd"

}
# output of generated id for newly created datastore
output "datastore_created_id" {
  value = satori_datastore.datastore0.id
}
