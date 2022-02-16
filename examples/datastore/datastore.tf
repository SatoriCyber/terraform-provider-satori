locals {
  dataaccess_controller_id = "<assigned dataaccess_controller_id>"
}
resource "satori_datastore" "datastore0" {
  name                     = "exampleDatastore"
  hostname                 = "data.source.target.hostname"
  dataaccess_controller_id = local.dataaccess_controller_id
  type                     = "SNOWFLAKE"
  origin_port               = 8081
  identity_provider_id     = "aaaaaaaaaaaaa-ddddd-ddddddddd-dddddddd"

}
# output of generated id for newly created datastore
output "datastore_created_id" {
  value = satori_datastore.datastore0.id
}
