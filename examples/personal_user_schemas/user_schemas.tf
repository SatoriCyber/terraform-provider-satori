/*
Problem: When each user in a data warehouse has their own personal schema, sort of a scratch space where
they need full control over, it can be challenging to manage at scale in the UI, as each user 
would require a separate dataset and user access rule.

Solution: use Terraform to automate the process of managing these datasets. A simple file with the information
on the users who need personal schemas is loaded, and the corresponding Satori objects are created.
*/

locals {
  # Loads the users file
  raw_users_data              = yamldecode(file("users.yaml"))
  users = (
    length(keys(local.raw_users_data)) > 0 && local.raw_users_data["users"] != null
  ) ? local.raw_users_data["users"] : []

  # Assuming the data store was already created either manually
  # (enter its ID) or by another resource (take it from the resource's output)
  datastore_id       = "fda6cd66-67da-490c-be16-92eab04c2136"

  # The name of the database that contains the personal schemas
  database_name      = "scratchspace"
}

resource "satori_dataset" "dataset_for_personal_schema" {
  for_each = { for user in local.users : user["name"] => user }
  
  definition {
    name = "Personal Schema - ${each.value.name}"
    description = "This is a dataset granting access for ${each.value.name}'s personal schema"

    include_location {
      datastore = local.datastore_id
      location_parts = [local.database_name, each.value.schema] // S3 example
    }      
  }

  access_control_settings {
    enable_access_control = true
  }

}

resource "satori_access_rule" "access_rule_for_personal_schame" {
  for_each = { for user in local.users : user["name"] => user }

  parent_data_policy = satori_dataset.dataset_for_personal_schema[each.key].data_policy_id
  access_level = "OWNER"
  identity {
    type = "USER"
    name = each.value.email
  }

}