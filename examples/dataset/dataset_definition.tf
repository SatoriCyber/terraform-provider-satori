//get ID of user by email
data "satori_user" "data_steward1" {
  email = "data-steward@acme.organization"
}

resource "satori_dataset_definition" "dataset_definition1" {
  definition {
    name        = "satori_dataset_definition terraform test"
    description = "from satori terraform provider"
    owners = ["12345678-8d7b-4498-b39d-6911e2839253", data.satori_user.data_steward1.id]

    include_location {
      datastore = "12345678-95cf-474f-a1d6-d5084810dd95"
    }

    include_location {
      datastore = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location_parts = ["db1", "schema1"]
    }

    include_location {
      datastore     = "3go33ff5-95cf-474f-a1d6-d5084810dd5k"
      location_path = "db1.collection1"
    }

    include_location {
      datastore = "8kl43ff5-95cf-474f-a1d6-d508481049lw"
      location_parts = ["bucket1", "a/b/c"]
    }
  }
}

// Example with deprecated usage of location field (use location_path, location_parts or location_parts_full instead)
resource "satori_dataset_definition" "dataset_definition2" {
  definition {
    name        = "satori_dataset_definition terraform test"
    description = "from satori terraform provider"
    owners = ["12345678-8d7b-4498-b39d-6911e2839253", data.satori_user.data_steward1.id]

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
      location {
        relational_location {
          db     = "db2"
          schema = "schema1"
          table  = "tableX"
        }
      }
    }
  }
}
