{{- define "gateway.configureEnv" -}}
{{- $env := list -}}

{{- $agentsGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.agentsGrpcTarget) -}}
{{- $env = append $env (dict "name" "AGENTS_GRPC_TARGET" "value" $agentsGrpcTarget) -}}

{{- $appsGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.appsGrpcTarget) -}}
{{- $env = append $env (dict "name" "APPS_GRPC_TARGET" "value" $appsGrpcTarget) -}}

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

{{- $meteringGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.meteringGrpcTarget) -}}
{{- $env = append $env (dict "name" "METERING_GRPC_TARGET" "value" $meteringGrpcTarget) -}}

{{- $secretsGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.secretsGrpcTarget) -}}
{{- $env = append $env (dict "name" "SECRETS_GRPC_TARGET" "value" $secretsGrpcTarget) -}}

{{- $usersGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.usersGrpcTarget) -}}
{{- $env = append $env (dict "name" "USERS_GRPC_TARGET" "value" $usersGrpcTarget) -}}

{{- $organizationsGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.organizationsGrpcTarget) -}}
{{- $env = append $env (dict "name" "ORGANIZATIONS_GRPC_TARGET" "value" $organizationsGrpcTarget) -}}

{{- $runnersGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.runnersGrpcTarget) -}}
{{- $env = append $env (dict "name" "RUNNERS_GRPC_TARGET" "value" $runnersGrpcTarget) -}}

{{- $exposeGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.exposeGrpcTarget) -}}
{{- $env = append $env (dict "name" "EXPOSE_GRPC_TARGET" "value" $exposeGrpcTarget) -}}

{{- $egressRulesGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.egressRulesGrpcTarget) -}}
{{- $env = append $env (dict "name" "EGRESS_RULES_GRPC_TARGET" "value" $egressRulesGrpcTarget) -}}

{{- $groupsGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.groupsGrpcTarget) -}}
{{- $env = append $env (dict "name" "GROUPS_GRPC_TARGET" "value" $groupsGrpcTarget) -}}

{{- $networksGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.networksGrpcTarget) -}}
{{- $env = append $env (dict "name" "NETWORKS_GRPC_TARGET" "value" $networksGrpcTarget) -}}

{{- $oidcIssuerUrl := trimAll " \n\t" (default "" .Values.gateway.oidcIssuerUrl) -}}
{{- if $oidcIssuerUrl -}}
{{- $env = append $env (dict "name" "OIDC_ISSUER_URL" "value" $oidcIssuerUrl) -}}
{{- end -}}

{{- $oidcClientId := trimAll " \n\t" (default "" .Values.gateway.oidcClientId) -}}
{{- if $oidcClientId -}}
{{- $env = append $env (dict "name" "OIDC_CLIENT_ID" "value" $oidcClientId) -}}
{{- end -}}

{{- $clusterAdminToken := trimAll " \n\t" (default "" .Values.gateway.clusterAdminToken) -}}
{{- if $clusterAdminToken -}}
{{- $env = append $env (dict "name" "CLUSTER_ADMIN_TOKEN" "value" $clusterAdminToken) -}}
{{- end -}}

{{- $clusterAdminIdentityId := trimAll " \n\t" (default "" .Values.gateway.clusterAdminIdentityId) -}}
{{- if $clusterAdminIdentityId -}}
{{- $env = append $env (dict "name" "CLUSTER_ADMIN_IDENTITY_ID" "value" $clusterAdminIdentityId) -}}
{{- end -}}

{{- $userEnv := .Values.env | default (list) -}}
{{- $_ := set .Values "env" (concat $env $userEnv) -}}
{{- end -}}
