#!/bin/bash

# Set the environment variables
VARIABLES=(
  "SERVER_PORT=8080"
  "AUTH_PATH=http://auth:8081/auth/token"
  "AUTH_HOST=kubernetes.auth.internal.eng.com"
  "MONGODB_URI=mongodb://admin:password@mongodb:27017/"
  "MONGODB_DATABASE=person"
  "MONGODB_COLLECTION=person"
  "MONGODB_USERNAME=admin"
  "MONGODB_PASSWORD=password"
  "MONGODB_ROLE=userAdminAnyDatabase"
)

echo "Setting environment variables in namespace $1"

for variable in "${VARIABLES[@]}"; do
  kubectl -n $1 set env deployment/people-charts ${variable}
done

echo "Done!"