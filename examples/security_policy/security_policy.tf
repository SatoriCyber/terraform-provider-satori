resource "satori_security_policy" "security_policy" {
  name = "security policy terraform"
  profile {
    masking {
      active = true
      rule {
        id          = "3"
        description = "rule 1"
        active      = true
        criteria {
          condition = "IS"
          identity {
            type = "USER"
            name = "test-user"
          }
        }
        action {
          type = "APPLY_MASKING_PROFILE" # optional
          masking_profile_id = satori_masking_profile.masking_profile.id # as reference of previously created masking profile
        }
      }
      rule {
        id          = "1"
        description = "rule 2"
        active      = true
        criteria {
          condition = "IS_NOT"
          identity {
            type = "IDP_GROUP"
            name = "test-group"
          }
        }
        action {
          type               = "APPLY_MASKING_PROFILE"
          masking_profile_id = satori_masking_profile.masking_profile.id # as reference of previously created masking profile
        }
      }
      rule {
        id          = "2"
        description = "rule 3"
        active      = true
        criteria {
          condition = "IS_NOT"
          identity {
            type     = "GROUP"
            group_id = satori_directory_group.group2.id
          }
        }
        action {
          type               = "APPLY_MASKING_PROFILE"
          masking_profile_id = satori_masking_profile.masking_profile.id
        }
      }
      rule {
        id          = "4"
        description = "rule 4"
        active      = true
        criteria {
          condition = "IS_NOT"
          identity {
            type     = "GROUP"
            group_id = satori_directory_group.group3.id
          }
        }
        action {
          # type = "APPLY_MASKING_PROFILE"
          masking_profile_id = satori_masking_profile.masking_profile.id
        }
      }
    }
    row_level_security {
      active = false
      rule {
        id          = "rls1"
        description = "rls 1 description"
        active      = true
        filter {
          datastore_id = local.datastore_id
          advanced     = false
          logic_yaml   = <<EOT
                          field:
                            name: '33'
                          filterName: Filter 1
                        EOT
          location {
            relational_location {
              db     = "db2"
              schema = "schema2"
              table  = "table"
            }
          }
        }
      }
      rule {
        id          = "rls2"
        description = "rls 1 description"
        active      = true
        filter {
          datastore_id = local.datastore_id
          logic_yaml   = <<EOT
and:
  - or:
    - and:
      - field:
          name: c1
          path: $.a['b']
        filterName: Filter 1
      - field:
          name: c2
        filterName: with space
    - and:
      - field:
          name: c3
        filterName: Filter 1
      - field:
          name: c4
          path: $.a.b
        filterName: Filter 2
  - or:
    - field:
        name: lala
      filterName: Filter 1
    - field:
        name: lala
        path: $.a.b
      filterName: Filter 1
    - field:
        name: d3
        path: $.a['b']
      filterName: Filter 1
  EOT
          location_prefix { // usage of the deprecated field 'location_prefix'
            db     = "db2"
            schema = "schema2"
            table  = "table1"
          }
        }
      }
      mapping {
        name = "Filter 1"
        filter {
          criteria {
            condition = "IS_NOT"
            identity {
              type     = "GROUP"
              group_id = satori_directory_group.group3.id
            }
          }
          values {
            type  = "STRING"
            value = ["a", "b"]
          }
        }
        defaults {
          type  = "ALL_OTHER_VALUES"
          value = []
        }
      }
      mapping {
        name = "Filter 2"
        filter {
          criteria {
            condition = "IS_NOT"
            identity {
              type     = "GROUP"
              group_id = satori_directory_group.group3.id
            }
          }
          values {
            type  = "STRING"
            value = ["a", "b"]
          }
        }
        defaults {
          type  = "NO_VALUE"
          value = []
        }
      }
      mapping {
        name = "with space"
        filter {
          criteria {
            condition = "IS_NOT"
            identity {
              type     = "GROUP"
              group_id = satori_directory_group.group3.id
            }
          }
          values {
            type  = "STRING"
            value = ["a", "b"]
          }
        }
        defaults {
          type  = "NUMERIC"
          value = ["1", "2.4", "1.343434"]
        }
      }
    }
  }
}