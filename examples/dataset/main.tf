resource "satori_dataset" "dataset1" {
  definition {
    name = "terraform test"
    description = "from satori terraform provider"
    owners_ids = [ "522fb8ab-8d7b-4498-b39d-6911e2839253" ]

    include_location {
      datastore_id = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location = "bla"
    }

    include_location {
      datastore_id = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location = "bla.bla"
    }

    exclude_location {
      datastore_id = "80f33ff5-95cf-474f-a1d6-d5084810dd95"
      location = "bla.ddd"
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
