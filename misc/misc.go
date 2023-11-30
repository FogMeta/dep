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

func EncodeStructValues(data any, tag, sep string) (encoded string, err error) {
	rv := reflect.Indirect(reflect.ValueOf(data))
	rt := rv.Type()
	if rv.Kind() != reflect.Struct {
		return "", fmt.Errorf("only support struct type")
	}
	list := make([]string, 0, rt.NumField())
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
		list = append(list, fmt.Sprintf("%s=%v", key, value))
	}
	return strings.Join(list, sep), nil
}

func CompareStructValues(src, dst any, tag string, excludeCols ...string) (values map[string]any, err error) {
	sv := reflect.Indirect(reflect.ValueOf(src))
	dv := reflect.Indirect(reflect.ValueOf(dst))
	if sv.Kind() != reflect.Struct || dv.Kind() != reflect.Struct {
		return nil, errors.New("invalid type, only support struct")
	}
	st := sv.Type()
	dt := dv.Type()
	if st != dt {
		return nil, errors.New("src and dst must be the same type")
	}
	exm := make(map[string]bool, len(excludeCols))
	for _, col := range excludeCols {
		exm[col] = true
	}
	sm := make(map[string]any, st.NumField())
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if field.IsExported() {
			sm[field.Name] = sv.Field(i).Interface()
		}
	}

	values = make(map[string]any)
	for i := 0; i < dt.NumField(); i++ {
		field := dt.Field(i)
		if !field.IsExported() {
			continue
		}
		name := field.Name
		fv := dv.Field(i).Interface()
		if sm[name] == fv {
			continue
		}
		if tag != "" {
			name = field.Tag.Get(tag)
			if tag == "gorm" {
				name = ParseGormColumn(field.Name, name)
			}
		}
		if name != "" && !exm[name] {
			values[name] = fv
		}
	}
	return
}

func ParseGormColumn(name, tag string) string {
	list := strings.Split(tag, ";")
	columnKV := ""
	for _, v := range list {
		if strings.Contains(strings.ToLower(v), "column") {
			columnKV = v
			break
		}
	}
	if columnKV == "" {
		return name
	}
	values := strings.Split(columnKV, ":")
	if len(values) == 2 {
		return values[1]
	}
	return name
}
