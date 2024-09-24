package global

{{- if .HasGlobal }}

import "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/plugin/{{ .Snake}}/config"

var GlobalConfig = new(config.{{ .PlugName}})
{{ end -}}