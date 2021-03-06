resource "satori_masking_profile" "masking_profile" {
  name = "masking profile1"

  condition {
    tag = "c12n.pii::address"
    type = "TRUNCATE"
    truncate = 1
  }

  condition {
    tag = "c12n.pii::address"
    type = "REPLACE_CHAR"
    replacement = "a"
  }

  condition {
    tag = satori_custom_taxonomy_classifier.cls2.tag # as reference of previously created custom classifier
    type = "EMAIL_PREFIX"
  }

  condition {
    tag = satori_custom_taxonomy_classifier.cls2.tag
    type = "REPLACE_CHAR"
    replacement = "a"
  }

}