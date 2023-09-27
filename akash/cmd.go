package akash

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CreateAccount(name string) (err error) {
	cmd := exec.Command("bash", "akash/create_account.sh", name)
	return DoCmd(cmd)
}

func Deploy(name, ymlPath string) (err error) {
	_, err = os.Stat(ymlPath)
	if err != nil {
		return
	}
	cmd := exec.Command("bash", "akash/deploy.sh", name, ymlPath)
	return DoCmd(cmd)
}

func syncLog(reader io.ReadCloser) {
	cache := ""
	buf := make([]byte, 1024)
	for {
		strNum, err := reader.Read(buf)
		if strNum > 0 {
			outputByte := buf[:strNum]
			outputSlice := strings.Split(string(outputByte), "\n")
			logText := strings.Join(outputSlice[:len(outputSlice)-1], "\n")
			log.Printf("%s%s", cache, logText)
			cache = outputSlice[len(outputSlice)-1]
		}
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "file already closed") {
				err = nil
				if cache != "" {
					log.Println(cache)
				}
			}
		}
	}
}

func DoCmd(cmd *exec.Cmd) (err error) {
	cmdStdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}
	cmdStderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		return
	}
	if err = cmd.Start(); err != nil {
		log.Println(err)
		return
	}
	go syncLog(cmdStdoutPipe)
	go syncLog(cmdStderrPipe)
	return cmd.Wait()
}
