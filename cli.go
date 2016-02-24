package main

import (
	"flag"
	"github.com/hjin-me/Dirup/config"
	"github.com/hjin-me/Dirup/tool"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"

	"fmt"
	"golang.org/x/net/context"
)

var (
	path = flag.String("c", "", "tool config, default is ~/.dirup.yaml")
	dir  = flag.String("i", "./statics", "the directory to be uploaded")
)

func main() {
	flag.Parse()
	if *path == "" {
		usr, err := user.Current()
		if err != nil {
		}
		*path = filepath.Join(usr.HomeDir, "./.dirup.yaml")
	}
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
	cfg, err := config.ReadConfig(conf)
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
				fmt.Fprintln(os.Stderr, err)
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
		os.Exit(1)
	}
}
