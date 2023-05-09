resource "satori_user_settings" "settings_for_user_json_path" {
  user_id  = local.user_id_to_configure_settings_on

  // You may define a path to a json file containing a JSON object with the attributes for the user.
  attributes = "./attributes/user_a.json"
}


resource "satori_user_settings" "settings_for_user_raw_json" {
  user_id  = local.user_id_to_configure_settings_on

  // You may define the attributes in a raw JSON object string.
  attributes = "{\"name\": \"Williham\", \"age\": [3.5, 5], \"cities\": [\"Tel Aviv\", \"Lisbon\"]}"
}