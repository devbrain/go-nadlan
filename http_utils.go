package gonadlan

import (
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"io"
	"math/rand/v2"
	"net/http"
	"strings"
)

func SetStandardHeaders(req *http.Request, originUrl string) {
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("sec-ch-ua", "Chromium\";v=\"104\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"104")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	if rand.IntN(100)%2 == 0 {
		req.Header.Add("mobile-app", "false")
	} else {
		req.Header.Add("mobile-app", "true")
	}
	req.Header.Add("sec-ch-ua-mobile", "?0")
	if rand.IntN(100)%2 == 0 {
		req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.79 Safari/537.36")
	} else {
		req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:124.0) Gecko/20100101 Firefox/124.0")
	}
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("sec-ch-ua-platform", "Linux")
	req.Header.Add("Origin", originUrl)
	req.Header.Add("Sec-Fetch-Site", "same-site")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Referer", originUrl)
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
}

var noSpace = strings.NewReplacer(" ", "")

func splitEncodingHeader(raw string) []string {
	if raw == "" {
		return []string{}
	}
	return strings.Split(noSpace.Replace(raw), ",")
}

func ReadHTTPResponse(res *http.Response, err error) ([]byte, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error in body.Close")
		}
	}(res.Body)
	values := splitEncodingHeader(res.Header.Get("Content-Encoding"))
	var reader io.Reader = nil
	for i := len(values) - 1; i >= 0; i-- {
		v := values[i]
		switch v {
		case "br":
			reader = io.NopCloser(brotli.NewReader(res.Body))
		case "gzip", "x-gzip":
			reader, err = gzip.NewReader(res.Body)
			if err != nil {
				return nil, err
			}
		case "zstd":
			reader, err = zstd.NewReader(res.Body)
			if err != nil {
				return nil, err
			}
		case "", "identity":
		default:
			reader = nil
		}
	}
	var body []byte
	if reader == nil {
		body, err = io.ReadAll(res.Body)
	} else {
		body, err = io.ReadAll(reader)
	}

	if err != nil {
		return nil, err
	}
	return body, nil
}
