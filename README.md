<br>
<a href="https://terraform.io">
    <img src=".github/tf.png" alt="Terraform logo" title="Terraform" height="50" />
</a>
&nbsp;&nbsp;
<a href="https://www.satoricyber.com/">
    <img src="https://avatars.githubusercontent.com/u/59790990" alt="Satori logo" title="Satori" height="50" />
</a>

<!-- TOC -->
* [Terraform Provider for Satori](#terraform-provider-for-satori)
  * [Local development](#local-development)
    * [First time setup:](#first-time-setup)
    * [Run the following command to build the provider](#run-the-following-command-to-build-the-provider)
    * [Generate/update documentation](#generateupdate-documentation)
    * [Test sample configuration](#test-sample-configuration)
    * [Local vulnerabilities check](#local-vulnerabilities-check-)
  * [Create a Configuration File and State for Existing Resources](#create-a-configuration-file-and-state-for-existing-resources)
    * [Motivation:](#motivation)
    * [Prerequisite:](#prerequisite-)
    * [Step 1: Generate a new `main.tf` file with imported resources’ configuration](#step-1-generate-a-new-maintf-file-with-imported-resources-configuration)
    * [Step 2: Apply the import command in-order to sync the terraform’s state](#step-2-apply-the-import-command-in-order-to-sync-the-terraforms-state)
    * [Step 3: Setting up the files names to work against the new generated configuration file](#step-3-setting-up-the-files-names-to-work-against-the-new-generated-configuration-file)
      * [Option 1: Updating the current `main.tf` to the proper content](#option-1-updating-the-current-maintf-to-the-proper-content)
      * [Option 2: Updating `main.tf` to be an import file for future reference](#option-2-updating-maintf-to-be-an-import-file-for-future-reference)
    * [Step 4: Validate the configuration file](#step-4-validate-the-configuration-file)
<!-- TOC -->

# Terraform Provider for Satori

## Local development

### First time setup:
```shell
make init
```

### Run the following command to build the provider

```shell
make build
```

### Generate/update documentation

Do not edit files under `docs`, they are generated from `templates` and the source code.
To preview how the docs will look in the terraform registry, paste them here https://registry.terraform.io/tools/doc-preview

***Important:*** Run this command before commiting changes to git, to update the docs for recent changes.

```shell
make docs
```

### Test sample configuration

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

### Local vulnerabilities check 

Run the following command to check the vulnerabilities of the provider locally.
`govulncheck` is used to check for vulnerabilities in the Go code of the provider.

```shell
make vuln
```

## Create a Configuration File and State for Existing Resources

### Motivation:
In-order to have up to date resource state and configuration file on resources that were not managed by terraform, the client will be required to define which resources he would like to manage by terraform and import them along with updating the Terraform’s state.

This document will instruct the user how to conduct this procedure.

These instructions will have 3 main parts:
Generating a new `main.tf` file that will include all of the desired imported resources’ configuration.
Applying the import command in-order to sync the terraform’s state.
Setting up the files names to work against the new generated configuration file.

### Prerequisite: 
Terraform version 1.5.x+ is required, as this feature is not supported on older versions.


### Step 1: Generate a new `main.tf` file with imported resources’ configuration
*Note: If you want, you can define a `satori_imports` module and use it to import the resources. This will allow you to keep the main.tf file clean and organized, and would require you to execute in step 3 only option 1. So it will save some work & leave the imports module as it was for future reference.*

1. Create a new `main.tf` with the Satori provider information.
More info can be found here https://registry.terraform.io/providers/SatoriCyber/satori/latest.

    Initialize the terraform directory with the following command:

    ```shell
    terraform init
    ```

2. For each resource you want to manage using terraform add the terraform import resource:

    ```terraform
    import {
    to = ${terraform_resource_type}.${terraform_resource_unique_name_per_type}
    id = "${satori_resource_id}"
    }
    ```

    For example: 
    
    ```terraform
    import {
     to = satori_dataset.my_imported_dataset
     id = "fvf9e3a7-5bcd-4b7f-9230-3e5f5wce0b85"
    } 
    ```
    
    In the above example, we would like to import a `satori_dataset` resource type, and its terraform resource name will be `my_imported_dataset`. To this resource we would like to import a satori dataset resource which has the ID `fvf9e3a7-5bcd-4b7f-9230-3e5f5wce0b85`.

3. Once you’ve configured for each resource that you want to import the proper terraform import resource which has the proper details. Run the following command:

    ```shell
    terraform plan -generate-config-out=new_main.tf
    ``` 
   The command above will create under the working directory a new file, `new_main.tf` that will contain all the resources’ configurations of the `import resource` you’ve defined above.

### Step 2: Apply the import command in-order to sync the terraform’s state

Now after creating the terraform configuration file, we would like to synchronize the terraform’s state. In-order to do so, run the following command:

```terraform
terraform apply
```

The command above will update the terraform state to match the configuration file generated by the previous command. (As long that no changes occurred on the imported resource between the generation of the configuration file and the above command)

### Step 3: Setting up the files names to work against the new generated configuration file

This part has 2 options available for it:

#### Option 1: Updating the current `main.tf` to the proper content

1. From the `main.tf` file remove all of the terraform import resources you’ve generated in step 2 of Part 1.
2. From the `new_main.tf` copy the entire content and paste it into the end of `main.tf` file.

#### Option 2: Updating `main.tf` to be an import file for future reference

1. From `main.tf` copy all the details for the provider to the `new_main.tf`:
    For example:

    ```terraform
    terraform { 
        required_providers {
            satori = {
                source  = "Satoricyber/satori"
                version = "version-your-using"
            }
        }
    }

    provider "satori" {
        service_account         = "service-account-id"
        service_account_key     = "service-account-key"
        satori_account          = "satori-account-id"
    }
    ```

2. Change the name of the current `main.tf` file to `imports_main.tf` (or any other name that is not `main.tf`)
3. Change the name of `new_main.tf` to `main.tf`


### Step 4: Validate the configuration file

Now, when running the `terraform plan` command, you should see the imported resources’ configuration file and the state should be in sync.

```shell
terraform plan
```

meaning that the configuration file, the state, and the resources in our system are in sync.

`From this point forward you are able to manage those resources using the terraform command`

Happy terraforming… 😀
