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
    owners = [ "522fb8ab-8d7b-4498-b39d-6911e2839253", data.satori_user.data_steward1.id ]

    include_location {
      datastore = "12345678-95cf-474f-a1d6-d5084810dd95"
    }

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
          db = "db2"
          schema = "schema1"
        }
      }
    }

    include_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location {
        relational_location {
          db = "db2"
          schema = "schema2"
          table = "table"
        }
      }
    }

    exclude_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location {
        relational_location {
          db = "db2"
          schema = "schema1"
          table = "tableX"
        }
      }
    }
  }

  access_control_settings {
    enable_access_control = false
    enable_user_requests = false
    enable_self_service = false
  }

  custom_policy {
    #default priority is 100
    #priority = 100
    rules_yaml = file("${path.module}/rules.yaml")
    tags_yaml = file("${path.module}/tags.yaml")
  }

  security_policies = [ "56412aff-6ecf-4eff-9b96-2e0f6ec36c42" ]
}

// Example with different location types
resource "satori_dataset" "dataset2" {
  definition {
    name = "satori_dataset terraform test"
    description = "from satori terraform provider"
    #the service account must also be an owner to be able to modify settings beyond definition
    owners = [ "522fb8ab-8d7b-4498-b39d-6911e2839253", data.satori_user.data_steward1.id ]

    include_location {
      datastore = "12345678-95cf-474f-a1d6-d5084810dd95"
    }

    include_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location {
        relational_location {
          db = "db1"
          schema = "schema1"
        }
      }
    }

    include_location {
      datastore = "3go33ff5-95cf-474f-a1d6-d5084810dd5k"
      location {
        mongo_location {
          db = "db1"
          collection = "collection1"
        }
      }
    }

    include_location {
      datastore = "8kl43ff5-95cf-474f-a1d6-d508481049lw"
      location {
        s3_location {
          bucket = "bucket1"
          object_key = "a/b/c"
        }
      }
    }

    exclude_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location {
        relational_location {
          db = "db1"
          schema = "schema1"
          table = "tableX"
        }
      }
    }
  }

  access_control_settings {
    enable_access_control = false
    enable_user_requests = false
    enable_self_service = false
  }

  custom_policy {
    #default priority is 100
    #priority = 100
    rules_yaml = file("${path.module}/rules.yaml")
    tags_yaml = file("${path.module}/tags.yaml")
  }

  security_policies = [ "56412aff-6ecf-4eff-9b96-2e0f6ec36c42" ]
}

// Example with deprecated usage of relational_location field
resource "satori_dataset" "dataset3" {
  definition {
    name = "satori_dataset terraform test"
    description = "from satori terraform provider"
    #the service account must also be an owner to be able to modify settings beyond definition
    owners = [ "522fb8ab-8d7b-4498-b39d-6911e2839253", data.satori_user.data_steward1.id ]

    include_location {
      datastore = "12345678-95cf-474f-a1d6-d5084810dd95"
    }

    include_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      relational_location {
        db = "db1"
      }
    }

    include_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      relational_location {
        db = "db2"
        schema = "schema1"
      }
    }

    include_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      relational_location {
        db = "db2"
        schema = "schema2"
        table = "table"
      }
    }

    exclude_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      relational_location {
        db = "db2"
        schema = "schema1"
        table = "tableX"
      }
    }
  }

  access_control_settings {
    enable_access_control = false
    enable_user_requests = false
    enable_self_service = false
  }

  custom_policy {
    #default priority is 100
    #priority = 100
    rules_yaml = file("${path.module}/rules.yaml")
    tags_yaml = file("${path.module}/tags.yaml")
  }

  security_policies = [ "56412aff-6ecf-4eff-9b96-2e0f6ec36c42" ]
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
- **enable_self_service** (Boolean, Deprecated) Allow users to grant themselves access to this dataset. Defaults to `false`. The 'enable_self_service' field has been deprecated. Please check the Dataset Permissions section in the documentation.
- **enable_user_requests** (Boolean, Deprecated) Allow users to request access to this dataset. Defaults to `false`. The 'enable_user_requests' field has been deprecated. Please check the Dataset Permissions section in the documentation.


<a id="nestedblock--definition"></a>
### Nested Schema for `definition`

Required:

- **name** (String) Dataset name.

Optional:

- **description** (String) Dataset description.
- **exclude_location** (Block List) Location to exclude from dataset. (see [below for nested schema](#nestedblock--definition--exclude_location))
- **include_location** (Block List) Location to include in dataset. (see [below for nested schema](#nestedblock--definition--include_location))
- **owners** (List of String) IDs of Satori users that will be set as dataset owners.

<a id="nestedblock--definition--exclude_location"></a>
### Nested Schema for `definition.exclude_location`

Required:

- **datastore** (String) Data store ID.

Optional:

- **location** (Block List, Max: 1) Location for a data store. Can include only one location type field from the above: relational_location, mysql_location, athena_location, mongo_location and s3_location . Conflicts with 'relational_location' field. (see [below for nested schema](#nestedblock--definition--exclude_location--location))
- **relational_location** (Block List, Max: 1, Deprecated) Location for a relational data store. Conflicts with 'location' field. The 'relational_location' field has been deprecated. Please use the 'location' field instead. (see [below for nested schema](#nestedblock--definition--exclude_location--relational_location))

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



<a id="nestedblock--definition--exclude_location--relational_location"></a>
### Nested Schema for `definition.exclude_location.relational_location`

Required:

- **db** (String) Database name.

Optional:

- **schema** (String) Schema name.
- **table** (String) Table name.



<a id="nestedblock--definition--include_location"></a>
### Nested Schema for `definition.include_location`

Required:

- **datastore** (String) Data store ID.

Optional:

- **location** (Block List, Max: 1) Location for a data store. Can include only one location type field from the above: relational_location, mysql_location, athena_location, mongo_location and s3_location . Conflicts with 'relational_location' field. (see [below for nested schema](#nestedblock--definition--include_location--location))
- **relational_location** (Block List, Max: 1, Deprecated) Location for a relational data store. Conflicts with 'location' field. The 'relational_location' field has been deprecated. Please use the 'location' field instead. (see [below for nested schema](#nestedblock--definition--include_location--relational_location))

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



<a id="nestedblock--definition--include_location--relational_location"></a>
### Nested Schema for `definition.include_location.relational_location`

Required:

- **db** (String) Database name.

Optional:

- **schema** (String) Schema name.
- **table** (String) Table name.




<a id="nestedblock--custom_policy"></a>
### Nested Schema for `custom_policy`

Optional:

- **priority** (Number) Dataset custom policy priority. Defaults to `100`.
- **rules_yaml** (String) Custom policy rules YAML.
- **tags_yaml** (String) Custom policy tags YAML.