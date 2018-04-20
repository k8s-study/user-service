package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Api struct {
	Name         string `json:"name"`
	Uris         string `json:"uris"`
	Methods      string `json:"methods"`
	Upstream_Url string `json:"upstream_url"`
}

type Plugin struct {
	Name string `json:"name"`
}

func createApi(name string, path string, method string) {
	url := fmt.Sprintf("%s/apis/", os.Getenv("KONG_HOST"))
	// FIXME: change upstream URL for k8s
	api := Api{
		name,
		path,
		method,
		"http://user-service"}
	pbytes, _ := json.Marshal(api)
	buff := bytes.NewBuffer(pbytes)

	if _, err := http.Post(url, "application/json", buff); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("APIs registered.")
}

func enableApi(name string) {
	url := fmt.Sprintf("%s/apis/%s/plugins", os.Getenv("KONG_HOST"), name)
	plugin := Plugin{"key-auth"}
	pbytes, _ := json.Marshal(plugin)
	buff := bytes.NewBuffer(pbytes)

	if _, err := http.Post(url, "application/json", buff); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Key-auth plugin enabled.")
}

func main() {
	createApi("user-service-health", "/v1/health", "GET")
	createApi("user-service-signup", "/v1/signup", "POST")
	createApi("user-service-login", "/v1/login", "POST")
	createApi("user-service-userinfo", "/v1/users/[0-9]+", "GET")

	enableApi("user-service-userinfo")
}
