package bsd_testtool

type Feature interface {
	Do(command []byte) error
}

type RepeatFeature struct {
	Times    int `toml:"times" mapstructure:"times"`
	Interval int `toml:"interval" mapstructure:"interval"`
}
