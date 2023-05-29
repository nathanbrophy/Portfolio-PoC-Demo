#!/usr/bin/env bash

# set -e is explicitly not set here because there are some tolerable errors we can silently ignore

VERSION="v1.0.0"
CRITICAL_ERROR=0
ROOT_DIR=`pwd`
_SKIP_BUILD="$(cat deploy_config.yaml | yq -Mr '.deploy.no_build')"
BLANK='\033[0m'
RED='\033[1;31m'
CYAN='\033[1;36m'

function build_and_push {
    printf "${_SKIP_BUILD}" | grep '[tT][rR][uU][eE]' > /dev/null && return

    cmd="make docker-build ${1}"
    eval "${cmd}"
    cmd="make docker-push  ${1}"
    eval "${cmd}"
}

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
    local reqs=(docker go kubectl terraform yq aws eksctl)
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

    build_and_push "${make_flags}"

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

    build_and_push "${make_flags}"

    info "installing CRDs to cluster"

    make install

    info "deploying controller to environment"

    kubectl kustomize config/default | sed "s%REPLACE_IMAGE%${image_registry}/${image_repository}:${image_tag}%g" | kubectl apply -f -

    info "create the CR to deploy the application to the environment"

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

    [[ "${CONTINUE}" != "y" ]] && error "user did not accept the terraform plan, canceling apply"

    [[ "${CRITICAL_ERROR}" -eq 1 ]] && exit 1

    terraform apply -var region="${region}" -var cluster_name="${cluster_name}" -auto-approve

    cd "${ROOT_DIR}"
    info "done standing up the AWS infrastructure"

    rm -rf "${work_dir}"
}

function configure_cluster_credentials {
    info "configuring the kubeconfig to authenticate to current cluster"

    local region="$(cat deploy_config.yaml | yq -Mr '.deploy.terraform.region')"
    local cluster_name="$(cat deploy_config.yaml | yq -Mr '.deploy.terraform.cluster_name')"

    aws eks --region "${region}" update-kubeconfig --name "${cluster_name}"


    info "credentials configured, happy k8s-ing!"
}

function setup_load_balancer {
    local work_dir="$(mktemp -d)"

    info "begin instalation of the nginx controller"
    info "working directory available at (${work_dir})"
    
    local region="$(cat deploy_config.yaml | yq -Mr '.deploy.terraform.region')"
    local cluster_name="$(cat deploy_config.yaml | yq -Mr '.deploy.terraform.cluster_name')"
    local iam_policy_path="${work_dir}/iam-policy.json"
    local controller_path="${work_dir}/ingress-controller.yaml"
    local iam_url="https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/main/docs/install/iam_policy.json"
    local jet_stack_release="https://github.com/jetstack/cert-manager/releases/download/v1.6.0/cert-manager.yaml"
    local load_balancer_controller_release="https://github.com/kubernetes-sigs/aws-load-balancer-controller/releases/download/v2.4.1/v2_4_1_full.yaml"
    local role_name_eks="eks-${cluster_name}-loadbalancing-role"

    info "downloading the iam policy requirements from the latest stable release (${iam_url})"

    curl -o ${iam_policy_path} ${iam_url}

    info "creating the iam policy"

    cmd="aws iam create-policy --policy-name AWSLoadBalancerControllerIAMPolicy --policy-document file://${iam_policy_path}"
    
    out="$(eval ${cmd})"
    local arn_lookup_required=$?
    local policy_arn=""

    if [[ $arn_lookup_required -ne 0 ]]; then
        policy_arn=$(aws iam list-policies --no-cli-pager --output yaml | yq -Mr '.Policies[] | select(.PolicyName == "AWSLoadBalancerControllerIAMPolicy") | .Arn')
    else 
        policy_arn="$(printf ${out} | jq -Mr '.Policy.Arn')"   
    fi

    info "policy ARN is (${policy_arn})"

    eksctl create iamserviceaccount --cluster="${cluster_name}" --namespace=kube-system --name=aws-load-balancer-controller --attach-policy-arn="${policy_arn}" --region="${region}" --override-existing-serviceaccounts --approve --role-name="${role_name_eks}"

    info "downloading jetstack release from (${jet_stack_release})"
    # Install jetstack certmanager as a requirement for enabling TLS on exposed routes
    kubectl apply --validate=false -f "${jet_stack_release}"

    # Force an update to the aws load balancer CRDs
    cmd="kubectl apply -f https://raw.githubusercontent.com/aws/eks-charts/master/stable/aws-load-balancer-controller/crds/crds.yaml"
    eval $cmd

    sleep 120 

    Install the controller 
    info "downloading load balancer controller release from (${load_balancer_controller_release})"

    curl -Lo ${controller_path} ${load_balancer_controller_release}

    cat "${controller_path}" | sed "s/your-cluster-name/${cluster_name}/g" | kubectl apply -f -

    # Patch the service account to allow the ARN access to complete and actually
    # act on behalf of the user for the cluster. 
    #
    # Alternatively the eksctl command can be ran agian to update the metadata dynamically
    eksctl create iamserviceaccount --cluster="${cluster_name}" --namespace=kube-system --name=aws-load-balancer-controller --attach-policy-arn="${policy_arn}" --region="${region}" --override-existing-serviceaccounts --approve --role-name="${role_name_eks}"

    info "load balancer configuration complete"

    rm -rf "${work_dir}"
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
    check_env_settings

    [[ "${CRITICAL_ERROR}" -eq 1 ]] && exit 1

    info "all prerequisites have been verified"

    server
    infrastructure
    configure_cluster_credentials

    info "checking to see if load balancer requires instalation"
    cat deploy_config.yaml | yq -Mr '.deploy.load_balancer.deploy' | grep '[tT][rR][uU][eE]' > /dev/null &&
    setup_load_balancer

    context_namespace
    controller

    display_rest_endpoint_with_sample
}

run

# Verify the install was successfule
source ./verify_e2e.sh 
