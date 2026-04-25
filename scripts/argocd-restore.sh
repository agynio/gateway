#!/usr/bin/env bash
set -euo pipefail

if [ "${CI:-}" = "true" ]; then
  echo "Skipping ArgoCD auto-sync restore in CI."
  exit 0
fi

if kubectl get application gateway -n argocd >/dev/null 2>&1; then
  echo "Re-enabling ArgoCD auto-sync for gateway..."
  kubectl patch application gateway -n argocd \
    --type merge \
    -p '{"spec":{"syncPolicy":{"automated":{"prune":true,"selfHeal":true}}}}'
  echo "ArgoCD auto-sync restored."
fi
