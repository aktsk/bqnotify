package config

import (
	"io/ioutil"
	"os"

	"github.com/aktsk/bqnotify/lib/notify"

	"gopkg.in/yaml.v2"
)

// Config has configurations of BigQuery and Slack
type Config struct {
	Project string
	Queries []Query
	Slack   *notify.Slack
}

type Query struct {
	SQL   string
	Slack *notify.Slack
}

// Parse parses config.yaml
func Parse(file string) (*Config, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	conf := &Config{}
	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return nil, err
	}

	if conf.Project == "" {
		conf.Project = os.Getenv("BQNOTIFY_PROJECT")
	}

	if conf.Project == "" {
		conf.Project = os.Getenv("GCP_PROJECT")
	}

	return conf, nil
}
