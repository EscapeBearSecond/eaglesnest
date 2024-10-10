package config

type Login struct {
	Attempts         int   `mapstructure:"attempts" json:"attempts" yaml:"attempts"`                         // 尝试次数
	DisabledDuration int64 `mapstructure:"disabledDuration" json:"disabledDuration" yaml:"disabledDuration"` // 禁用时长
}
