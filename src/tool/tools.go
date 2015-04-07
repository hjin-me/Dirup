package tool

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	relativePath, err := filepath.Rel(root, filename)
	if err != nil {
		return err
	}
	cfg := LoadConfig()
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", cfg.Prefix+relativePath, fd)
	if err != nil {
		return err
	}
	l, err := time.LoadLocation("GMT")
	if err != nil {
		return err
	}
	gmt := time.Now().In(l)

	mm := MIME(filename)

	sign, err := HeaderSign(cfg.AccessKey, cfg.SecretKey, "PUT", "", mm, filepath.Join("/", cfg.Bucket, relativePath), "", gmt)
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", sign)
	request.Header.Add("Cache-Control", "max-age=864000")
	request.Header.Add("Expires", "864000")
	request.Header.Add("Date", gmt.Format(time.RFC1123))
	request.Header.Add("Content-Type", mm)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(body) != "" {
		log.Printf("%s\n", body)
		return errors.New("response not empty")
	}
	log.Printf("upload [%s] completed", filename)
	return nil
}
