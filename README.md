# Portfolio Demo

- [Portfolio Demo](#portfolio-demo)
  - [Overview](#overview)
    - [Infrastructure](#infrastructure)
    - [K8s](#k8s)
    - [Server](#server)
  - [Architecture Diagrams](#architecture-diagrams)
  - [Running the PoC](#running-the-poc)

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

To run the PoC sample the [deploy_config.yaml](./deploy_config.yaml) __must__ be configured to provide script configuration options.  Once the configuration is completed, the [deploy_poc.sh](./deploy_poc.sh) script can be ran to fully stand the PoC up on a public cloud environment. 

Before running the PoC please ensure the following credential steps are taken:

1. `docker login` performed to ensure the images can be pushed to their repositories
2. `aws configure` performed to ensure access to the configured S3 bucket to store the TF state (optional)
