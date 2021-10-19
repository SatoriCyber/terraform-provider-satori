resource "satori_security_policy" "security_policy" {
  name = "security policy terraform"
  profile {
    masking {
      active = true
      rule {
        id = "3"
        description = "rule 1"
        active = true
        criteria {
          condition = "IS"
          identity {
            type = "USER"
            name = "test-user"
          }
        }
        action {
          type = "APPLY_MASKING_PROFILE" # optional
          masking_profile_id = satori_masking_profile.masking_profile.id # as reference of previously created masking profile
        }
      }
      rule {
        id = "1"
        description = "rule 2"
        active = true
        criteria {
          condition = "IS_NOT"
          identity {
            type = "IDP_GROUP"
            name = "test-group"
          }
        }
        action {
          type = "APPLY_MASKING_PROFILE"
          masking_profile_id = satori_masking_profile.masking_profile.id
        }
      }
      rule {
        id = "2"
        description = "rule 3"
        active = true
        criteria {
          condition = "IS_NOT"
          identity {
            type = "GROUP"
            group_id = satori_directory_group.group2.id
          }
        }
        action {
          type = "APPLY_MASKING_PROFILE"
          masking_profile_id = satori_masking_profile.masking_profile.id
        }
      }
      rule {
        id = "4"
        description = "rule 4"
        active = true
        criteria {
          condition = "IS_NOT"
          identity {
            type = "GROUP"
            group_id = satori_directory_group.group3.id
          }
        }
        action {
          masking_profile_id = satori_masking_profile.masking_profile.id
        }
      }
    }
  }
}