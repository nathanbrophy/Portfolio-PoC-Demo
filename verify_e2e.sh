#!/usr/bin/env bash

function fail {
    echo "[FAIL]"
    exit 1
}

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

kubectl -n acme-portfolio-example-manager get po -l control-plane=controller-manager | grep Running > /dev/null || fail

echo "[PASS]"

# Verify application is running

echo "Verify application is running"

kubectl get po -l app=acme-application | grep Running > /dev/null || fail

echo "[PASS]"

# Verify application is working

echo "Verify application is working"

url="localhost:8081/example"
# Check for load balancer in cluster env
lb=$(kubectl get ing acme-application -o yaml | yq -Mr '.status.loadBalancer.ingress[].hostname')
if [[ ! -z $lb ]]; then 
    url="${lb}:80/example"
else 
    # Begin the port forward 
    kubectl port-forward service/acme-application 8081:8081 > /dev/null &
    PID=$!
fi

curl "${url}" | grep Automate || fail

if [[ ! -z $PID ]]; then kill -9 $PID; fi

echo "[PASS]"