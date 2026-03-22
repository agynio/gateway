{{- define "gateway.configureEnv" -}}
{{- $env := list -}}

{{- $agentsGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.agentsGrpcTarget) -}}
{{- $env = append $env (dict "name" "AGENTS_GRPC_TARGET" "value" $agentsGrpcTarget) -}}

{{- $threadsGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.threadsGrpcTarget) -}}
{{- $env = append $env (dict "name" "THREADS_GRPC_TARGET" "value" $threadsGrpcTarget) -}}

{{- $chatGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.chatGrpcTarget) -}}
{{- $env = append $env (dict "name" "CHAT_GRPC_TARGET" "value" $chatGrpcTarget) -}}

{{- $notificationsGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.notificationsGrpcTarget) -}}
{{- $env = append $env (dict "name" "NOTIFICATIONS_GRPC_TARGET" "value" $notificationsGrpcTarget) -}}

{{- $filesGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.filesGrpcTarget) -}}
{{- $env = append $env (dict "name" "FILES_GRPC_TARGET" "value" $filesGrpcTarget) -}}

{{- $agentStateGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.agentStateGrpcTarget) -}}
{{- $env = append $env (dict "name" "AGENT_STATE_GRPC_TARGET" "value" $agentStateGrpcTarget) -}}

{{- $tokenCountingGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.tokenCountingGrpcTarget) -}}
{{- $env = append $env (dict "name" "TOKEN_COUNTING_GRPC_TARGET" "value" $tokenCountingGrpcTarget) -}}

{{- $llmGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.llmGrpcTarget) -}}
{{- $env = append $env (dict "name" "LLM_GRPC_TARGET" "value" $llmGrpcTarget) -}}

{{- $secretsGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.secretsGrpcTarget) -}}
{{- $env = append $env (dict "name" "SECRETS_GRPC_TARGET" "value" $secretsGrpcTarget) -}}

{{- $usersGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.usersGrpcTarget) -}}
{{- $env = append $env (dict "name" "USERS_GRPC_TARGET" "value" $usersGrpcTarget) -}}

{{- $oidcIssuerUrl := trimAll " \n\t" (default "" .Values.gateway.oidcIssuerUrl) -}}
{{- if $oidcIssuerUrl -}}
{{- $env = append $env (dict "name" "OIDC_ISSUER_URL" "value" $oidcIssuerUrl) -}}
{{- end -}}

{{- $oidcClientId := trimAll " \n\t" (default "" .Values.gateway.oidcClientId) -}}
{{- if $oidcClientId -}}
{{- $env = append $env (dict "name" "OIDC_CLIENT_ID" "value" $oidcClientId) -}}
{{- end -}}

{{- $zitiEnabled := .Values.gateway.zitiEnabled | default false -}}
{{- $env = append $env (dict "name" "ZITI_ENABLED" "value" (printf "%t" $zitiEnabled)) -}}

{{- $zitiManagementGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.zitiManagementGrpcTarget) -}}
{{- $env = append $env (dict "name" "ZITI_MANAGEMENT_GRPC_TARGET" "value" $zitiManagementGrpcTarget) -}}

{{- $zitiLeaseRenewalInterval := trimAll " \n\t" (default "" .Values.gateway.zitiLeaseRenewalInterval) -}}
{{- if $zitiLeaseRenewalInterval -}}
{{- $env = append $env (dict "name" "ZITI_LEASE_RENEWAL_INTERVAL" "value" $zitiLeaseRenewalInterval) -}}
{{- end -}}

{{- $userEnv := .Values.env | default (list) -}}
{{- $_ := set .Values "env" (concat $env $userEnv) -}}
{{- end -}}
