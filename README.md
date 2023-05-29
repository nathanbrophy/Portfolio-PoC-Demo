# Portfolio Demo

- [Portfolio Demo](#portfolio-demo)
  - [Overview](#overview)
    - [Infrastructure](#infrastructure)
    - [K8s](#k8s)
    - [Server](#server)
  - [Architecture Diagrams](#architecture-diagrams)
  - [Running the PoC](#running-the-poc)
    - [Prerequisites:](#prerequisites)
      - [Tools](#tools)
      - [Env Vars](#env-vars)
    - [What is Installed](#what-is-installed)
    - [Verifying the Install](#verifying-the-install)
  - [Tearing Down the PoC](#tearing-down-the-poc)
  - [Testing](#testing)
    - [Infrastructure](#infrastructure-1)
    - [K8s](#k8s-1)
    - [Server](#server-1)

## Overview

This repository is meant to be a demonstration / PoC repository for a sample DevOps project.  The parts of the repository are as follows:

1. Infrastructure
2. K8s
3. Server

### Infrastructure

Infrastructure contains the IaC code (terraform) needed to provision the k8s environment in the public cloud of your choosing (currently only AWS is supported). 

### K8s

K8s contains the manfiests and codes needed to deploy the server to the k8s environment. 

### Server

Server contains the code needed to standup an example REST server container that can be deployed to a running kubernetes environment. 

## Architecture Diagrams

![High Level Arch Diagram for the PoC Repository](images/highlevel.png)

## Running the PoC 

Please read [the prerequisites](#prerequisites) before installing the application. 

![](./images/poc_run.gif)

To run the PoC sample the [deploy_config.yaml](./deploy_config.yaml) __must__ be configured to provide script configuration options.  Once the configuration is completed, the [deploy_poc.sh](./deploy_poc.sh) script can be ran to fully stand the PoC up on a public cloud environment. 

Before running the PoC please ensure the following credential steps are taken:

1. `docker login` performed to ensure the images can be pushed to their repositories
2. `aws configure` performed to ensure access to the configured S3 bucket to store the TF state (optional)

### Prerequisites:

#### Tools

| Requisite | Version (tested) | Reason |
| :-------- | :--------------- | :----- |
| `docker` | `20.10.14` | To build and push container images on change, Docker is needed.  If running PoC with `.deploy.no_build` set, this can be ignored. |
| `go` | `1.20.4` | Required to build binaries from source and perform testing, if running PoC with `.deploy.no_build` set, this can be ignored. |
| `kubectl` | `1.27` | Required to interact with the cluster and deploy resources. |
| `terraform` | `1.4.6` | Required to provision the AWS environment in the public cloud. |
| `yq` | `4.34.1` | Required to read configuration information for the PoC script. |
| `aws` | `2.11.23` | Required to get and set cluster credential information to deploy the application. |
| `eksctl` | `0.143.0` | Required to enable AWS load balancing. |

#### Env Vars

| Requisite | Value |
| :-------- | :---- |
| `AWS_ACCESS_KEY_ID` | The API key ID required for the terraform module to access the AWS environment |
| `AWS_SECRET_ACCESS_KEY` | The API key secret value required for the terraform module to access the AWS environment |

### What is Installed

When running the PoC the following items will be installed:

1. A k8s cluster installed to the public cloud
   1. Node for running the code
      1. Node group for tying the nodes to the cluster
   2. VPC created to allo networking for the cluster
   3. 4 subnets
      1. 2 private
      2. 2 public
   4. Security IAM groups and bindings allowing resource connectivity 
   5. TLS certificates for the cluster connection
   6. Cloudwatch log groups
2. An operator controller image that is built and pushed to a docker repository with anonymous pull 
   1. Can be skipped if `deploy.no_build` is set to a truthy value in the config file
3. An operator controller image is deployed to the running environment
   1. k8s manifests:
      1. `ServiceAccount`
      2. `Deployment`
         1. `ReplicaSet`
            1. `Pod`
      3. `ClusterRole`
      4. `ClusterRoleBinding`
      5. `CRD`
4. A REST server container image that is built and pushed to a docker repository with anonymous pull 
   1. Can be skipped if `deploy.no_build` is set to a truthy value in the config file
5. A REST server image is deployed to the running environment
   1. k8s manifests:
      1. `Deployment`
      2. `ServiceAccount`
      3. `Service`
      4. `Application`
         1. The resources above are generated by the controller based on the definition defined in the application CR
6. If the load balancer is installed with the PoC script (`.deploy.load_balancer.deploy`), then the following resources are additional configured:
   1. AWS:
      1. AWSLoadBalancerControllerIAMPolicy IAM policy
      2. A dynamic load balancer resource is created from the ingress
   2. K8s:
      1. Load balancer IAM Service Account binding through OpenID connect policy bindings
      2. Jetstack Cert Manager Operator Controller
      3. AWS ELB load balancer controller

### Verifying the Install

The provided [verify_e2e.sh](./verify_e2e.sh) script can be used to verify the deployed environment. The load balancer install option is optional, and in the absence of the load balancer, the verify script will port forward the Pod's running port and verify from there.

## Tearing Down the PoC

To tear down the PoC the following can be ran:

```sh
# The application must be deleted to cascade the deletion of cluster resources
# the load balancer in the AWS env must be deleted or it can cause 
# the subnet finalizers to not resolve properly resulting in orphaned resources. 
#
# This kubectl command will delete the application from the cluster.
% kubectl delete Application application-sample

# The following commands will remove the infrastructure
# provisioned in the PoC
% cd infrastructure/aws

% terraform destroy -auto-approve

# To cleanup the IAM permissions from the load balancer
# the following can be ran
% aws iam delete-policy --policy-arn=$(aws iam list-policies --no-cli-pager --output yaml | yq -Mr '.Policies[] | select(.PolicyName == "AWSLoadBalancerControllerIAMPolicy") | .Arn')
```

The above terraform command will completely remove all created public cloud resources, which will destroy the running applications as well. 

## Testing

Test practices and procedures are documented in each sections's own `README.md`, however and overview of this will be given here.

### Infrastructure

[Full Details](./infrastructure/aws/README.md)

The testing for this is handled via linter enforcement of terraform best practices, and can be run anytime with the [package lint script](./infrastructure/aws/lint).

### K8s

[Full Details](./k8s/README.md)

The testing for this is handled via unit tests for the controller package to codify the deployment pipeline, and the package provides an e2e test suite as well for ensuring reconciliation properly handles. 

To run the unit tests for the package perform the following:

```sh
% cd k8s

% make test
```

To run the e2e test suite for the package perform the following:

> NOTE: `minikube` is a prerequisite for this test, as we stand up an ephemeral, local k8s cluster to perform tests on, as to not potentially conflict with a production environment. 

```sh
% cd k8s/test_e2e.sh

% ./test_e2e.sh
```

### Server

[Full Details](./server/README.md)

The testing for this is handled through golang unit testing and an e2e test provided for the router handling of the REST server.  

To run the unit and e2e tests for the package perform the following:

```
% cd server

% make test
```
