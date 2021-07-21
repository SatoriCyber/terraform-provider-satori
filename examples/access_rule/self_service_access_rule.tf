resource "satori_self_service_access_rule" "perm1_dataset1" {
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
}

resource "satori_self_service_access_rule" "perm2_dataset1" {
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
}

resource "satori_self_service_access_rule" "perm3_dataset1" {
  parent_data_policy = satori_dataset.dataset1.data_policy_id
  access_level = "READ_WRITE"
  identity {
    type = "IDP_GROUP"
    name = "groupName"
  }
  //no security policies
  security_policies = [ "none" ]
}

resource "satori_self_service_access_rule" "perm1_dataset_definition1" {
  parent_data_policy = satori_dataset.dataset_definition1.data_policy_id
  access_level = "READ_ONLY"
  identity {
    type = "EVERYONE"
  }
  //specific security policies
  security_policies = [ "8c4745f5-a21e-4b7a-bb21-83c54351539f" ]
}