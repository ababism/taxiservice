package mstracer

type Config struct {
	Enable    bool   `mapstructure:"enable"`
	ExpTarget string `mapstructure:"exp_target"`
	// "host.docker.internal:4317"
	StdOut    bool   `mapstructure:"stdout"`
	TraceFile string `mapstructure:"trace_file"`
}
