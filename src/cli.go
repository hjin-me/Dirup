package main

import (
	"flag"
	"log"
	"path/filepath"
)

var (
	path = flag.String("c", "./conf.yaml", "tool config")
	dir  = flag.String("i", "./statics", "the directory to be uploaded")
)

func main() {
	flag.Parse()
	conf, err := filepath.Abs(*path)
	if err != nil {
		flag.Usage()
		return
	}
	directory, err := filepath.Abs(*dir)
	if err != nil {
		flag.Usage()
		return
	}
	_ = directory
	cfg, err := ReadConfig(conf)
	if err != nil {
		log.Println(err)
		return
	}

	_ = cfg
}
