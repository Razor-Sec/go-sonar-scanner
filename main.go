package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	//"github.com/TwiN/go-color"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

type healthSQ struct {
	Health string `json:"health"`
}

func systemHealth(username string, password string, url, method string) string {
	url = url + "/api/system/health"
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Errorf("Got error %s", err)
	}

	req.SetBasicAuth(username, password)
	//var bearer = "Bearer " + token
	//req.Header.Add("Authorization", bearer)
	response, err := client.Do(req)

	if err != nil {
		fmt.Errorf("Got error %s", err.Error())
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Errorf("Getting error %s", err)
	}
	//fmt.Println(string(body))
	var healthcheck healthSQ
	if err := json.Unmarshal(body, &healthcheck); err != nil {
		fmt.Errorf("Get error %s", err)
	}
	var stats = response.StatusCode

	if stats == 403 {
		println("403 UnAuthorized")
	}
	return string(healthcheck.Health)
}

func QG(username string, password string, url_base string, qgname string, projectKey string) {
	var newUrlbase = url_base + "/api/projects/create?project=" + projectKey + "&name=" + projectKey
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", newUrlbase, nil)
	if err != nil {
		fmt.Errorf("Got error %s", err.Error())
	}
	req.SetBasicAuth(username, password)
	response, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Got error %s", err.Error())
	}

	defer response.Body.Close()

	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Errorf("Getting error %s", err)
	// }
	// println(string(body))
	ChangeQG(username, password, url_base, qgname, projectKey)
}

func ChangeQG(username string, password string, url_base string, qgname string, projectKey string) {
	url_base = url_base + "/api/qualitygates/select?gateName=" + qgname + "&projectKey=" + projectKey
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("POST", url_base, nil)
	if err != nil {
		fmt.Errorf("Got error %s", err.Error())
	}
	req.SetBasicAuth(username, password)
	response, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Got error %s", err.Error())
	}

	defer response.Body.Close()

	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Errorf("Getting error %s", err)
	// }
	// println(string(body))
}

func scanner(baseurl string, args string) {
	var sonar string = basepath + "/resources/bin/sonar-scanner"
	baseurl = "-Dsonar.host.url=" + baseurl
	testargs := strings.Split(args, " ")
	newargs := append([]string{sonar, baseurl}, testargs...)
	cmd := exec.Command("bash", newargs...)

	println(cmd.String())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}
	err = cmd.Start()
	fmt.Println("INFO: EXEC sonar-scanner")
	if err != nil {
		fmt.Println(err)
	}

	// print the output of the subprocess
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	scanner2 := bufio.NewScanner(stderr)
	for scanner2.Scan() {
		m := scanner2.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}

func main() {

	projectKey := flag.String("projectKey", "", "For project key sonarqube")
	qgname := flag.String("qualityGate", "", "For quality gate name")
	args := flag.String("args", "", "For Arguments sonar-scanner")
	methodAuth := flag.String("auth", "", "For authentication env or flag")
	username := flag.String("username", "", "For Arguments sonar-scanner")
	password := flag.String("password", "", "For Arguments sonar-scanner")
	baseurl := flag.String("baseurl", "http://127.0.0.1:9001", "For url sonar-scanner")
	flag.Parse()
	if *methodAuth == "env" {
		if os.Getenv("sonaruser") != "" && os.Getenv("sonarpass") != "" {
			println("System Health : ", systemHealth(os.Getenv("sonaruser"), os.Getenv("sonarpass"), *baseurl, "GET"))
			QG(os.Getenv("sonaruser"), os.Getenv("sonarpass"), *baseurl, *qgname, *projectKey)
			scanner(*baseurl, *args)
		} else {
			fmt.Println("[!] env sonaruser & sonarpass not found")
		}
	}
	testdoang, err := exec.Command("pwd").Output()
	if err != nil {
		fmt.Println(err)
	}
	println(string(testdoang))
	if *methodAuth == "flag" {
		println("System Health : ", systemHealth(*username, *password, *baseurl, "GET"))
		QG(*username, *password, *baseurl, *qgname, *projectKey)
		scanner(*baseurl, *args)
	}
}
