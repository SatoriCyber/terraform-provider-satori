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
