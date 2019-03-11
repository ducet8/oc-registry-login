// Create ~/.oc.info with a single line containing <user>:<password>:<registry>
// Make that file 600 permissions

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func getFileName() (fileName string) {
	usr, userErr := user.Current()

	if userErr == nil {
		fileName = usr.HomeDir + "/.oc.info"
	}

	return fileName
}

func getToken() (token string) {
	ocTokenCmd := exec.Command("oc", "whoami", "-t")
	ocTokenOutput, ocTokenErr := ocTokenCmd.CombinedOutput()

	if ocTokenErr == nil {
		token = strings.TrimRight(string(ocTokenOutput), "\r\n")
		return
	}

	fmt.Println("Error retriving oc token")
	os.Exit(1)
	return
}

func getVars(fileName string) (usr, key, registry string) {
	fileData, fileErr := ioutil.ReadFile(fileName)

	if fileErr == nil {
		info := strings.Split(string(fileData), ":")
		usr, key, registry = info[0], info[1], strings.TrimRight(info[2], "\r\n")
	}

	return usr, key, registry
}

func loginDocker(dockerRegistry, ocUser, token string) () {
	if token != "" {
		dockerCmd := exec.Command("docker", "login", "-u", ocUser, "-p", token, dockerRegistry)
		dockerOutput, dockerErr := dockerCmd.CombinedOutput()
		
		if dockerErr == nil {
			fmt.Println("Logged in to Docker Registry")
			os.Exit(0)
		} 
		
		fmt.Println("Docker Error: " + fmt.Sprint(dockerErr) + ": " + string(dockerOutput))
		os.Exit(1)
	}
}

func loginOC(ocUser, ocKey string) (bool) {
	ocCmd := exec.Command("oc", "login", "-u", ocUser, "-p", ocKey)
	ocOutput, ocErr := ocCmd.CombinedOutput()
	
	if ocErr == nil {
		fmt.Println("Logged in to OC")
		return true
	}
	
	fmt.Println("oc Error: " + fmt.Sprint(ocErr) + ": " + string(ocOutput))
	os.Exit(1)
	return false
}

func main() {
	ocUser, ocKey, dockerRegistry := getVars(getFileName())
	if loginOC(ocUser, ocKey) {
		loginDocker(dockerRegistry, ocUser, getToken())
	}
}