package tool

import (
	"github.com/hjin-me/Dirup/tool/oss"
	"io/ioutil"
	"path/filepath"
	"strings"

	"golang.org/x/net/context"
)

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
	return oss.UploadFile(ctx, root, filename)
}
