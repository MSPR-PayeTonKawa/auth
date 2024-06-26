#!/bin/bash

for file in *.yaml; do
    kubectl apply -f "$file" || { echo "Error applying $file"; exit 1; }
done

kubectl rollout restart $DEPLOYMENT_NAME