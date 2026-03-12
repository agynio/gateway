{{- define "gateway.configureEnv" -}}
{{- $env := list -}}

{{- $baseURL := required "gateway.platformBaseUrl is required" (trimAll " \n\t" (default "" .Values.gateway.platformBaseUrl)) -}}
{{- $env = append $env (dict "name" "PLATFORM_BASE_URL" "value" $baseURL) -}}

{{- $teamsGrpcTarget := required "gateway.teamsGrpcTarget is required" (trimAll " \n\t" (default "" .Values.gateway.teamsGrpcTarget)) -}}
{{- $env = append $env (dict "name" "TEAMS_GRPC_TARGET" "value" $teamsGrpcTarget) -}}

{{- $filesGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.filesGrpcTarget) -}}
{{- if $filesGrpcTarget }}
{{- $env = append $env (dict "name" "FILES_GRPC_TARGET" "value" $filesGrpcTarget) -}}
{{- end }}

{{- $llmGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.llmGrpcTarget) -}}
{{- if $llmGrpcTarget }}
{{- $env = append $env (dict "name" "LLM_GRPC_TARGET" "value" $llmGrpcTarget) -}}
{{- end }}

{{- $llmHttpBaseUrl := trimAll " \n\t" (default "" .Values.gateway.llmHttpBaseUrl) -}}
{{- if $llmHttpBaseUrl }}
{{- $env = append $env (dict "name" "LLM_HTTP_BASE_URL" "value" $llmHttpBaseUrl) -}}
{{- end }}

{{- $secretsGrpcTarget := trimAll " \n\t" (default "" .Values.gateway.secretsGrpcTarget) -}}
{{- if $secretsGrpcTarget }}
{{- $env = append $env (dict "name" "SECRETS_GRPC_TARGET" "value" $secretsGrpcTarget) -}}
{{- end }}

{{- $authSecret := trim (default "" .Values.gateway.authToken.existingSecret) -}}
{{- $authVar := dict "name" "PLATFORM_AUTH_TOKEN" -}}
{{- if $authSecret }}
  {{- $secretKey := default "platform-auth-token" .Values.gateway.authToken.existingSecretKey -}}
  {{- $_ := set $authVar "valueFrom" (dict "secretKeyRef" (dict "name" $authSecret "key" $secretKey)) -}}
{{- else }}
  {{- $_ := set $authVar "value" (default "" .Values.gateway.authToken.value) -}}
{{- end }}
{{- $env = append $env $authVar -}}

{{- $timeout := int (default 10000 .Values.gateway.timeoutMs) -}}
{{- $env = append $env (dict "name" "PLATFORM_TIMEOUT_MS" "value" (printf "%d" $timeout)) -}}

{{- $retries := int (default 2 .Values.gateway.retries) -}}
{{- $env = append $env (dict "name" "PLATFORM_RETRIES" "value" (printf "%d" $retries)) -}}

{{- $headers := default "" .Values.gateway.requestHeadersJson -}}
{{- $env = append $env (dict "name" "PLATFORM_REQUEST_HEADERS_JSON" "value" $headers) -}}

{{- $validate := printf "%t" (default false .Values.gateway.openapiValidateResponse) -}}
{{- $env = append $env (dict "name" "OPENAPI_VALIDATE_RESPONSE" "value" $validate) -}}

{{- $userEnv := .Values.env | default (list) -}}
{{- $_ := set .Values "env" (concat $env $userEnv) -}}
{{- end -}}
