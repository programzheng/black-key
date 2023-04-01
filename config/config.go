package config

var Cfg = New()

type Config interface {
	GetString(name string) string
	GetInt(name string) int
	GetBool(name string) bool
	Set(key string, value interface{})
}

func New() Config {
	config := &Instance{
		Package: NewViper(),
	}

	return config
}

func GetProductionStatus() bool {
	return Cfg.GetString("ENV") == "production"
}
