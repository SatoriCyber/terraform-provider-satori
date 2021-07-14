<br>
<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terraform" height="50" />
</a>
&nbsp;&nbsp;
<a href="https://www.satoricyber.com/">
    <img src="https://avatars.githubusercontent.com/u/59790990" alt="Satori logo" title="Satori" height="50" />
</a>

# Terraform Provider for Satori

#### First time setup:
```shell
make init
```

#### Run the following command to build the provider

```shell
make build
```

#### Generate/update documentation

Do not edit files under `docs`, they are generated from `templates` and the source code.
To preview how the docs will look in the terraform registry, paste them here https://registry.terraform.io/tools/doc-preview

***Important:*** Run this command before commiting changes to git, to update the docs for recent changes.

```shell
make docs
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```

For local binary, configure terraform as:
```terraform
terraform {
  required_providers {
    satori = {
      version = "~>1.0.0"
      source  = "satoricyber.com/terraform/satori"
    }
  }
}
```