#!/usr/bin/env bash

#
# Vars
#
_NS="default"
CONTROLLER_NS=
APP_NS=

function fail {
    echo "[FAIL]"
    exit 1
}

function help {
    echo "Usage: $(basename ${BASH_SOURCE}) --controller-namespace STR --application-namespace STR"
    echo ""
    echo "Arguments:"
    echo ""
    echo "Flags:"
    echo "--controller-namespace     the namespace the controller is running in"
    echo "--application-namespace    the namespace the application is running in"
    echo ""
}

function e2e {
    # Verify cluster connection
    echo "Veriying cluster connection"

    kubectl cluster-info > /dev/null || fail

    echo "[PASS]"

    # Verify CRDs are present

    echo "Verify CRDs are present"

    kubectl get crd applications.acme.io > /dev/null || fail

    echo "[PASS]"

    # Verify controller running

    echo "Verify controller running"

    kubectl -n ${CONTROLLER_NS} get po -l control-plane=controller-manager | grep Running > /dev/null || fail

    echo "[PASS]"

    # Verify application is running

    echo "Verify application is running"

    kubectl -n ${APP_NS} get po -l app=acme-application | grep Running > /dev/null || fail

    echo "[PASS]"

    # Verify application is working

    echo "Verify application is working"

    url="localhost:8081/example"
    # Check for load balancer in cluster env
    lb=$(kubectl -n ${APP_NS} get ing acme-application -o yaml | yq -Mr '.status.loadBalancer.ingress[].hostname')
    if [[ ! -z $lb ]] && [[ "${lb}" != "" ]]; then 
        url="${lb}:80/example"
    else 
        # Begin the port forward 
        kubectl -n ${APP_NS} port-forward service/acme-application 8081:8081 > /dev/null &
    PID=$!
    fi

    curl "${url}" | grep Automate || fail

    if [[ ! -z $PID ]]; then kill -9 $PID; fi

    echo "[PASS]"
}

#
# Main
#
while [ -n "${1}" ]; do
    case "${1}" in
        --controller-namespace)
        shift
        CONTROLLER_NS="${1}"
        shift
        ;;
        --application-namespace)
        shift
        APP_NS="${1}"
        shift
        ;;
        -h|--help)
        help
        exit 0
        ;;
        *)
        echo "[ERRO] ${1} is not a supported flag"
        help
        exit 1
        ;;
    esac
done

if [[ -z "${CONTROLLER_NS}" ]]; then 
    echo "[WARN] --controller-namespace not provided, using default namespace"
    CONTROLLER_NS="${_NS}"
fi

if [[ -z "${APP_NS}" ]]; then 
    echo "[WARN] --application-namespace not provided, using default namespace"
    APP_NS="${_NS}"
fi

e2e