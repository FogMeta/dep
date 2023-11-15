package misc

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

func ValidHttpURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

const (
	CharsetNum = "0123456789"
	CharsetAll = "0123456789abcdefghijklmnopqrABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func RandomString(n int, charset string) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func DownloadLinkSize(url string) (size int64, err error) {
	resp, err := http.Head(url)
	if err != nil {
		return size, err
	}
	contentLength := resp.Header.Get("Content-Length")
	size, err = strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, errors.New("invalid download url")
	}
	return
}

func EncodeStructValues(data any, tag string) (encoded string, err error) {
	rv := reflect.Indirect(reflect.ValueOf(data))
	rt := rv.Type()
	if rv.Kind() != reflect.Struct {
		return "", fmt.Errorf("only support struct type")
	}
	for i := 0; i < rt.NumField(); i++ {
		key := ""
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}
		if tag != "" {
			key = field.Tag.Get(tag)
			if key == "-" {
				continue
			}
		} else {
			key = field.Name
		}
		value := rv.Field(i).Interface()
		encoded += fmt.Sprintf("%s=%v", key, value)
	}
	return
}
