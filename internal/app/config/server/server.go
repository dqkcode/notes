package serverconf

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type (
	Config struct {
		Addr         string        `envconfig:"ADDR" default:":8080"`
		ReadTimeOut  time.Duration `envconfig:"READ_TIMEOUT" default:"5m"`
		WriteTimeOut time.Duration `envconfig:"WRITE_TIMEOUT" default:"5m"`
	}
)

func Load() (Config, error) {
	var conf Config
	if err := envconfig.Process("", &conf); err != nil {
		return conf, err
	}
	return conf, nil
}
