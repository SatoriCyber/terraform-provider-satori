resource "satori_custom_taxonomy_category" "cat1" {
  name = "cat1"
  color = "#808080"
}

resource "satori_custom_taxonomy_category" "cat2" {
  name = "sub-cat2"
  description = "category from terraform"
  color = "#000000"
  parent_category = satori_custom_taxonomy_category.cat1.id
}

resource "satori_custom_taxonomy_classifier" "cls1" {
  name = "cls1"
  description = "classifier from terraform"
  type = "NON_AUTOMATIC"
  parent_category = satori_custom_taxonomy_category.cat1.id
  additional_satori_categories = ["pii"]
}

resource "satori_custom_taxonomy_classifier" "cls2" {
  name = "cls2"
  description = "classifier from terraform"
  type = "SATORI_BASED"
  parent_category = satori_custom_taxonomy_category.cat1.id
  satori_based_config {
    satori_base_classifier = "EMAIL"
  }
  scope {
    datasets = [satori_dataset.dataset1.id]
  }
}

resource "satori_custom_taxonomy_classifier" "cls3" {
  name = "cls3"
  type = "CUSTOM"
  parent_category = satori_custom_taxonomy_category.cat1.id
  custom_config {
    field_name_pattern = "abc.*xyz"
    field_type = "ANY"
  }
}
