package config

import (
	"io/ioutil"

	"github.com/aktsk/bqnotify/lib/notify"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Project string
	SQL     string
	Slack   notify.Slack
}

func Parse() (*Config, error) {
	buf, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
