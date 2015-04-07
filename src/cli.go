package main

import (
	"flag"
	"log"
	"path/filepath"
	"sync"
	"tool"

	"golang.org/x/net/context"
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
	cfg, err := tool.ReadConfig(conf)
	if err != nil {
		log.Println(err)
		return
	}

	_ = cfg
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	in := tool.ScanDir(ctx, directory)

	var wg sync.WaitGroup
	for filename := range in {
		wg.Add(1)
		go func(f string) {
			tool.UploadFile(ctx, f)
			wg.Done()
		}(filename)
	}

	wg.Wait()
}
