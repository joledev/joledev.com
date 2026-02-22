#!/bin/bash
# Manual deploy script for joledev K3s cluster
# Usage: ssh to VPS then run this, or run remotely:
#   ssh -p 2222 joel@69.62.68.130 'bash -s' < scripts/deploy.sh
set -euo pipefail

echo "==> Restarting joledev deployments..."
kubectl rollout restart deployment/web -n joledev
kubectl rollout restart deployment/api-quoter -n joledev
kubectl rollout restart deployment/api-scheduler -n joledev
kubectl rollout restart deployment/gatus -n joledev

echo "==> Waiting for rollout..."
kubectl rollout status deployment -n joledev --timeout=120s

echo "==> Current pod status:"
kubectl get pods -n joledev

echo "==> Deploy complete!"
