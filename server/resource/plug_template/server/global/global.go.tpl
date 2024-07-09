package global

{{- if .HasGlobal }}

import "47.103.136.241/goprojects/curesan/server/plugin/{{ .Snake}}/config"

var GlobalConfig = new(config.{{ .PlugName}})
{{ end -}}