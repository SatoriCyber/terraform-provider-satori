---
layout: ""
page_title: "Provider: Satori"
description: |-
The Satori provider provides resources to interact with the Satori management API.
---

# Satori Provider

The Satori provider provides resources to interact with the Satori management API.

A service account must be created for Terraform.

## Example Usage

```terraform
provider "satori" {
  #can be provided via environment variable: SATORI_SA
  service_account = "522fb8ab-8d7b-4498-b39d-6911e2839253"
  #can be provided via environment variable: SATORI_SA_KEY
  service_account_key = "OZhw6ImBHXWMf51oICtfMoSYmm8gq9VxbYZTZjzaSO5NT0EHxbopnpLBuXQJo6aS"
  #satori account id for all resources in this terraform
  satori_account = "7cb42d6f-4d74-46c2-86c3-718116c1f5a1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **satori_account** (String) Your Satori account ID.

### Optional

- **service_account** (String) Service account ID with administrative privileges. Can be specified with the `SATORI_SA` environment variable.
- **service_account_key** (String, Sensitive) Service account key. Can be specified with the `SATORI_SA_KEY` environment variable.
- **url** (String) Defaults to `https://app.satoricyber.com`.
- **verify_tls** (Boolean) Defaults to `true`.