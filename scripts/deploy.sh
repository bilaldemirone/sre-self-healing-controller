#!/bin/bash

set -e

echo "Deploying Controller..."

kubectl apply -f deploy/controller/rbac.yaml
kubectl apply -f deploy/controller/configmap.yaml
kubectl apply -f deploy/controller/deployment.yaml

echo "Deploying Dummy App..."

kubectl apply -f deploy/dummy-app/deployment.yaml

echo "Deployment completed."