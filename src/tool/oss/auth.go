package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
)

func HeaderSign(ak, sk, method, contentMD5, contentType, path, header string, expires time.Time) (string, error) {
	h1 := hmac.New(sha1.New, []byte(sk))
	if z, _ := expires.Zone(); z != "GMT" {
		l, err := time.LoadLocation("GMT")
		if err != nil {
			return "", err
		}
		expires = expires.In(l)
	}
	rawStr := fmt.Sprintf("%s\n%s\n%s\n%s\n", method, contentMD5, contentType, expires.Format(time.RFC1123))
	if header != "" {
		rawStr += header + "\n" + path
	} else {
		rawStr += path
	}
	// log.Printf("%s\n", rawStr)
	// log.Printf("% x\n", []byte(rawStr))

	h1.Write([]byte(rawStr))
	// log.Printf("% x\n", h1.Sum(nil))
	sign := base64.StdEncoding.EncodeToString(h1.Sum(nil))
	// log.Println(sign)
	return fmt.Sprintf("OSS %s:%s", ak, sign), nil
	/*
			"Authorization: OSS " + Access Key Id + ":" + Signature

			Signature = base64(hmac-sha1(AccessKeySecret,
		            VERB + "\n"
				    + CONTENT-MD5 + "\n"
					+ CONTENT-TYPE + "\n"
					+ DATE + "\n"
					+ CanonicalizedOSSHeaders
					+ CanonicalizedResource))
	*/

}

func QuerySign(sk, method, contentMD5, contentType, path, headers string, expires time.Time) string {

	h1 := hmac.New(sha1.New, []byte(sk))
	rawStr := fmt.Sprintf("%s\n%s\n%s\n%d\n%s\n%s", method, contentMD5, contentType, expires.Unix(), headers, path)

	return base64.StdEncoding.EncodeToString(h1.Sum([]byte(rawStr)))

	/*
			Signature = base64(hmac-sha1(AccessKeySecret,
		          VERB + "\n"
			            + CONTENT-MD5 + "\n"
					          + CONTENT-TYPE + "\n"
						            + EXPIRES + "\n"
								          + CanonicalizedOSSHeaders
									            + CanonicalizedResource))
	*/
}
