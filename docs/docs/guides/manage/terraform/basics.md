---
title: Terraform Provider Basics
---

It covers how to:

- manage ZITADEL resources through the ZITADEL Terraform provider

Prerequisites:

- existing ZITADEL Instance, if not present follow [this guide](../../start/quickstart)
- existing user with enough authorization to manage the desired resources, if not present follow [this guide](../../integrate/serviceusers)
- installed Terraform, if not present follow [this guide](https://learn.hashicorp.com/tutorials/terraform/install-cli)

## Manage ZITADEL resources through terraform

The full documentation and examples are available [here](https://registry.terraform.io/providers/zitadel/zitadel/latest/docs).

To provide a small guide to where to start:

1. Create a folder where all the terraform files reside.
2. Configure the provider to use the right domain, port and token, with for example a `main.tf`file [as shown in the example](https://registry.terraform.io/providers/zitadel/zitadel/latest/docs).
3. Add a `zitadel_org` resource to the `main.tf` file, to create and manage a new organization in the instance, [as shown in the example](https://registry.terraform.io/providers/zitadel/zitadel/latest/docs/resources/org).
4. Add any resources to the organization in the `main.tf` file, [as example a human user](https://registry.terraform.io/providers/zitadel/zitadel/latest/docs/resources/human_user).
5. (Optional) Use Terraform in the directory with the command `terraform plan`, to see which resources would be created and how.
6. Apply the changes and start managing your resources with terraform with `terraform apply`.
7. (Optional) Delete your created resources with `terraform destroy` to clean-up.
