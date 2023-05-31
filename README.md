# Portfolio Demo

- [Portfolio Demo](#portfolio-demo)
  - [Overview](#overview)
    - [Infrastructure](#infrastructure)
    - [K8s](#k8s)
    - [Server](#server)
  - [Architecture Diagrams](#architecture-diagrams)
  - [Running the PoC](#running-the-poc)
    - [Prerequisites:](#prerequisites)
      - [Running the PoC](#running-the-poc-1)
      - [Building the PoC](#building-the-poc)
    - [What is Installed](#what-is-installed)
      - [(1) Infrastructure / K8s Env (pipeline)](#1-infrastructure--k8s-env-pipeline)
      - [(2) Configure Cluster Credentials](#2-configure-cluster-credentials)
      - [(3) Setup / Deploy PoC (k8s)](#3-setup--deploy-poc-k8s)
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

[Infrastructure](./infrastructure/) contains the IaC code (terraform) needed to provision the k8s environment in the public cloud of your choosing (currently only AWS is supported). 

### K8s

[K8s](./k8s/) contains the manfiests and codes needed to deploy the server to the k8s environment. 

### Server

[Server](./server/) contains the code needed to standup an example REST server container that can be deployed to a running kubernetes environment. 

## Architecture Diagrams

![High Level Arch Diagram for the PoC Repository](images/highlevel.png)

## Running the PoC 

Please read [the prerequisites](#prerequisites) before installing the application. 

![Proof of Concept GIF to show an example run](./images/poc_run.gif)

To run the PoC sample, a Helm chart is provided to assist in installation.  The following can be ran to completely install the PoC to the environment:

```sh
# Configure the kube-config to connect to the deployed AWS EKS cluster
$ aws eks --region "${region}" update-kubeconfig --name "${cluster_name}"

# Install the actual PoC
# This installs the following:
#   1. CRDs for the controller
#   2. The controller itself
#   3. A sample CR to deploy the application
$ helm upgrade pocsample operator-controller \
   -i \
   --namespace acme-controller-manager \
   --create-namespace
```

### Prerequisites:

#### Running the PoC

##### Tools

| Requisite | Version (tested) | Reason |
| :-------- | :--------------- | :----- |
| `aws` | `2.11.23` | Required to get and set cluster credential information to deploy the application. |
| `kubectl` | `1.27` | Required to interact with the cluster and deploy resources. |
| `helm` | `v3.12.0` | Required to deploy the manifests to the cluster for the PoC. |

##### Env Vars

| Requisite | Value |
| :-------- | :---- |
| `AWS_ACCESS_KEY_ID` | The API key ID required for the terraform module to access the AWS environment |
| `AWS_SECRET_ACCESS_KEY` | The API key secret value required for the terraform module to access the AWS environment |

#### Building the PoC

| Requisite | Version (tested) | Reason |
| :-------- | :--------------- | :----- |
| `docker` | `20.10.14` | To build and push container images on change, Docker is needed.  If running PoC with `.deploy.no_build` set, this can be ignored. |
| `go` | `1.20.4` | Required to build binaries from source and perform testing, if running PoC with `.deploy.no_build` set, this can be ignored. |
| `minikube` | `v1.30.1` | This is only required for local development testing of the k8s controller. |
| `terraform` | `1.4.6` | Required to provision the AWS environment in the public cloud. |

##### Tools

##### Env Vars

| Requisite | Value |
| :-------- | :---- |
| `AWS_ACCESS_KEY_ID` | The API key ID required for the terraform module to access the AWS environment |
| `AWS_SECRET_ACCESS_KEY` | The API key secret value required for the terraform module to access the AWS environment |

### What is Installed

When running the PoC there are logical steps that are taken in the installation script.  This section will describe each logical step, what actions they perform, how to enable/disable them, and what resources are deployed as part of the step being ran. 

#### (1) Infrastructure / K8s Env (pipeline)

##### Description

This step is responsible for deploying the public cloud infrastructure.  Currently, only AWS is supported as a public cloud vendor, but the [infrastructure](./infrastructure/) folder is designed to hold more than one public cloud.  This allows extensibility into other vendors for future enhancements. 

This step will also install the AWS ALB load balancing controller to the provisioned cluster, and setup the appropriate OIDC and IAM bindings to allow the controller to create load balancer resources in the AWS env.  The resulting ingress objects will dynamically generate a host name to allow the REST server to be routable on the public internet. 

> NOTE: This step currently is configured to use an S3 bucket as the backend that the owner of this git repository owns.  If you do not have IAM delegated access to this bucket and perform an `aws configure`, then remove the backend config section from the [terraform.tf](./infrastructure/aws/terraform.tf) file to proceed (or update to point to an S3 bucket owned by the runner, and perform the required login steps to allow terraform to access the bucket). 

##### Actions Performed 

1. Terraform plan on the declared resources to determine what changes will be made to the environment 
2. Displays the plan to the user, so that it can be reviewed and confirmed before executing the plan
3. On confirmation, perform a terraform apply for the generated plan

##### Resources Deployed

1. A Kubernetes EKS cluster in the defined region
   1. A node group for the EC2 virtual host nodes 
      1. A singlular node in the node group that servers as a master and a worker for the PoC setup
2. A VPC that is tied to the EKS cluster
   1. 4 subnets under the VPC
      1. 2 Public 
      2. 2 Private
3. Various IAM bindings and security groups to allow network traffic on the VPC/Cluster and allow different components to communicate
4. A NAT gateway for the VPC to allow egress from the network
   1. This includes an internet gateway as well, including multiple network interfaces
5. An IAM policy that allows the ALB load balancer controller to perform the required actions on the AWS account
2. A service account to the k8s cluster that is bound to the IAM policy using a role 
3. The ALB load balancer stack
   1. Please view full list of resources and their manifest definitions here [ALB Controller](https://github.com/kubernetes-sigs/aws-load-balancer-controller/releases/download/v2.4.1/v2_4_1_full.yaml)
4. The Jetstack Cert Manager stack
   1. Please view full list of resources and their manifest definitions here [Jetstack](https://github.com/jetstack/cert-manager/releases/download/v1.6.0/cert-manager.yaml)

#### (2) Configure Cluster Credentials

##### Description

This step uses the `aws` command line tool to configure the `$HOME/.kube/config` file for cluster access to the environment deployed in the previous step. 

##### Actions Performed 

1. Use public cloud vendor CLI to generate the kubernetes cluster credential file

##### Resources Deployed

_None_

#### (3) Setup / Deploy PoC (k8s)

##### Description

This step will create the k8s Operator controller to watch and manage the CRD instances we defined for the REST server.  Once the manager is deployed and running, this will create an instance of the CRD (a CR) as well to perform the actual deployment of the server. 

##### Actions Performed 

1. Runs a `helm upgrade --install` to install the `CRDs`, `Controller` and `Application CR` to the k8s environment

##### Resources Deployed

1. Controller manifests:
   1. ClusterRole
   2. ClusterRoleBinding
   3. Namespace
      1. `acme-portfolio-example-manager`
   4. CRD
      1. `applications.acme.io`
   5. Deployment
      1. `controller-manager`
   6. ServiceAccount
      1. `controller-manager`
2. The application CR and the downstream resources the controller reconciles
   1. `application.acme.io` (`application-sample`)
      1. Application namespace
      2. Deployment
      3. ServiceAccount
      4. Service
      5. Ingress (only valid when the ALB load balancer is deployed)

> NOTE: There is currently not a step to configure an image pull secret, so the destination repository must have anonymous image pulling enabled, or the resulting deploy will be in a constant back off due to `ImagePullBackOff` in the container create step of the Pods. 

### Verifying the Install

The provided [verify_e2e.sh](./verify_e2e.sh) script can be used to verify the deployed environment. The load balancer install option is optional, and in the absence of the load balancer, the verify script will port forward the Pod's running port and verify from there.

## Tearing Down the PoC

To tear down the PoC the following can be run:

```sh
##########################
#                        #
# Uninstall the PoC Only #
#                        #
##########################

# The application must be deleted to cascade the deletion of cluster resources
# the load balancer in the AWS env must be deleted or it can cause 
# the subnet finalizers to not resolve properly resulting in orphaned resources. 
$ helm uninstall pocsample -n acme-controller-manager

####################
#                  #
# Teardown the Env #
#                  #
####################

# The IAM policy binding MAY need to be cleaned up
# before running the terraform destroy as in some
# cases there can be a delete conflict returned 
# by the AWS server when removing the resources
#
# TODO(any): we can exec this in the terraform and put a dependend on 
#            statement in the iam role clause for the cleanup of the attachment
$ aws iam detach-role-policy \
   --role-name=ingress_alb_controller \
   --policy-arn=$(aws iam list-policies --no-cli-pager --output yaml | yq -Mr '.Policies[] | select(.PolicyName == "ingress_policy") | .Arn')

# The following commands will remove the infrastructure
# provisioned in the PoC
$ cd infrastructure/aws

$ terraform destroy -auto-approve
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
