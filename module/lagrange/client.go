package lagrange

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/FogMeta/libra-os/module/log"
)

const (
	contentTypeJson = "application/json"
	contentTypeForm = "application/x-www-form-urlencoded"
)

type Client struct {
	Host string
	http.Client
}

func (client *Client) postJson(token string, method string, data any, receiver any) error {
	return client.doRequest(token, "POST", contentTypeJson, data, method, receiver)
}

func (client *Client) postForm(token string, method string, data any, receiver any) error {
	return client.doRequest(token, "POST", contentTypeForm, data, method, receiver)
}

func (client *Client) get(token string, method string, data url.Values, receiver any) error {
	return client.doRequest(token, "GET", "", nil, method, receiver)
}

func (client *Client) doRequest(token string, method string, contentType string, data any, apiMethod string, receiver any) error {
	if receiver != nil && reflect.ValueOf(receiver).Kind() != reflect.Ptr {
		return errors.New("receiver must be a pointer")
	}
	link := client.Host + apiMethod
	var reader io.Reader
	switch method {
	case http.MethodPost:
		switch contentType {
		case contentTypeJson:
			if data != nil {
				b, _ := json.Marshal(data)
				reader = bytes.NewReader(b)
			}
		case contentTypeForm:
			switch data.(type) {
			case url.Values:
				values := data.(url.Values)
				reader = strings.NewReader(values.Encode())
			case string:
				reader = strings.NewReader(data.(string))
			default:
				return errors.New("invalid data type, must be url.Values or encoded values string in form mode")
			}
			values, ok := data.(url.Values)
			if !ok {
				return errors.New("invalid data type, must be url.Values in form mode")
			}
			reader = strings.NewReader(values.Encode())
		default:
			return fmt.Errorf("not supported content type: %s", contentType)
		}

	case http.MethodGet:
		values, ok := data.(url.Values)
		if ok {
			link += "?" + values.Encode()
		}
	default:
		return fmt.Errorf("not supported method %s", method)
	}
	req, err := http.NewRequest(method, link, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	log.Info("response:", string(b))
	var result Result
	result.Data = receiver
	if err = json.Unmarshal(b, receiver); err != nil {
		return err
	}
	if result.Status != "success" {
		return errors.New(result.Message)
	}
	return nil
}

type Result struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
