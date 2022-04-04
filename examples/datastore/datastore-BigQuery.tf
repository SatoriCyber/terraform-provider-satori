locals {
  dataaccess_controller_id = "<assigned dataaccess_controller_id>"
}

resource "satori_datastore" "datastore0" {
  name = "exampleDatastore"

  dataaccess_controller_id = local.dataaccess_controller_id
  # data source specific connection settings
  type                     = "BIGQUERY"
  project_ids              = ["abc", "cdf"] #  BigQuery affected project ids
  hostname                 = "data source target hostname"
  origin_port              = 8081 # data source server's ip
  ####### BASELINE_POLICY SETTINGS #########
  baseline_security_policy {

    unassociated_queries_category {
      query_action = "REDACT" #Allowed: PASS┃REDACT┃BLOCK
    }
    unsupported_queries_category {
      query_action = "REDACT"
    }
    exclusions {

      excluded_identities {
        identity_type = "USER"
        identity      = "user1"
      }
      excluded_identities {
        identity_type = "USER"
        identity      = "user2"
      }
      excluded_query_patterns {
        pattern = ".*a.*"
      }
      excluded_query_patterns {
        pattern = ".*b.*"
      }
    }
  }
  network_policy {
    allowed_rules {
      note = "desc1"
      ip_ranges {
        ip_range = "3.2.3.1"
      }
      ip_ranges {
        ip_range = "3.2.3.1"
      }
    }
    allowed_rules {
      note = "desc2"
      ip_ranges {
        ip_range = "3.2.3.1"
      }
    }
    blocked_rules {
      note = "desc1"
      ip_ranges {
        ip_range = "3.2.3.1"
      }
      ip_ranges {
        ip_range = "3.2.3.1"
      }
    }
    blocked_rules {
      note = "desc2"
      ip_ranges {
        ip_range = "3.2.3.1"
      }
    }
}


output "datastore_created_id" {
  value = satori_datastore.datastore0.id
}
