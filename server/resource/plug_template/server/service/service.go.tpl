package service

 {{- if .NeedModel }}
import (
   "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/plugin/{{ .Snake}}/model"
)
{{ end }}

type {{ .PlugName}}Service struct{}

func (e *{{ .PlugName}}Service) PlugService({{- if .HasRequest }}req model.Request {{ end -}}) ({{- if .HasResponse }}res model.Response,{{ end -}} err error) {
    // 写你的业务逻辑
	return {{- if .HasResponse }} res,{{ end }} nil
}
