package global

{{- if .HasGlobal }}

import "github.com/EscapeBearSecond/eaglesnest/server/plugin/{{ .Snake}}/config"

var GlobalConfig = new(config.{{ .PlugName}})
{{ end -}}