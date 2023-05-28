# k8s

- [k8s](#k8s)
  - [Description](#description)
    - [Controllers](#controllers)
    - [APIs](#apis)
    - [DriftDtection](#driftdtection)
    - [Generators](#generators)
  - [Getting Started](#getting-started)
    - [Running on the cluster](#running-on-the-cluster)
    - [Uninstall CRDs](#uninstall-crds)
    - [Undeploy controller](#undeploy-controller)
  - [Contributing](#contributing)
    - [How it works](#how-it-works)
    - [Test It Out](#test-it-out)
    - [Modifying the API definitions](#modifying-the-api-definitions)
  - [License](#license)

k8s is a sample operator controller that will reconcile a CRD definition with boilerplate information to standup an application on a k8s cluster. 

## Description

The operator in this repository has a few key components:

1. Controllers
2. APIs
3. DriftDetection
4. Generators

### Controllers

Holds a collection of tests and functions to act as the operator controller, that runs in the manager to reconcile the cluster state.

### APIs

Holds a collection of API versions that implement the overall `Application` API, that is used in reconciliation to generate the correct downstream manifests.  Currently, only `v1beta1` is a supported API version.

### DriftDtection

Holds a collection of boolean functions to test two objects for drift.  This is used to determine if reconciliation needs to be ran for a generated object, or if the cluster state is at parity with the expected state.

### Generators

Holds a collection of kubernetes object generators that are used to derive the downstream manifests needed to deploy the application from the CR coolected from the cluster. 

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/k8s:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/k8s:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

The controller also has an [e2e](./e2e/) test suite that can be ran to prove out integration functionality of the reconcile logic.  

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023 Nathan Brophy.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

