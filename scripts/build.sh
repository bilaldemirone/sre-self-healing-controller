#!/bin/bash

set -e

echo "======================================"
echo "Building SRE Controller"
echo "======================================"

cd controller-runtime

docker build -t sre-controller:v1 .

kind load docker-image sre-controller:v1 --name sre-case

echo "Controller image loaded successfully."