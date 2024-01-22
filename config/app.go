package config

type App struct {
	Env        string `mapstructure:"env" json:"env" yaml:"env"`
	Port       string `mapstructure:"port" json:"port" yaml:"port"`
	AppName    string `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
	RunLogType string `mapstructure:"run_log_type" json:"run_log_type" yaml:"run_log_type"`
}
