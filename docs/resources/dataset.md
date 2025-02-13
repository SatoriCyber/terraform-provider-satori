---
layout: ""
page_title: "satori_dataset (Resource)"
description: |-
satori_dataset resource allows defining datasets.
---

# satori_dataset (Resource)

Datasets are collections of data store locations that are meant to be governed as a single unit.
The **satori_dataset** resource allows defining datasets.

<br />
<br />
The resource output includes **data_policy_id** which is mandatory ID for future access rule resources creation.
See Read-Only section and **satori_request_access_rule** Resource examples.

## Example Usage

```terraform
//get ID of user by email
data "satori_user" "data_steward1" {
  email = "data-steward@acme.organization"
}

resource "satori_dataset" "dataset1" {
  definition {
    name = "satori_dataset terraform test"
    description = "from satori terraform provider"
    #the service account must also be an owner to be able to modify settings beyond definition
    owners = ["522fb8ab-8d7b-4498-b39d-6911e2839253", data.satori_user.data_steward1.id]


    approvers {
      # Currently can be only IdP groups
      type = "GROUP"
      id   = "788680b7-461d-493a-a3d6-86e71fd01ff2"
    }

    approvers {
      type = "USER"
      id   = "3d174db4-4526-4469-2fda-46d0dd2a7f7d"
    }

    approvers {
      type = "USER"
      id   = data.satori_user.data_steward1.id
    }

    include_location {
      datastore = "12345678-95cf-474f-a1d6-d5084810dd95"
    }

    include_location {
      datastore     = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location_path = "db1"
    }

    include_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location_parts = ["db2", "schema1"]
    }

    include_location {
      datastore     = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location_path = "db2.schema2.table"
    }

    exclude_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location_parts = ["db2", "schema1", "tableX"]
    }
  }

  access_control_settings {
    enable_access_control = false
    enable_user_requests  = false
    enable_self_service   = false
  }

  custom_policy {
    #default priority is 100
    #priority = 100
    rules_yaml = file("${path.module}/rules.yaml")
    tags_yaml = file("${path.module}/tags.yaml")
  }

  security_policies = ["56412aff-6ecf-4eff-9b96-2e0f6ec36c42"]
}

// Example with different location types
resource "satori_dataset" "dataset2" {
  definition {
    name = "satori_dataset terraform test"
    description = "from satori terraform provider"
    #the service account must also be an owner to be able to modify settings beyond definition
    owners = ["522fb8ab-8d7b-4498-b39d-6911e2839253", data.satori_user.data_steward1.id]

    include_location {
      datastore = "12345678-95cf-474f-a1d6-d5084810dd95"
    }

    include_location {
      datastore     = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location_path = "db2.schema1"
    }

    include_location {
      // MongoDB example
      datastore     = "3go33ff5-95cf-474f-a1d6-d5084810dd5k"
      location_path = "db1.collection1"
    }

    include_location {
      // MongoDB example
      datastore = "3go33ff5-95cf-474f-a1d6-d5084810dd5k"
      location_parts = ["db1", "collection1"]
    }

    include_location {
      // MongoDB example
      datastore = "3go33ff5-95cf-474f-a1d6-d5084810dd5k"
      location_parts_full {
        name = "db1"
        type = "DATABASE"
      }
      location_parts_full {
        name = "collection1"
        type = "COLLECTION"
      }
    }

    include_location {
      datastore = "8kl43ff5-95cf-474f-a1d6-d508481049lw"
      // S3 example
      location_parts = ["bucket1", "a/b/c"] // S3 example
    }

    include_location {
      datastore = "8kl43ff5-95cf-474f-a1d6-d508481049lw"
      // S3 example
      location_parts_full {
        name = "bucket1"
        type = "BUCKET"
      }
      location_parts_full {
        name = "a/b/c"
        type = "OBJECT"
      }
    }

    exclude_location {
      datastore     = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location_path = "db1.schema1.tableX"
    }
  }

  access_control_settings {
    enable_access_control = false
    enable_user_requests  = false
    enable_self_service   = false
  }

  custom_policy {
    #default priority is 100
    #priority = 100
    rules_yaml = file("${path.module}/rules.yaml")
    tags_yaml = file("${path.module}/tags.yaml")
  }

  security_policies = ["56412aff-6ecf-4eff-9b96-2e0f6ec36c42"]
}

// Example with deprecated usage of location field (use location_path, location_parts or location_parts_full instead)
resource "satori_dataset" "dataset3" {
  definition {
    name = "satori_dataset terraform test"
    description = "from satori terraform provider"
    #the service account must also be an owner to be able to modify settings beyond definition
    owners = ["522fb8ab-8d7b-4498-b39d-6911e2839253", data.satori_user.data_steward1.id]

    include_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location {
        relational_location {
          db = "db1"
        }
      }
    }

    include_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location {
        relational_location {
          db     = "db2"
          schema = "schema1"
        }
      }
    }

    include_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location {
        relational_location {
          db     = "db2"
          schema = "schema2"
          table  = "table"
        }
      }
    }

    exclude_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      relational_location {
        db     = "db2"
        schema = "schema1"
        table  = "tableX"
      }
    }
  }

  access_control_settings {
    enable_access_control = false
    enable_user_requests  = false
    enable_self_service   = false
  }

  custom_policy {
    #default priority is 100
    #priority = 100
    rules_yaml = file("${path.module}/rules.yaml")
    tags_yaml = file("${path.module}/tags.yaml")
  }

  security_policies = ["56412aff-6ecf-4eff-9b96-2e0f6ec36c42"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **access_control_settings** (Block List, Min: 1, Max: 1) Dataset access controls. (see [below for nested schema](#nestedblock--access_control_settings))
- **definition** (Block List, Min: 1, Max: 1) Parameters for dataset definition. (see [below for nested schema](#nestedblock--definition))

### Optional

- **custom_policy** (Block List, Max: 1) Dataset custom policy. (see [below for nested schema](#nestedblock--custom_policy))
- **id** (String) The ID of this resource.
- **security_policies** (List of String) IDs of security policies to apply to this dataset.

### Read-Only

- **data_policy_id** (String) Parent ID for dataset permissions.

<a id="nestedblock--access_control_settings"></a>
### Nested Schema for `access_control_settings`

Optional:

- **enable_access_control** (Boolean) Enforce access control to this dataset. Defaults to `false`.


<a id="nestedblock--definition"></a>
### Nested Schema for `definition`

Required:

- **name** (String) Dataset name.

Optional:

- **approvers** (Block List) Identities of Satori users/groups that will be set as dataset approvers. (see [below for nested schema](#nestedblock--definition--approvers))
- **description** (String) Dataset description.
- **exclude_location** (Block List) Location to exclude from dataset. (see [below for nested schema](#nestedblock--definition--exclude_location))
- **include_location** (Block List) Location to include in dataset. (see [below for nested schema](#nestedblock--definition--include_location))
- **owners** (List of String) IDs of Satori users that will be set as dataset owners.

<a id="nestedblock--definition--approvers"></a>
### Nested Schema for `definition.approvers`

Required:

- **id** (String) The ID of the approver entity
- **type** (String) Approver type, can be either `GROUP` (IdP Group alone) or `USER`


<a id="nestedblock--definition--exclude_location"></a>
### Nested Schema for `definition.exclude_location`

Required:

- **datastore** (String) Data store ID.

Optional:

- **location** (Block List, Max: 1, Deprecated) Location for a data store. Can include only one location type field from the above: relational_location, mysql_location, athena_location, mongo_location and s3_location . Conflicts with 'location_path', 'location_parts' and 'location_parts_full' fields. The 'location' field has been deprecated. Please use the 'location_path', `location_parts` or `location_parts_full` fields instead. (see [below for nested schema](#nestedblock--definition--exclude_location--location))
- **location_parts** (List of String) The part separated location path in the data store. Includes an array of path parts when part types are defined with default definitions. For example ['a', 'b', 'c'] in Snowflake data store will path to table 'a' under schema 'b' under database 'a'. Conflicts with 'location', 'location_path', and 'location_parts_full' fields
- **location_parts_full** (Block List) The full location path definition in the data store. Includes an array of objects with path name and path type. Can be used when the path type should be defined explicitly and not as defined by default. For example [{name= 'a', type= 'DATABASE'},{name= 'b', type= 'SCHEMA'},{name= 'view.c', type= 'VIEW'}]. Conflicts with 'location', 'location_path', and 'location_parts' fields. (see [below for nested schema](#nestedblock--definition--exclude_location--location_parts_full))
- **location_path** (String) The short presentation of the location path in the data store. Includes `.` separated string when part types are defined with default definitions. For example 'a.b.c' in Snowflake data store will path to table 'a' under schema 'b' under database 'a'.  Conflicts with 'location', 'location_parts', and 'location_parts_full' fields.

<a id="nestedblock--definition--exclude_location--location"></a>
### Nested Schema for `definition.exclude_location.location`

Optional:

- **athena_location** (Block List, Max: 1) Location for Athena data store. (see [below for nested schema](#nestedblock--definition--exclude_location--location--athena_location))
- **mongo_location** (Block List, Max: 1) Location for MongoDB data store. (see [below for nested schema](#nestedblock--definition--exclude_location--location--mongo_location))
- **mysql_location** (Block List, Max: 1) Location for MySql and MariaDB data stores. (see [below for nested schema](#nestedblock--definition--exclude_location--location--mysql_location))
- **relational_location** (Block List, Max: 1) Location for a relational data store. (see [below for nested schema](#nestedblock--definition--exclude_location--location--relational_location))
- **s3_location** (Block List, Max: 1) Location for S3 data store. (see [below for nested schema](#nestedblock--definition--exclude_location--location--s3_location))

<a id="nestedblock--definition--exclude_location--location--athena_location"></a>
### Nested Schema for `definition.exclude_location.location.s3_location`

Required:

- **catalog** (String) Catalog name.

Optional:

- **db** (String) Database name.
- **table** (String) Table name.


<a id="nestedblock--definition--exclude_location--location--mongo_location"></a>
### Nested Schema for `definition.exclude_location.location.s3_location`

Required:

- **db** (String) Database name.

Optional:

- **collection** (String) Collection name.


<a id="nestedblock--definition--exclude_location--location--mysql_location"></a>
### Nested Schema for `definition.exclude_location.location.s3_location`

Required:

- **db** (String) Database name.

Optional:

- **table** (String) Table name.


<a id="nestedblock--definition--exclude_location--location--relational_location"></a>
### Nested Schema for `definition.exclude_location.location.s3_location`

Required:

- **db** (String) Database name.

Optional:

- **schema** (String) Schema name.
- **table** (String) Table name.


<a id="nestedblock--definition--exclude_location--location--s3_location"></a>
### Nested Schema for `definition.exclude_location.location.s3_location`

Required:

- **bucket** (String) Bucket name.

Optional:

- **object_key** (String) Object Key name.



<a id="nestedblock--definition--exclude_location--location_parts_full"></a>
### Nested Schema for `definition.exclude_location.location_parts_full`

Required:

- **name** (String) The name of the location part.
- **type** (String) The type of the location part. Optional values: TABLE, COLUMN, SEMANTIC_MODEL, REPORT, DASHBOARD, DATABASE, SCHEMA, JSON_PATH, WAREHOUSE, ENDPOINT, TYPE, FIELD, EXTERNAL_LOCATION, CATALOG, BUCKET, OBJECT, COLLECTION, VIEW, etc



<a id="nestedblock--definition--include_location"></a>
### Nested Schema for `definition.include_location`

Required:

- **datastore** (String) Data store ID.

Optional:

- **location** (Block List, Max: 1, Deprecated) Location for a data store. Can include only one location type field from the above: relational_location, mysql_location, athena_location, mongo_location and s3_location . Conflicts with 'location_path', 'location_parts' and 'location_parts_full' fields. The 'location' field has been deprecated. Please use the 'location_path', `location_parts` or `location_parts_full` fields instead. (see [below for nested schema](#nestedblock--definition--include_location--location))
- **location_parts** (List of String) The part separated location path in the data store. Includes an array of path parts when part types are defined with default definitions. For example ['a', 'b', 'c'] in Snowflake data store will path to table 'a' under schema 'b' under database 'a'. Conflicts with 'location', 'location_path', and 'location_parts_full' fields
- **location_parts_full** (Block List) The full location path definition in the data store. Includes an array of objects with path name and path type. Can be used when the path type should be defined explicitly and not as defined by default. For example [{name= 'a', type= 'DATABASE'},{name= 'b', type= 'SCHEMA'},{name= 'view.c', type= 'VIEW'}]. Conflicts with 'location', 'location_path', and 'location_parts' fields. (see [below for nested schema](#nestedblock--definition--include_location--location_parts_full))
- **location_path** (String) The short presentation of the location path in the data store. Includes `.` separated string when part types are defined with default definitions. For example 'a.b.c' in Snowflake data store will path to table 'a' under schema 'b' under database 'a'.  Conflicts with 'location', 'location_parts', and 'location_parts_full' fields.

<a id="nestedblock--definition--include_location--location"></a>
### Nested Schema for `definition.include_location.location`

Optional:

- **athena_location** (Block List, Max: 1) Location for Athena data store. (see [below for nested schema](#nestedblock--definition--include_location--location--athena_location))
- **mongo_location** (Block List, Max: 1) Location for MongoDB data store. (see [below for nested schema](#nestedblock--definition--include_location--location--mongo_location))
- **mysql_location** (Block List, Max: 1) Location for MySql and MariaDB data stores. (see [below for nested schema](#nestedblock--definition--include_location--location--mysql_location))
- **relational_location** (Block List, Max: 1) Location for a relational data store. (see [below for nested schema](#nestedblock--definition--include_location--location--relational_location))
- **s3_location** (Block List, Max: 1) Location for S3 data store. (see [below for nested schema](#nestedblock--definition--include_location--location--s3_location))

<a id="nestedblock--definition--include_location--location--athena_location"></a>
### Nested Schema for `definition.include_location.location.s3_location`

Required:

- **catalog** (String) Catalog name.

Optional:

- **db** (String) Database name.
- **table** (String) Table name.


<a id="nestedblock--definition--include_location--location--mongo_location"></a>
### Nested Schema for `definition.include_location.location.s3_location`

Required:

- **db** (String) Database name.

Optional:

- **collection** (String) Collection name.


<a id="nestedblock--definition--include_location--location--mysql_location"></a>
### Nested Schema for `definition.include_location.location.s3_location`

Required:

- **db** (String) Database name.

Optional:

- **table** (String) Table name.


<a id="nestedblock--definition--include_location--location--relational_location"></a>
### Nested Schema for `definition.include_location.location.s3_location`

Required:

- **db** (String) Database name.

Optional:

- **schema** (String) Schema name.
- **table** (String) Table name.


<a id="nestedblock--definition--include_location--location--s3_location"></a>
### Nested Schema for `definition.include_location.location.s3_location`

Required:

- **bucket** (String) Bucket name.

Optional:

- **object_key** (String) Object Key name.



<a id="nestedblock--definition--include_location--location_parts_full"></a>
### Nested Schema for `definition.include_location.location_parts_full`

Required:

- **name** (String) The name of the location part.
- **type** (String) The type of the location part. Optional values: TABLE, COLUMN, SEMANTIC_MODEL, REPORT, DASHBOARD, DATABASE, SCHEMA, JSON_PATH, WAREHOUSE, ENDPOINT, TYPE, FIELD, EXTERNAL_LOCATION, CATALOG, BUCKET, OBJECT, COLLECTION, VIEW, etc




<a id="nestedblock--custom_policy"></a>
### Nested Schema for `custom_policy`

Optional:

- **priority** (Number) Dataset custom policy priority. Defaults to `100`.
- **rules_yaml** (String) Custom policy rules YAML.
- **tags_yaml** (String) Custom policy tags YAML.

~> **Note: The dataset resource is stateful:** The dataset resource is stateful, deletion or terraform resource name change should be avoided.