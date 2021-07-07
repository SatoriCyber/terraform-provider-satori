terraform {
  required_providers {
    satori = {
      version = "1.0.0"
      source  = "satoricyber.com/edu/satori"
    }
  }
}

provider "satori" {
  #can be provided via environment variable: SATORI_SA_ID
  service_account_id = "522fb8ab-8d7b-4498-b39d-6911e2839253"
  #can be provided via environment variable: SATORI_SA_KEY
  service_account_key = "OZhw6ImBHXWMf51oICtfMoSYmm8gq9VxbYZTZjzaSO5NT0EHxbopnpLBuXQJo6aS"
  #satori account id for all resources in this terraform
  account_id = "7cb42d6f-4d74-46c2-86c3-718116c1f5a1"

  #for local env only:
  url = "https://app.satoricyber.example:8014"
  verify_tls = false
}

resource "satori_dataset" "test1" {
  definition {
    name = "terraform test 5"
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
      location = "bla.ccc"
    }
  }

  access_control_settings {
    enable_access_control = true
    enable_user_requests = true
    enable_self_service = true
  }

  custom_policy {
    #default priority is 100
    #priority =
    rules_yaml = file("${path.module}/rules.yaml")
    tags_yaml = file("${path.module}/tags.yaml")
  }

  security_policies = [ "56412aff-6ecf-4eff-9b96-2e0f6ec36c42" ]
}
