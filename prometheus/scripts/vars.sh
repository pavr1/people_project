#!/bin/bash

# Set the environment variables
VARIABLES=(
  "PROMETHEUS_PORT=9000"
)

echo "Setting environment variables in namespace $1"

for variable in "${VARIABLES[@]}"; do
  kubectl -n $1 set env deployment/prometheus ${variable}
done

echo "Done!"