package lagrange

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/FogMeta/libra-os/module/log"
)

const (
	tokenURL = "/jwt_info"
)

func TokenValidate(host, jwt string) (wallet string, err error) {
	host += tokenURL
	client := &http.Client{}
	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+jwt)
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Error("resp body :", string(body))
		return "", fmt.Errorf("%d", resp.StatusCode)
	}
	var res result
	var info tokenInfo
	res.Data = &info
	if err = json.Unmarshal(body, &res); err != nil {
		log.Error(err)
		return
	}
	if res.Status != "success" {
		return "", errors.New(res.Message)
	}
	return info.WalletAddress, nil
}

type result struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type tokenInfo struct {
	WalletAddress string `json:"wallet_address"`
}
