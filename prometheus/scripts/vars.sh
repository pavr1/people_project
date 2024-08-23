#!/bin/bash

# Set the environment name
ENV_NAME="snbx" 

# Set the environment variables
VARIABLES=(
  "PROMETHEUS_PORT=9000"
)

for variable in "${VARIABLES[@]}"; do
  kubectl -n ${ENV_NAME} set env deployment/prometheus ${variable}
done