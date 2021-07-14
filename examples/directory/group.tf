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

resource "satori_directory_group" "group2" {
  name = "group2"
  description = "group with group from terraform"

  member {
    type = "USERNAME"
    name = "name"
  }

  member {
    type = "DIRECTORY_GROUP"
    name = satori_directory_group.group1.name
    group_id = satori_directory_group.group1.id
  }

}
