resource "satori_custom_taxonomy_category" "cat1" {
  name = "cat1"
  color = "#808080"
}

resource "satori_custom_taxonomy_category" "cat2" {
  name = "sub-cat2"
  color = "#000000"
  parent_category_id = satori_custom_taxonomy_category.cat1.id
}