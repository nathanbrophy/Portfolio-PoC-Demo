#!/usr/bin/env sh

set -e

VERSION="v1.0.0"
CRITICAL_ERROR=0
ROOT_DIR=`pwd`

function info {
    log "INFO" "${1}"
}

function error {
    log "ERRO" "${1}"
    CRITICAL_ERROR=1
}

function log {
    echo "[${1}] ${2}"
}

function verify_requisite {
    info "checking prerequisite ${1}"
    which "${1}" > /dev/null || error "required prerequisite ${1} not found on machine"
}

function check_prereqs {
    local reqs=(minikube docker go kubectl terraform yq aws)
    for r in "${reqs[@]}"; do 
        verify_requisite "${r}"
    done
}

function check_env_settings {
    local vars=(AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY)
    for v in in "${vars[@]}"; do 
        env | grep "${v}" > /dev/null || error "required environment setting ${v} not found"
    done
}

function server {
    info "building and deploying the server to remote registry"

    local make_flags="VERSION=$(cat deploy_config.yaml | yq -Mr '.deploy.server.tag') IMAGE_REPO=$(cat deploy_config.yaml | yq -Mr '.deploy.server.repository') IMAGE_REGISTRY=$(cat deploy_config.yaml | yq -Mr '.deploy.registry') OS=$(cat deploy_config.yaml | yq -Mr '.deploy.os') ARCH=$(cat deploy_config.yaml | yq -Mr '.deploy.arch')"

    cd server

    cmd="make docker-build ${make_flags}"
    eval "${cmd}"
    cmd="make docker-push  ${make_flags}"
    eval "${cmd}"

    cd "${ROOT_DIR}"
}

function controller {
    info "building and deploying controller to remote registry"

    local image_registry="$(cat deploy_config.yaml | yq -Mr '.deploy.registry')"
    local image_repository="$(cat deploy_config.yaml | yq -Mr '.deploy.controller.repository')"
    local image_tag="$(cat deploy_config.yaml | yq -Mr '.deploy.controller.tag')"
    local image="${image_registry}/${image_repository}:${image_tag}"
    local server_tag="$(cat deploy_config.yaml | yq -Mr '.deploy.server.tag') "
    local server_repository="$(cat deploy_config.yaml | yq -Mr '.deploy.server.repository')"
    local server_image="${image_registry}/${server_repository}:${server_tag}"
    server_image="$(printf ${server_image} | tr -d ' ')"

    local make_flags="IMG=${image}"

    cd k8s

    cmd="make docker-build ${make_flags}"
    eval "${cmd}"
    cmd="make docker-push  ${make_flags}"
    eval "${cmd}"

    info "installing CRDs to cluster"

    make install

    info "deploying controller to environment"

    kubectl kustomize config/default | sed "s%REPLACE_IMAGE%${image_registry}/${image_repository}:${image_tag}%g" | kubectl apply -f -

    info "create the CR to deploy the application to the environment"

    echo "\"${server_image}\""

    kubectl apply -f - <<EOF
apiVersion: acme.io/v1beta1
kind: Application
metadata:
  name: application-sample
spec:
  application:
    image: "${server_image}"
    port: 8081

EOF

    cd "${ROOT_DIR}"
}

function infrastructure {
    info "deploying the infrastructure to the cloud provider"

    local work_dir="$(mktemp -d)"
    local tf_plan_file="${work_dir}/state.tfplan"
    local region="$(cat deploy_config.yaml | yq -Mr '.deploy.terraform.region')"
    local cluster_name="$(cat deploy_config.yaml | yq -Mr '.deploy.terraform.cluster_name')"

    info "the working directory for the terraform plan is ${work_dir}"

    cd "infrastructure/aws"

    terraform plan -no-color -var region="${region}" -var cluster_name="${cluster_name}" > "${tf_plan_file}"

    # Block operation of the apply until the plan is approved
    echo 'Please press any key to continue and review the terraform plan'
    read DUMMY

    cat "${tf_plan_file}" | less

    echo "Conintue [y/N]: "
    read CONTINUE

    [[ "${CONTINUE}" != "y" ]] && fatal "user did not accept the terraform plan, canceling apply"

    terraform apply -var region="${region}" -var cluster_name="${cluster_name}" -auto-approve

    cd "${ROOT_DIR}"
    info "done standing up the AWS infrastructure"
}

function configure_cluster_credentials {
    info "configuring the kubeconfig to authenticate to current cluster"

    local region="$(cat deploy_config.yaml | yq -Mr '.deploy.terraform.region')"
    local cluster_name="$(cat deploy_config.yaml | yq -Mr '.deploy.terraform.cluster_name')"

    aws eks --region "${region}" update-kubeconfig --name "${cluster_name}"


    info "credentials configured, happy k8s-ing!"
}

function run {
    info "Version ${VERSION}"

    check_prereqs 
    check_env_settings

    [[ "${CRITICAL_ERROR}" -eq 1 ]] && exit 1

    info "all prerequisites have been verified"

    server
    infrastructure
    configure_cluster_credentials
    controller
}

run