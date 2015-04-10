package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
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

	var (
		success []string
		fail    []string
	)

	var process = func(ctx context.Context, f <-chan string) {
		for filename := range f {
			err := tool.UploadFile(ctx, directory, filename)
			if err != nil {
				fail = append(fail, filename)
			} else {
				success = append(success, filename)
			}

			time.Sleep(time.Millisecond * 300)
		}
		wg.Done()
	}
	wg.Add(5)
	go process(ctx, in)
	go process(ctx, in)
	go process(ctx, in)
	go process(ctx, in)
	go process(ctx, in)

	wg.Wait()

	if len(fail) > 0 {

		fname, err := filepath.Abs("./log")
		if err != nil {
			log.Println(success)
			log.Println(fail)
			log.Fatal(err)
		}
		output := ""
		for _, s := range success {
			output += s + " success\n"
		}
		for _, s := range fail {
			output += s + " fail\n"
		}

		ioutil.WriteFile(fname, []byte(output), os.ModePerm)
	}
}
