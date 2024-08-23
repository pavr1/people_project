#!/bin/bash

# Set the environment name
ENV_NAME="snbx" 

# Set the environment variables
VARIABLES=(
  "AUTH_PORT=8081"
)

# Loop through the variables and add them to the environment
for variable in "${VARIABLES[@]}"; do
  kubectl -n ${ENV_NAME} set env deployment/auth ${variable}
done