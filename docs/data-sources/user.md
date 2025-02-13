---
layout: ""
page_title: "satori_user (Data Source)"
description: |-
satori_user data source allows finding user ID by user email.
---

# satori_user (Data Source)

The **satori_user** data source allows finding user ID by user email.

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

- **email** (String) User's email address.

### Read-Only

- **id** (String) User's ID.