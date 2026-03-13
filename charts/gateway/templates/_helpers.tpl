{{- define "gateway.configureEnv" -}}
{{- $env := list -}}

{{- $teamsGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.teamsGrpcTarget) -}}
{{- $env = append $env (dict "name" "TEAMS_GRPC_TARGET" "value" $teamsGrpcTarget) -}}

{{- $filesGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.filesGrpcTarget) -}}
{{- $env = append $env (dict "name" "FILES_GRPC_TARGET" "value" $filesGrpcTarget) -}}

{{- $llmGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.llmGrpcTarget) -}}
{{- $env = append $env (dict "name" "LLM_GRPC_TARGET" "value" $llmGrpcTarget) -}}

{{- $secretsGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.secretsGrpcTarget) -}}
{{- $env = append $env (dict "name" "SECRETS_GRPC_TARGET" "value" $secretsGrpcTarget) -}}

{{- $validate := printf "%t" (default false .Values.gateway.openapiValidateResponse) -}}
{{- $env = append $env (dict "name" "OPENAPI_VALIDATE_RESPONSE" "value" $validate) -}}

{{- $userEnv := .Values.env | default (list) -}}
{{- $_ := set .Values "env" (concat $env $userEnv) -}}
{{- end -}}
