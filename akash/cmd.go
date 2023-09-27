package akash

import (
	"log"
	"os"
	"os/exec"
)

func CreateAccount(name string) (err error) {
	cmd := exec.Command("bash", "akash/create_account.sh", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(out))
	return
}

func Deploy(name, ymlPath string) (err error) {
	_, err = os.Stat(ymlPath)
	if err != nil {
		return
	}
	cmd := exec.Command("bash", "akash/deploy.sh", name, ymlPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(out))
	return
}
