package tool

import (
	"path/filepath"
	"testing"
	"time"
)

func TestUploadOneFile(t *testing.T) {
	path, err := filepath.Abs("test/conf.yaml")
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := ReadConfig(path)
	if err != nil {
		t.Fatal(err)
	}

	expires, err := time.Parse(time.RFC1123, "Tue, 07 Apr 2015 08:32:42 GMT")
	if err != nil {
		t.Fatal(err)
	}

	sign, err := HeaderSign(cfg.AccessKey, cfg.SecretKey, "PUT", "", "application/octet-stream", "/bucket/test.log", "", expires)
	t.Log(sign)
	if err != nil {
		t.Fatal(err)
	}
	if sign != "OSS abc:obrLliKQha0F7QUwlLx7BXJ5Jxg=" {
		t.Error("sign not ok")
	}

}
