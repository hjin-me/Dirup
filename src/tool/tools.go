package tool

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"tool/bos"

	"golang.org/x/net/context"
	"gopkg.in/yaml.v2"
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
	globalCfg = cfg
	return
}

func readDir(ctx context.Context, ch chan string, dir string) {
	fileList, _ := ioutil.ReadDir(dir)
	for _, v := range fileList {
		fp := filepath.Join(dir, v.Name())
		if strings.HasPrefix(v.Name(), ".") {
			continue
		}
		if v.IsDir() {
			readDir(ctx, ch, fp)
		} else {
			select {
			case <-ctx.Done():
				return
			case ch <- fp:
			}
		}
	}
}

func ScanDir(ctx context.Context, dir string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		readDir(ctx, ch, dir)
	}()

	return ch

}

func UploadFile(ctx context.Context, root, filename string) error {
	return bos.UploadFile(ctx, root, filename)
}
