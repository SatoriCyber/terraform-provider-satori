---
layout: ""
page_title: "security_policy (Resource)"
description: |-
security_policy resource allows defining security policies.
---

# security_policy (Resource)

The Satori Security Policy is a re-usable object that can be configured to contain multiple sets of dynamic masking configurations and data filtering configurations.
A Security Policy can be applied to on one or more datasets.

The **security_policy** resource allows defining security policies.

## Example Usage

```terraform
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
        conditional_masking {
          where_condition = "country = 'US'"
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
          location_path = "db2.schema2.table"
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
          location_parts = ["db2", "schema2", "table1"]
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String) Security policy name.

### Optional

- **id** (String) The ID of this resource.
- **profile** (Block List, Max: 1) Security policy profile. (see [below for nested schema](#nestedblock--profile))

<a id="nestedblock--profile"></a>
### Nested Schema for `profile`

Optional:

- **masking** (Block List, Max: 1) Masking profile. (see [below for nested schema](#nestedblock--profile--masking))
- **row_level_security** (Block List, Max: 1) Row level security profile (see [below for nested schema](#nestedblock--profile--row_level_security))

<a id="nestedblock--profile--masking"></a>
### Nested Schema for `profile.masking`

Required:

- **active** (Boolean) Is active.

Optional:

- **rule** (Block List) Masking Rule. (see [below for nested schema](#nestedblock--profile--masking--rule))

<a id="nestedblock--profile--masking--rule"></a>
### Nested Schema for `profile.masking.rule`

Required:

- **action** (Block List, Min: 1, Max: 1) Rule action. (see [below for nested schema](#nestedblock--profile--masking--rule--action))
- **active** (Boolean) Is active rule.
- **criteria** (Block List, Min: 1, Max: 1) Masking criteria. (see [below for nested schema](#nestedblock--profile--masking--rule--criteria))
- **description** (String) Rule description.
- **id** (String) Rule id, has to be unique.

Optional:

- **conditional_masking** (Block List, Max: 1) Conditional masking. (see [below for nested schema](#nestedblock--profile--masking--rule--conditional_masking))

<a id="nestedblock--profile--masking--rule--action"></a>
### Nested Schema for `profile.masking.rule.conditional_masking`

Required:

- **masking_profile_id** (String) The reference id to be applied as masking profile.

Optional:

- **type** (String) Rule type. Defaults to `APPLY_MASKING_PROFILE`.


<a id="nestedblock--profile--masking--rule--criteria"></a>
### Nested Schema for `profile.masking.rule.conditional_masking`

Required:

- **condition** (String) Identity condition, for example IS_NOT, IS, etc.
- **identity** (Block List, Min: 1, Max: 1) Identity to apply the rule for. (see [below for nested schema](#nestedblock--profile--masking--rule--conditional_masking--identity))

<a id="nestedblock--profile--masking--rule--conditional_masking--identity"></a>
### Nested Schema for `profile.masking.rule.conditional_masking.identity`

Required:

- **type** (String) Identity type, valid types are: USER, DB_USER, IDP_GROUP, GROUP, DATABRICKS_GROUP, DATABRICKS_SERVICE_PRINCIPAL, SNOWFLAKE_ROLE, CEL, EVERYONE.
Can not be changed after creation.

Optional:

- **group_id** (String) Directory group ID for identity of type GROUP.
Can not be changed after creation.
- **name** (String) User/group name for identity types of USER and IDP_GROUP or a custom expression based on attributes of the identity for CEL identity type.
Can not be changed after creation.



<a id="nestedblock--profile--masking--rule--conditional_masking"></a>
### Nested Schema for `profile.masking.rule.conditional_masking`

Optional:

- **where_condition** (String) Where condition.




<a id="nestedblock--profile--row_level_security"></a>
### Nested Schema for `profile.row_level_security`

Required:

- **active** (Boolean) Row level security activation.

Optional:

- **mapping** (Block List) Row Level Security Mapping. (see [below for nested schema](#nestedblock--profile--row_level_security--mapping))
- **rule** (Block List) Row Level Security Rule definition. (see [below for nested schema](#nestedblock--profile--row_level_security--rule))

<a id="nestedblock--profile--row_level_security--mapping"></a>
### Nested Schema for `profile.row_level_security.mapping`

Required:

- **defaults** (Block List, Min: 1, Max: 1) A list of default values to be applied in this filter if there was no match. Values are dependent on their type and has to be homogeneous (see [below for nested schema](#nestedblock--profile--row_level_security--mapping--defaults))
- **filter** (Block List, Min: 1, Max: 1) Filter definition. (see [below for nested schema](#nestedblock--profile--row_level_security--mapping--filter))
- **name** (String) Filter name, has to be unique in this policy.

<a id="nestedblock--profile--row_level_security--mapping--defaults"></a>
### Nested Schema for `profile.row_level_security.mapping.name`

Required:

- **type** (String) Default values type. Allowed options: STRING, NUMERIC, CEL, SQL, ANY_VALUE, ALL_OTHER_VALUES
- **value** (List of String) List of values, when NO_VALUE or ALL_OTHER_VALUES are defined, the list has to be empty


<a id="nestedblock--profile--row_level_security--mapping--filter"></a>
### Nested Schema for `profile.row_level_security.mapping.name`

Required:

- **criteria** (Block List, Min: 1, Max: 1) Filter criteria. (see [below for nested schema](#nestedblock--profile--row_level_security--mapping--name--criteria))
- **values** (Block List, Min: 1, Max: 1) A list of values to be applied in this filter. Values are dependent on their type and has to be homogeneous (see [below for nested schema](#nestedblock--profile--row_level_security--mapping--name--values))

<a id="nestedblock--profile--row_level_security--mapping--name--criteria"></a>
### Nested Schema for `profile.row_level_security.mapping.name.criteria`

Required:

- **condition** (String) Identity condition, for example IS_NOT, IS, etc.
- **identity** (Block List, Min: 1, Max: 1) Identity to apply the rule for. (see [below for nested schema](#nestedblock--profile--row_level_security--mapping--name--criteria--identity))

<a id="nestedblock--profile--row_level_security--mapping--name--criteria--identity"></a>
### Nested Schema for `profile.row_level_security.mapping.name.criteria.identity`

Required:

- **type** (String) Identity type, valid types are: USER, DB_USER, IDP_GROUP, GROUP, DATABRICKS_GROUP, DATABRICKS_SERVICE_PRINCIPAL, SNOWFLAKE_ROLE, CEL, EVERYONE.
Can not be changed after creation.

Optional:

- **group_id** (String) Directory group ID for identity of type GROUP.
Can not be changed after creation.
- **name** (String) User/group name for identity types of USER and IDP_GROUP or a custom expression based on attributes of the identity for CEL identity type.
Can not be changed after creation.



<a id="nestedblock--profile--row_level_security--mapping--name--values"></a>
### Nested Schema for `profile.row_level_security.mapping.name.values`

Required:

- **type** (String) Values type. Allowed options: STRING, NUMERIC, CEL, SQL, ANY_VALUE, ALL_OTHER_VALUES
- **value** (List of String) List of values, when ANY_VALUE or ALL_OTHER_VALUES are defined, the list has to be empty




<a id="nestedblock--profile--row_level_security--rule"></a>
### Nested Schema for `profile.row_level_security.rule`

Required:

- **active** (Boolean) Is active rule.
- **description** (String) Rule description.
- **filter** (Block List, Min: 1, Max: 1) Rule filter. (see [below for nested schema](#nestedblock--profile--row_level_security--rule--filter))
- **id** (String) Rule id, has to be unique.

<a id="nestedblock--profile--row_level_security--rule--filter"></a>
### Nested Schema for `profile.row_level_security.rule.id`

Required:

- **datastore_id** (String) Datastore ID.
- **logic_yaml** (String) Conditional rule, for more info see https://satoricyber.com/docs/security-policies/#setting-up-data-filtering.

Optional:

- **advanced** (Boolean) Describes if logic yaml contains complex configuration. Defaults to `true`.
- **location** (Block List, Deprecated) Location to be included in the rule. The 'location' field has been deprecated. Please use the 'location_path', `location_parts` or `location_parts_full` fields instead. (see [below for nested schema](#nestedblock--profile--row_level_security--rule--id--location))
- **location_parts** (List of String) The part separated location path in the data store. Includes an array of path parts when part types are defined with default definitions. For example ['a', 'b', 'c'] in Snowflake data store will path to table 'a' under schema 'b' under database 'a'. Conflicts with 'location', 'location_path', and 'location_parts_full' fields
- **location_parts_full** (Block List) The full location path definition in the data store. Includes an array of objects with path name and path type. Can be used when the path type should be defined explicitly and not as defined by default. For example [{name= 'a', type= 'DATABASE'},{name= 'b', type= 'SCHEMA'},{name= 'view.c', type= 'VIEW'}]. Conflicts with 'location', 'location_path', and 'location_parts' fields. (see [below for nested schema](#nestedblock--profile--row_level_security--rule--id--location_parts_full))
- **location_path** (String) The short presentation of the location path in the data store. Includes `.` separated string when part types are defined with default definitions. For example 'a.b.c' in Snowflake data store will path to table 'a' under schema 'b' under database 'a'.  Conflicts with 'location', 'location_parts', and 'location_parts_full' fields.

<a id="nestedblock--profile--row_level_security--rule--id--location"></a>
### Nested Schema for `profile.row_level_security.rule.id.location`

Optional:

- **athena_location** (Block List, Max: 1) Location for Athena data store. (see [below for nested schema](#nestedblock--profile--row_level_security--rule--id--location--athena_location))
- **mongo_location** (Block List, Max: 1) Location for MongoDB data store. (see [below for nested schema](#nestedblock--profile--row_level_security--rule--id--location--mongo_location))
- **mysql_location** (Block List, Max: 1) Location for MySql and MariaDB data stores. (see [below for nested schema](#nestedblock--profile--row_level_security--rule--id--location--mysql_location))
- **relational_location** (Block List, Max: 1) Location for a relational data store. (see [below for nested schema](#nestedblock--profile--row_level_security--rule--id--location--relational_location))
- **s3_location** (Block List, Max: 1) Location for S3 data store. (see [below for nested schema](#nestedblock--profile--row_level_security--rule--id--location--s3_location))

<a id="nestedblock--profile--row_level_security--rule--id--location--athena_location"></a>
### Nested Schema for `profile.row_level_security.rule.id.location.s3_location`

Required:

- **catalog** (String) Catalog name.

Optional:

- **db** (String) Database name.
- **table** (String) Table name.


<a id="nestedblock--profile--row_level_security--rule--id--location--mongo_location"></a>
### Nested Schema for `profile.row_level_security.rule.id.location.s3_location`

Required:

- **db** (String) Database name.

Optional:

- **collection** (String) Collection name.


<a id="nestedblock--profile--row_level_security--rule--id--location--mysql_location"></a>
### Nested Schema for `profile.row_level_security.rule.id.location.s3_location`

Required:

- **db** (String) Database name.

Optional:

- **table** (String) Table name.


<a id="nestedblock--profile--row_level_security--rule--id--location--relational_location"></a>
### Nested Schema for `profile.row_level_security.rule.id.location.s3_location`

Required:

- **db** (String) Database name.

Optional:

- **schema** (String) Schema name.
- **table** (String) Table name.


<a id="nestedblock--profile--row_level_security--rule--id--location--s3_location"></a>
### Nested Schema for `profile.row_level_security.rule.id.location.s3_location`

Required:

- **bucket** (String) Bucket name.

Optional:

- **object_key** (String) Object Key name.



<a id="nestedblock--profile--row_level_security--rule--id--location_parts_full"></a>
### Nested Schema for `profile.row_level_security.rule.id.location_parts_full`

Required:

- **name** (String) The name of the location part.
- **type** (String) The type of the location part. Optional values: TABLE, COLUMN, SEMANTIC_MODEL, REPORT, DASHBOARD, DATABASE, SCHEMA, JSON_PATH, WAREHOUSE, ENDPOINT, TYPE, FIELD, EXTERNAL_LOCATION, CATALOG, BUCKET, OBJECT, COLLECTION, VIEW, etc