# tf-aws-module_primitive-appmesh_virtual_service

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)

## Overview

Creates an [AWS App Mesh Virtual Service](https://docs.aws.amazon.com/app-mesh/latest/userguide/virtual_services.html), an abstraction of a real service that other resources in the mesh discover by name — a DNS-style hostname — without needing to know how that service is actually implemented.

Traffic destined for the virtual service is provided by either a single virtual node (`virtual_node_name`) or a virtual router (`virtual_router_name`); these are mutually exclusive and exactly one should be set for the virtual service to route traffic anywhere.

This module creates only the virtual service itself. The mesh, virtual node, and virtual router it references are separate primitives — see [tf-aws-module_primitive-appmesh](https://github.com/launchbynttdata/tf-aws-module_primitive-appmesh), [tf-aws-module_primitive-virtual_node](https://github.com/launchbynttdata/tf-aws-module_primitive-virtual_node), and [tf-aws-module_primitive-virtual_router](https://github.com/launchbynttdata/tf-aws-module_primitive-virtual_router).

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.0.2, < 2.0.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 5.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_appmesh_virtual_service.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/appmesh_virtual_service) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_app_mesh_name"></a> [app\_mesh\_name](#input\_app\_mesh\_name) | Id of the App Mesh where the virtual service will reside | `string` | n/a | yes |
| <a name="input_name"></a> [name](#input\_name) | Name of the virtual service | `string` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | A map of custom tags to be attached to this resource | `map(string)` | `{}` | no |
| <a name="input_virtual_node_name"></a> [virtual\_node\_name](#input\_virtual\_node\_name) | Name of the virtual Node to associate with the virtual service. Conflicts with virtual\_router | `string` | `""` | no |
| <a name="input_virtual_router_name"></a> [virtual\_router\_name](#input\_virtual\_router\_name) | Name of the virtual Router to associate with the virtual service. Conflicts with virtual\_node | `string` | `""` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_arn"></a> [arn](#output\_arn) | ARN of the virtual service |
| <a name="output_id"></a> [id](#output\_id) | ID of the virtual service |
| <a name="output_name"></a> [name](#output\_name) | Name of the virtual service |
<!-- END_TF_DOCS -->

## Module Development

### Pre-Requisites

The following commands should be available on your system:

- `asdf` or `mise`
- `make`
- `python3` (for pre-commit)

Additionally, your `git` user and email must be configured. Run the `make configure` command from the root of the repository to ensure that you meet these requirements.

### Pre-Commit hooks

The [.pre-commit-config.yaml](.pre-commit-config.yaml) file defines certain `pre-commit` hooks that are relevant to Terraform and Golang, as well as some common linting tasks. These will be configured for you when you run `make configure`.

### Local Validation

You should validate the changes you make to any module locally, prior to pushing your changes in a branch to GitHub.

1. Ensure that you have run `make configure` successfully.

2. Ensure you are signed into the appropriate cloud provider (e.g. AWS or Azure) for the module under test in your current console session.

3. Run the Terraform and Golang linters with the following command:

```
make lint
```

4. Once you have satisfied the linters, the following command will build example infrastructure in your configured cloud, run the tests, and then tear down the infrastructure it created:

```
make test
```

The pre-commit validations, as well as the `make lint` and `make test` targets, will all be performed in CI. Running these validations locally prior to opening a PR helps ensure a smooth review and merge process.

### Review & Merge Process

Once your change has been tested locally and your branch pushed up, open a new Pull Request for your branch to the default (main) branch of this repository.

The title of your Pull Request will determine the version bump for this change, and the title must be in [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#specification) format in order to merge. A breaking change will trigger a major version bump, a feature will trigger a minor version bump, and all other types will trigger a patch version bump.

Ensure your CI workflows are passing; seek approval from teammates and address any feedback; seek any explicit approvals required by the CODEOWNERS file. You may merge the PR as soon as all requirements are met, and a new release and tag will be automatically created for you.

### Automatic Updates

The shared configuration and workflow files in this repository are largely managed through the [launch-terraform-skeleton](https://github.com/launchbynttdata/launch-terraform-skeleton) repository. Outside of perhaps the `.gitignore` to account for specific files being generated by certain Terraform modules (e.g. Lambda functions), there should not be much cause to update these files on a per-repo basis, and making changes to them individually is discouraged.

If desired, you can check for and run these updates locally in a branch if you have the `copier` tool installed. Some example commands are included below:

```
# Check for updates, optionally checking prerelease versions
copier check-update [--prereleases]

# Run an update, using default answers if there are any. We use tasks, which requires --trust to be set.
copier update --defaults --trust [--prereleases]

# Recopy from the source, and --overwrite all templated files in the process
copier recopy --defaults --trust --overwrite [--prereleases]
```

Automatic updates will run through a scheduled workflow, and if the post-update tests are successful, the Pull Request created will automatically merge. Conflicts in the update or failures to test may leave a Pull Request outstanding, which needs to be addressed by a Launch Engineer.
