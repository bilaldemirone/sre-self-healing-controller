#!/bin/bash

echo "Removing deployments..."

kubectl delete -f deploy/controller --ignore-not-found=true
kubectl delete -f deploy/dummy-app --ignore-not-found=true

echo "Cleanup completed."