{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "sealed-secrets-httppubcert.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "sealed-secrets-httppubcert.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create a fully qualified app name with a suffix
*/}}
{{- define "sealed-secrets-httppubcert.fullnameSuffix" -}}
{{- $trimSuffix := .suffix | trunc 63 -}}
{{- $effectiveTruncLen := int (sub 63 (len $trimSuffix | int)) -}}
{{- if .scope.Values.fullnameOverride -}}
{{- print (.scope.Values.fullnameOverride | trunc $effectiveTruncLen | trimSuffix "-") $trimSuffix -}}
{{- else -}}
{{- $name := default .scope.Chart.Name .scope.Values.nameOverride -}}
{{- if contains $name .scope.Release.Name -}}
{{- print (.scope.Release.Name | trunc $effectiveTruncLen | trimSuffix "-") $trimSuffix -}}
{{- else -}}
{{- printf "%s-%s" .scope.Release.Name (print ($name | trunc $effectiveTruncLen | trimSuffix "-") $trimSuffix) -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "sealed-secrets-httppubcert.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "sealed-secrets-httppubcert.labels" -}}
app.kubernetes.io/name: {{ include "sealed-secrets-httppubcert.name" . }}
helm.sh/chart: {{ include "sealed-secrets-httppubcert.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}
