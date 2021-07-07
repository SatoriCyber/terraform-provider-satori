<br>
<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terraform" height="50" />
</a>
&nbsp;&nbsp;
<a href="https://www.satoricyber.com/">
    <img src="https://avatars.githubusercontent.com/u/59790990" alt="Satori logo" title="Satori" height="50" />
</a>

# Terraform Provider for Satori

Run the following command to build the provider

```shell
go build
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