resource "satori_request_access_rule" "request_access1_dataset1" {
  //reference to owning dataset
  parent_data_policy = satori_dataset.dataset1.data_policy_id
  //granted access level, OWNER, READ_WRITE, READ_ONLY
  access_level = "OWNER"
  //identity can not be changed after creation
  identity {
    type = "USER"
    name = "test-user"
  }
  expire_in {
    unit_type = "MONTHS" //MINUTES, HOURS, DAYS, WEEKS, MONTHS, YEARS
    units = 3
  }
  revoke_if_not_used_in_days = 90
  require_approver_note = true //default is false

  // Optional to add approvers on an access-rule.
  approvers {
    type = "GROUP"
    id   = "781246b7-461d-493a-a2d6-86f2e5w01ff2"
  }

  approvers {
    type = "USER"
    id   = "78dc2cb7-461d-493a-a2d6-86e71fv4v5d2"
  }

  approvers {
    // The MANAGER approver type should not have `id` field set for it.
    type = "MANAGER"
  }
}

resource "satori_request_access_rule" "request_access2_dataset1" {
  parent_data_policy = satori_dataset.dataset1.data_policy_id
  access_level = "READ_ONLY"
  identity {
    type = "GROUP"
    group_id = satori_directory_group.group1.id
  }
  expire_in {
    unit_type = "DAYS" //MINUTES, HOURS, DAYS, WEEKS, MONTHS, YEARS
    units = 5
  }
  revoke_if_not_used_in_days = 90
  //dataset default security policies
  security_policies = [ ]
  require_approver_note = false // Optional, as default is false
}

resource "satori_request_access_rule" "request_access3_dataset1" {
  parent_data_policy = satori_dataset.dataset1.data_policy_id
  access_level = "READ_WRITE"
  identity {
    type = "IDP_GROUP"
    name = "groupName"
  }
  //no security policies
  security_policies = [ "none" ]
  // suspend this rule
  enabled = false
}

resource "satori_request_access_rule" "request_access1_dataset_definition1" {
  parent_data_policy = satori_dataset.dataset_definition1.data_policy_id
  access_level = "READ_ONLY"
  identity {
    type = "EVERYONE"
  }
  //specific security policies
  security_policies = [ "8c4745f5-a21e-4b7a-bb21-83c54351539f" ]
}