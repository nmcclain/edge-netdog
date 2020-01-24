package main

import (
	"fmt"
	"io/ioutil"
	"time"

	docopt "github.com/docopt/docopt-go"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Global struct {
		Debug              bool
		MonitorIntervalRaw string        `yaml:"monitor_interval"`
		MonitorInterval    time.Duration `yaml:"-"`
		TargetAttempts     int           `yaml:"target_attempts"`
		TargetUrl          string        `yaml:"target_url"`
		TargetMatch        string        `yaml:"target_match"`
		ActionDelayRaw     string        `yaml:"action_delay"`
		ActionDelay        time.Duration `yaml:"-"`
		configFile         string        `yaml:"-"`
	}
	Actions []string `yaml:"actions"`
}

func getcfg() (Config, error) {
	config := Config{}
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		return config, err
	}
	config.Global.configFile = args["<config>"].(string)
	configFile, err := ioutil.ReadFile(config.Global.configFile)
	if err != nil {
		return config, fmt.Errorf("opening config file: %s", err.Error())
	}
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return config, fmt.Errorf("parsing config file: %s", err.Error())
	}
	if config.Global.Debug {
		log.SetLevel(logrus.DebugLevel)
	}
	if config.Global.TargetAttempts < 1 {
		config.Global.TargetAttempts = 3
	}
	if len(config.Global.MonitorIntervalRaw) < 1 {
		config.Global.MonitorIntervalRaw = "1m"
	}
	config.Global.MonitorInterval, err = time.ParseDuration(config.Global.MonitorIntervalRaw)
	if err != nil {
		return config, fmt.Errorf("parsing monitor_interval: %v", err)
	}
	if len(config.Global.ActionDelayRaw) < 1 {
		config.Global.ActionDelayRaw = "30s"
	}
	config.Global.ActionDelay, err = time.ParseDuration(config.Global.ActionDelayRaw)
	if err != nil {
		return config, fmt.Errorf("parsing action_delay: %v", err)
	}
	log.Debugf("Using configuration: %+v", config)
	return config, nil
}
