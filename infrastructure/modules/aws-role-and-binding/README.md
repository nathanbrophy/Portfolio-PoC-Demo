# aws-role-and-binding

- [aws-role-and-binding](#aws-role-and-binding)
  - [Usage](#usage)
  - [Inputs](#inputs)
  - [Outputs](#outputs)
  - [Authors](#authors)


This module allows you to define an AWS policy and role, and perform the attachment for the policy to the role, so that it can be assumed by a principle later on. 

## Usage

```hcl
module "iam_binding_and_policy" {
    source = "../modules/aws-role-and-binding"
    iam_policy_name = "example-policy-name"
    iam_role_name = "example-role-name"
    policy = <<EOT
{
    # ... JSON policy here
}
EOT

    description = "This is an example description"
    assume_role_policy = <<EOT
{
    # ... JSON policy here
}
EOT

    tags = {
        "example-tag-key": "example-tag-value",
    }
}
```

## Inputs

| Name | Purpose | Type | Required | Default |
| :--- | :------ | :--- | :------- | :------ |
| `assume_role_policy` | This is the assume role policy passed into the generated role, used in the attachment to allow principles to assume the role. | `string` | Yes | |
| `description` | This is the description for the generated IAM policy | `string` | Yes | |
| `policy` | This is the actual JSON encoded policy ruleset that will be used to defined the genreated IAM policy. | `string` | Yes | |
| `iam_policy_name` | This variable is used to declare the generated IAM policy name. | `string` | No | `ingress_policy` |
| `iam_role_name` | This variable is used to declare the generated IAM role name. | `string` | No | `ingress_alb_controller` |
| `tags` | This variable is used to pass additional tags to the generated IAM role. | `map(string)` | No | `{}` |

## Outputs

| Name | Type | Value |
| :--- | :--- | :---- |
| `role_name` | `string` | This value generated as the role name used to create the IAM role object. |

## Authors

* `@nathanbrophy`