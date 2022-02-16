---
layout: ""
page_title: "satori_datstore (Resource)"
description: |-
satori_datastore resource allows defining datastore(s)
---

# satori_datastore (Resource)

Satori provides the ability to connect to broad range of data stores repositories.
The **satori_datastore** resource allows lifecycle management for the datastores.

## Example Usage

```terraform
provider "satori" {
  #can be provided via environment variable: SATORI_SA
  service_account     = "aaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaaa"
  #can be provided via environment variable: SATORI_SA_KEY
  service_account_key = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
  #satori account id for all resources in this terraform
  satori_account      = "cccccc-cccc-cccc-cccc-cccccccccccc"
  url        = "https://<satori-mgmt-server-address>:8014"
  verify_tls = true
}
locals {
  dataaccess_controller_id = "<assigned dataaccess_controller_id>"
}
resource "satori_datastore" "datastore0" {
  name                     = "exampleDatastore"
  hostname                 = "data.source.target.hostname"
  dataaccess_controller_id = local.dataaccess_controller_id
  type                     = "SNOWFLAKE"
  originPort                     = 8081
  identity_provider_id     = "aaaaaaaaaaaaa-ddddd-ddddddddd-dddddddd"

}
# output of generated id for newly created datastore
output "datastore_created_id" {
  value = satori_datastore.datastore0.id
}
```
```terraform
provider "satori" {
  #can be provided via environment variable: SATORI_SA
  service_account     = "aaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaaa"
  #can be provided via environment variable: SATORI_SA_KEY
  service_account_key = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
  #satori account id for all resources in this terraform
  satori_account      = "cccccc-cccc-cccc-cccc-cccccccccccc"
  url        = "https://<satori-mgmt-server-address>:8014"
  verify_tls = false
}
locals {
  dataaccess_controller_id = "<assigned dataaccess_controller_id>"
}

resource "satori_datastore" "datastore0" {
  name = "exampleDatastore"

  dataaccess_controller_id = local.dataaccess_controller_id
  # data source specific connection settings
  type                     = "BIGQUERY"
  project_ids              = [111, 222] #  BigQuery affected project ids
  hostname                 = "data source target hostname"
  origin_port              = 8081 # data source server's ip
  identity_provider_id     = "aaaaaaaaaaaaa-ddddd-ddddddddd-dddddddd"
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
        identity = "user1"
      }
      excluded_identities {
        identity_type = "USER"
        identity = "user2"
      }
      excluded_query_patterns {
        pattern = ".*a.*"
      }
      excluded_query_patterns {
        pattern = ".*b.*"
      }
    }
  }
}


output "datastore_created_id" {
  value = satori_datastore.datastore0.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **baseline_security_policy** (Block List, Min: 1, Max: 1) Baseline security policy. (see [below for nested schema](#nestedblock--baseline_security_policy))
- **dataaccess_controller_id** (String) Host FQDN name.
- **hostname** (String) Data provider's FQDN hostname.
- **name** (String) DataStore name.
- **type** (String) IDs of Satori users that will be set as DataStore owners.

### Optional

- **custom_ingress_port** (Number) Port number description.
- **identity_provider_id** (String) IDs of Satori users that will be set as DataStore owners.
- **origin_port** (Number) Port number description.
- **project_ids** (Set of String) ProjectIds list of project IDs

### Read-Only

- **id** (String) DataStore resource id.
- **parent_id** (String) Parent resource id.

<a id="nestedblock--baseline_security_policy"></a>
### Nested Schema for `baseline_security_policy`

Required:

- **exclusions** (Block List, Min: 1, Max: 1) Exempt users and patterns from baseline security policy (see [below for nested schema](#nestedblock--baseline_security_policy--exclusions))
- **unassociated_queries_category** (Block List, Min: 1, Max: 1) UnassociatedQueriesCategory (see [below for nested schema](#nestedblock--baseline_security_policy--unassociated_queries_category))
- **unsupported_queries_category** (Block List, Min: 1, Max: 1) UnsupportedQueriesCategory (see [below for nested schema](#nestedblock--baseline_security_policy--unsupported_queries_category))

Optional:

- **type** (String) DataStore basepolicy . Defaults to `BASELINE_POLICY`.

<a id="nestedblock--baseline_security_policy--exclusions"></a>
### Nested Schema for `baseline_security_policy.exclusions`

Optional:

- **excluded_identities** (Block List) Exempt Users from the Baseline Security Policy (see [below for nested schema](#nestedblock--baseline_security_policy--exclusions--excluded_identities))
- **excluded_query_patterns** (Block List) Exempt Queries from the Baseline Security Policy (see [below for nested schema](#nestedblock--baseline_security_policy--exclusions--excluded_query_patterns))

<a id="nestedblock--baseline_security_policy--exclusions--excluded_identities"></a>
### Nested Schema for `baseline_security_policy.exclusions.excluded_identities`

Optional:

- **identity** (String) Username
- **identity_type** (String) USER type are supported


<a id="nestedblock--baseline_security_policy--exclusions--excluded_query_patterns"></a>
### Nested Schema for `baseline_security_policy.exclusions.excluded_query_patterns`

Optional:

- **pattern** (String) Query pattern



<a id="nestedblock--baseline_security_policy--unassociated_queries_category"></a>
### Nested Schema for `baseline_security_policy.unassociated_queries_category`

Optional:

- **query_action** (String) Default policy action for querying locations that are not associated with a dataset.


<a id="nestedblock--baseline_security_policy--unsupported_queries_category"></a>
### Nested Schema for `baseline_security_policy.unsupported_queries_category`

Required:

- **query_action** (String) Default policy action for unsupported queries and objects