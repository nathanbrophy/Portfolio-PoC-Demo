#!/usr/bin/env bash

# set -e is explicitly not set here because there are some tolerable errors we can silently ignore

VERSION="v2.0.0"
CRITICAL_ERROR=0
ROOT_DIR=`pwd`
_SKIP_BUILD="$(cat deploy_config.yaml | yq -Mr '.deploy.no_build')"
BLANK='\033[0m'
RED='\033[1;31m'
CYAN='\033[1;36m'

function info {
    log "${CYAN}INFO${BLANK}" "${1}"
}

function error {
    log "${RED}ERRO${BLANK}" "${1}"
    CRITICAL_ERROR=1
}

function log {
    echo -e "[${1}] ${2}"
}

function verify_requisite {
    info "checking prerequisite ${1}"
    which "${1}" > /dev/null || error "required prerequisite ${1} not found on machine"
}

function check_prereqs {
    local reqs=(kubectl yq aws)
    for r in "${reqs[@]}"; do 
        verify_requisite "${r}"
    done
}

function controller {
    info "building and deploying controller to remote registry"

    helm upgrade -i operatorsample operator-controller --namespace acme-controller-manager --create-namespace 

    cd "${ROOT_DIR}"
}

function configure_cluster_credentials {
    info "configuring the kubeconfig to authenticate to current cluster"

    local region="$(cat deploy_config.yaml | yq -Mr '.deploy.terraform.region')"
    local cluster_name="$(cat deploy_config.yaml | yq -Mr '.deploy.terraform.cluster_name')"

    aws eks --region "${region}" update-kubeconfig --name "${cluster_name}"


    info "credentials configured, happy k8s-ing!"
}

function display_rest_endpoint_with_sample {
    info "gathering cluster information to display the application REST endpoint"

    info "REST server is reachable at:"
    info "    $(kubectl get ing acme-application -o yaml | yq -Mr '.status.loadBalancer.ingress[].hostname'):80/example"
    echo ""
}

function context_namespace {
    local ns="$(cat deploy_config.yaml | yq -Mr '.deploy.namespace')"

    info "setting kkubectl context to ${ns}"

    kubectl get ns "${ns}" > /dev/null || kubectl create ns "${ns}"

    kubectl config set-context --current --namespace=${ns}
}

function run {
    info "Version ${VERSION}"

    check_prereqs 

    [[ "${CRITICAL_ERROR}" -eq 1 ]] && exit 1

    info "all prerequisites have been verified"

    configure_cluster_credentials
    context_namespace
    controller

    # TODO: replace this with a until do loop in display REST information
    sleep 30s
    
    display_rest_endpoint_with_sample
}

run

# Verify the install was successfule
source ./verify_e2e.sh 
