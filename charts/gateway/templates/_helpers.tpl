{{- define "gateway.configureEnv" -}}
{{- $env := list -}}

{{- $baseURL := default "https://platform.example.com" .Values.gateway.platformBaseUrl -}}
{{- $env = append $env (dict "name" "PLATFORM_BASE_URL" "value" $baseURL) -}}

{{- $authSecret := trim (default "" .Values.gateway.authToken.existingSecret) -}}
{{- $authVar := dict "name" "PLATFORM_AUTH_TOKEN" -}}
{{- if $authSecret }}
  {{- $secretKey := default "platform-auth-token" .Values.gateway.authToken.existingSecretKey -}}
  {{- $_ := set $authVar "valueFrom" (dict "secretKeyRef" (dict "name" $authSecret "key" $secretKey)) -}}
{{- else }}
  {{- $_ := set $authVar "value" (default "" .Values.gateway.authToken.value) -}}
{{- end }}
{{- $env = append $env $authVar -}}

{{- $timeout := printf "%d" (default 10000 .Values.gateway.timeoutMs) -}}
{{- $env = append $env (dict "name" "PLATFORM_TIMEOUT_MS" "value" $timeout) -}}

{{- $retries := printf "%d" (default 2 .Values.gateway.retries) -}}
{{- $env = append $env (dict "name" "PLATFORM_RETRIES" "value" $retries) -}}

{{- $headers := default "" .Values.gateway.requestHeadersJson -}}
{{- $env = append $env (dict "name" "PLATFORM_REQUEST_HEADERS_JSON" "value" $headers) -}}

{{- $validate := printf "%t" (default false .Values.gateway.openapiValidateResponse) -}}
{{- $env = append $env (dict "name" "OPENAPI_VALIDATE_RESPONSE" "value" $validate) -}}

{{- $userEnv := .Values.env | default (list) -}}
{{- $_ := set .Values "env" (concat $env $userEnv) -}}
{{- end -}}
