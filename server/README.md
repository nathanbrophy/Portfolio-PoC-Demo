# Server

- [Server](#server)
  - [Run](#run)
    - [Binary](#binary)
    - [Container](#container)
  - [Test](#test)


The server project is an example REST server that handles requests at the `/example` route and returns a JSON payload with today's date. 

## Run

To run the project you can either run through the container or through the direct binary.

### Binary

```sh
% make run
```

### Container

```sh
% make docker-build IMAGE_TAG="${TAG}"

% docker run -it -d -p 8081:8081 "${TAG}"

% curl localhost:8081/example
```

## Test

All tests are runnable through the Makefile and this will run the `e2e` and the unit tests:

```
├── pkg
│   └── automationHandler
│       └── automationHandler_test.go  # Unit test for the server handlers
│   └── server
│       └── server_test.go  # Unit and e2e tests for the server
```

```sh
% make test
```
