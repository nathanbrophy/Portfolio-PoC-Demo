#!/usr/bin/env bash

MANAGER_PID=
ROOT_DIR=`pwd`
REPO_ROOT_DIR="$GOPATH/src/github.com/nathanbrophy/portfolio-demo/k8s"
NAMESPACE="e2e-controller-test"

function log {
    echo "$(date '+%d-%m-%Y.%H:%M:%S') [e2e test] $@" 2>&1
}

function info {
    log "[INFO]" $@
}

function error {
    log "[ERRO]" $@
}

function fatal {
    error $@
    kill "${MANAGER_PID}"
    exit 1
}

function check_prereqs {
    local reqs=(minikube kubectl)

    for r in "${reqs[@]}"; do
        which "${r}" > /dev/null || fatal "pre-requisite ${r} not met, cannot proceed with e2e test suite"
    done
}

function setup {
    info "setting up test environment"

    minikube status | grep Running > /dev/null || minikube start > /dev/null

    kubectl create ns "${NAMESPACE}"

    cd "${REPO_ROOT_DIR}"

    make install > /dev/null

    local tries=0
    until kubectl get crd 'applications.acme.io' > /dev/null; do
        [[ "${tries}" -eq 10 ]] && fatal "setup step failed, CRDs failed to come online in 10 attempts"
        info "waiting for CRD to come online"
        sleep 5
    done

    make build > /dev/null

    go run ./main.go > /dev/null 2>&1 &
    MANAGER_PID=$!

    info "manager PID is '${MANAGER_PID}'"

    sleep 10

    cd "${SCRIPT_DIR}"
}

function teardown {
    info "tearing down test env"
    kill "${MANAGER_PID}"
    minikube delete > /dev/null
}

function e2e {
    info "start e2e test"
    info "applying CR to cluster"

    kubectl apply -f "${REPO_ROOT_DIR}/config/samples/default_v1beta1_application.yaml" -n "${NAMESPACE}"

    local tries=0
    until kubectl -n "${NAMESPACE}" get po -l "app=acme-application" | grep Running; do
        [[ "${tries}" -eq 25 ]] && fatal "e2e test failed the downstream Pod never came online"
        info "waiting for pod to come online"
        sleep 5
    done 

    local pod_name="$(kubectl -n ${NAMESPACE} get po -l app=acme-application -o name | awk -F'/' '{print $2}')"

    local kc_exec="kubectl -n ${NAMESPACE} exec -it ${pod_name}"
    local data="$(${kc_exec} -- curl acme-application.${NAMESPACE}.svc.cluster.local:8081/example)"
    
    printf "${data}" | grep "time" > /dev/null || fatal "e2e test failed the data does not contain a time element in the json"

    local status="$(${kc_exec} -- curl -X GET -I acme-application.${NAMESPACE}.svc.cluster.local:8081/example)"
    echo "${status}"

    printf "${status}" | grep 200 | grep OK  > /dev/null || fatal "e2e test failed the return status for a GET command was not a 200 OK"

    info "e2e test finished... [PASS]"
}

function run {
    check_prereqs
    setup
    e2e
    teardown
}

run