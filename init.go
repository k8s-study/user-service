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

func createApi() {
	url := fmt.Sprintf("%s/apis/", os.Getenv("KONG_HOST"))
	// FIXME: change upstream URL for k8s
	api := Api{
		"users",
		"/v1/health",
		"GET",
		"http://127.0.0.1"}
	pbytes, _ := json.Marshal(api)
	buff := bytes.NewBuffer(pbytes)

	if _, err := http.Post(url, "application/json", buff); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("APIs registered.")
}

func enableApi() {
	url := fmt.Sprintf("%s/apis/users/plugins", os.Getenv("KONG_HOST"))
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
	createApi()
	enableApi()
}
