package global

{{- if .HasGlobal }}

import "github.com/EscapeBearSecond/curescan/server/plugin/{{ .Snake}}/config"

var GlobalConfig = new(config.{{ .PlugName}})
{{ end -}}