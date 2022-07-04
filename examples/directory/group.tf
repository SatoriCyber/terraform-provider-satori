resource "satori_directory_group" "group1" {
  name = "group1"
  description = "group from terraform"

  member {
    type = "USERNAME"
    name = "username"
  }

  member {
    type = "IDP_GROUP"
    name = "group_name1"
    #OKTA, AZURE, ONELOGIN
    identity_provider = "OKTA"
  }

  member {
    type = "DB_ROLE"
    name = "role_name"
    #SNOWFLAKE, REDSHIFT, BIGQUERY, POSTGRESQL, ATHENA, MSSQL, SYNAPSE
    data_store_type = "SNOWFLAKE"
  }
}

resource "satori_directory_group" "group_in_group" {
  name = "group_in_group"
  description = "group with group from terraform"

  member {
    type = "USERNAME"
    name = "name"
  }

  member {
    type = "DIRECTORY_GROUP"
    group_id = satori_directory_group.group1.id
  }
}

resource "satori_directory_group" "empty_group" {
  name        = "group4"
  description = "Empty directory group"
}
