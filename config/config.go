package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	globalCfg  Config
	ErrNotFile = errors.New("not a file")
)

type Config struct {
	AccessKey string `yaml:"ak"`
	SecretKey string `yaml:"sk"`
	Prefix    string `yaml:"domain"`
	Bucket    string `yaml:"bucket"`
	Workers   int    `yaml:"workers"`
}

func LoadConfig() Config {
	return globalCfg
}

func ReadConfig(filename string) (cfg Config, err error) {

	if !filepath.IsAbs(filename) {
		err = errors.New("filename is not abs")
		return
	}
	fi, err := os.Lstat(filename)
	if err != nil {
		return
	}
	if fi.IsDir() {
		err = ErrNotFile
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	bf, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(bf, &cfg)
	if err != nil {
		return
	}
	if cfg.Workers == 0 {
		cfg.Workers = 5
	}
	globalCfg = cfg
	return
}
