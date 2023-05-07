---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "satori_user_settings Resource - terraform-provider-satori"
subcategory: ""
description: |-
  Satori user settings allows to config existing user's configuration. Currently supports only user's attributes configuration
---

# satori_user_settings (Resource)

Satori user settings allows to config existing user's configuration. Currently supports only user's attributes configuration



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **attributes** (String) User's set of attributes in JSON object format. may include the following types: int, string, float, boolean, string[], number[], where number may be float/int. The value may be a path to a json file that contains the attributes for a user or a raw JSON string, for example: "./attribute_files/user_a.json" OR "{"company": "SatoriCyber","age": 30.5,"cities": ["Washington", "Lisbon"],"kids_age": [1, 3.14759, 7], "isActive": true}"
- **user_id** (String) User ID to manage settings for.

### Optional

- **id** (String) The ID of this resource.

