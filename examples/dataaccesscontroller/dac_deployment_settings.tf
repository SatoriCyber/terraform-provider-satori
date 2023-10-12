locals {
  dataaccess_controller_id = "<assigned dataaccess_controller_id>"
}

data "satori_dac_deployment_settings" "deployment_settings" {
  id = local.dataaccess_controller_id
}

# Value can be used as following:
# satori_dac_deployment_settings.deployment_settings.service_account
