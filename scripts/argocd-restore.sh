#!/usr/bin/env bash
set -euo pipefail

if kubectl get application gateway -n argocd >/dev/null 2>&1; then
  echo "Re-enabling ArgoCD auto-sync for gateway..."
  kubectl patch application gateway -n argocd \
    --type merge \
    -p '{"spec":{"syncPolicy":{"automated":{"prune":true,"selfHeal":true}}}}'
  echo "ArgoCD auto-sync restored."
else
  echo "WARNING: ArgoCD Application 'gateway' not found in argocd namespace."
fi
