# AWS

- [AWS](#aws)
  - [Provision Environment](#provision-environment)
  - [Teardown Environment](#teardown-environment)
  - [Linting](#linting)


This repository is a sample extension of the [Hashicorp Terraform Developer Tutorial](https://developer.hashicorp.com/terraform/tutorials/kubernetes/eks#eks) for creating an EKS cluster in AWS public cloud.  

## Provision Environment

The following steps can be taken to deploy the AWS environment:

```
% terraform plan 

# Review the above plan before applying

% terraform apply 
```

## Teardown Environment

The environment can be destroyed by using the following command:

```
% terraform destroy -auto-approve
```

## Linting

Terraform's builtin linter `terraform validate` is ran on the code from the `lint` command that is locally available to the working directory. 