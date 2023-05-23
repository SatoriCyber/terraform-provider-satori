resource "satori_user_settings" "settings_for_user_json_path" {
  user_id  = local.user_id_to_configure_settings_on

  // You may define a path to a json file containing a JSON object with the attributes for the user.
  attributes = file("${path.module}/attributes/user_a.json")
}


resource "satori_user_settings" "settings_for_user_raw_json" {
  user_id  = local.user_id_to_configure_settings_on

  // You may define the attributes in a raw JSON object using terraform's jsoncode({}).
  attributes = jsonencode({
    name      = "William"
    age       = 30.5
    cities    = ["Tel Aviv", "London", "Lisbon"]
    is_active  = true
    kids_ages = [1, 5, 6.5]
  })
}