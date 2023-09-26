package lagrange

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	match = "lagrangedao.org/spaces/0x\\w+/.+"
)

func DownloadSpace(spaceURL string, dir string) (path string, err error) {
	rxp, err := regexp.Compile(match)
	if err != nil {
		return
	}
	if !rxp.Match([]byte(spaceURL)) {
		return "", errors.New("invalid space url")
	}
	if strings.HasPrefix(spaceURL, "https://lagrangedao.org") {
		spaceURL = strings.ReplaceAll(spaceURL, "https://lagrangedao.org", "https://api.lagrangedao.org/")
		spaceURL = strings.TrimSuffix(spaceURL, "/app")
		spaceURL += "/files"
	}
	resp, err := http.Get(spaceURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var res SpaceFileResp
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return
	}
	if res.Status != "success" {
		return "", errors.New(res.Status)
	}
	if len(res.Data) == 0 {
		return "", errors.New("have no files")
	}
	name := res.Data[0].Name
	list := strings.Split(name, "/")
	if len(list) < 3 {
		return "", errors.New("invalid name")
	}
	prefix := filepath.Join(list[:2]...)
	for _, file := range res.Data {
		file.Name = strings.TrimPrefix(file.Name, prefix)
		if err = downloadFile(file, dir); err != nil {
			return
		}
	}
	return filepath.Join(dir, list[2]), nil
}

func downloadFile(file *File, dir string) error {
	path := filepath.Join(dir, file.Name)
	log.Println("path: ", path)
	if err := os.MkdirAll(filepath.Dir(path), 0766); err != nil {
		return err
	}
	resp, err := http.Get(file.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fw, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fw.Close()
	_, err = io.Copy(fw, resp.Body)
	return err
}

type SpaceFileResp struct {
	Data    []*File     `json:"data"`
	Message interface{} `json:"message"`
	Status  string      `json:"status"`
}

type File struct {
	Cid       string `json:"cid"`
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updated_at"`
	URL       string `json:"url"`
}
