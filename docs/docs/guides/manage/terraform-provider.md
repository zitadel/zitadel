---
title: ZITADEL Terraform Provider
sidebar_label: Terraform Provider
---

The [ZITADEL Terraform Provider](https://registry.terraform.io/providers/zitadel/zitadel/latest/docs) is a tool that allows you to manage ZITADEL resources through Terraform.
In other words, it lets you define and provision infrastructure for ZITADEL using Terraform configuration files.

This Terraform provider acts as a bridge, allowing you to manage various aspects of your ZITADEL instance directly through the [ZITADEL API](/docs/apis/introduction), using Terraform's declarative configuration language.
It can be used to create, update, and delete ZITADEL resources, as well as to manage the relationships between those resources.

## Before you start

Make sure you create the following resources in ZITADEL and have [Terraform installed](https://learn.hashicorp.com/tutorials/terraform/install-cli):

- [A ZITADEL Instance](../start/quickstart)
- [A service user](/docs/guides/integrate/service-users/authenticate-service-users) with [enough authorization](/docs/guides/manage/console/managers) to manage the desired resources

## Manage ZITADEL resources through terraform

The full documentation and examples are available on the [Terraform registry](https://registry.terraform.io/providers/zitadel/zitadel/latest/docs).

To provide a small guide to where to start:

1. Create a folder where all the terraform files reside.
2. Configure the provider to use the right domain, port and token, with for example a `main.tf`file [as shown in the example](https://registry.terraform.io/providers/zitadel/zitadel/latest/docs).
3. Add a `zitadel_org` resource to the `main.tf` file, to create and manage a new organization in the instance, [as shown in the example](https://registry.terraform.io/providers/zitadel/zitadel/latest/docs/resources/org).
4. Add any resources to the organization in the `main.tf` file, [as example a human user](https://registry.terraform.io/providers/zitadel/zitadel/latest/docs/resources/human_user).
5. (Optional) Use Terraform in the directory with the command `terraform plan`, to see which resources would be created and how.
6. Apply the changes and start managing your resources with terraform with `terraform apply`.
7. (Optional) Delete your created resources with `terraform destroy` to clean-up.

## References

- [Deploy ZITADEL in your infrastructure](/docs/self-hosting/deploy/overview)
- [ZITADEL CLI](/docs/self-hosting/manage/cli/overview)
- [Configuration Options in ZITADEL](/docs/self-hosting/manage/configure)
