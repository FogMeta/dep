package lagrange

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ResultURL(url string) (string, error) {
	var resp ResultResp
	if err := httpGet(url, &resp); err != nil {
		return "", err
	}
	return resp.JobResultURI, nil
}

type ResultResp struct {
	JobResultURI string `json:"job_result_uri"`
}

func httpGet(url string, receiver any) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	if receiver != nil {
		return json.NewDecoder(resp.Body).Decode(receiver)
	}
	return
}
