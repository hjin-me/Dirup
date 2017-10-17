package bos

import (
	"errors"
	"github.com/hjin-me/Dirup/mimes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"context"
)

func UploadFile(ctx context.Context, root, filename string) error {
	relativePath, err := filepath.Rel(root, filename)
	if err != nil {
		return err
	}
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", "http://bj.bcebos.com/v1/{bucket}/"+relativePath, fd)
	// request, err := http.NewRequest("PUT", "http://cp01-rdqa-dev301.cp01.baidu.com:8097/v1/staticstest/"+relativePath, fd)
	if err != nil {
		return err
	}

	mm := mimes.MIME(filename)
	fInfo, err := fd.Stat()
	if err != nil {
		return err
	}
	l, err := time.LoadLocation("GMT")
	if err != nil {
		return err
	}
	gmt := time.Now().In(l)

	request.ContentLength = fInfo.Size()
	request.Header.Add("Content-Type", mm)
	request.Header.Set("Date", gmt.Format(time.RFC1123))
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
	// log.Println(resp.StatusCode, resp.Header)
	if string(body) != "" {
		log.Printf("%s\n", body)
		return errors.New("response not empty")
	}
	log.Printf("upload [%s] completed", filename)
	return nil
}
