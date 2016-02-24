package oss

import (
	"errors"
	"github.com/hjin-me/Dirup/config"
	"github.com/hjin-me/Dirup/mimes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/context"
)

func UploadFile(ctx context.Context, root, filename string) error {
	relativePath, err := filepath.Rel(root, filename)
	if err != nil {
		return err
	}
	cfg := config.LoadConfig()
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	request, err := http.NewRequest("PUT", cfg.Prefix+relativePath, fd)
	if err != nil {
		return err
	}
	l, err := time.LoadLocation("GMT")
	if err != nil {
		return err
	}
	gmt := time.Now().In(l)

	mm := mimes.MIME(filename)

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
