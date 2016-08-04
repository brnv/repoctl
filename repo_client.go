package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kovetskiy/toml"
)

type APIResponse struct {
	Success bool                `json:"success"`
	Error   string              `json:"error"`
	Data    map[string][]string `json:"data"`
}

func (response *APIResponse) String() string {
	output := ""
	for _, list := range response.Data {
		if len(list) == 0 {
			return ""
		}

		for index, element := range list {
			output = output + element
			if index != len(list)-1 {
				output = output + "\n"
			}
		}
	}
	return output
}

func (response *APIResponse) toJSON() (string, error) {
	output, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

type RepoClient struct {
	repodURL string
	method   string
	resource string
	body     io.Reader
}

func NewRepodClient(address string, version string) *RepoClient {
	repod := &RepoClient{
		repodURL: address + "/" + version,
		resource: "/",
	}

	return repod
}

func (client *RepoClient) Do() (APIResponse, error) {
	response, err := client.doRequest()
	if err != nil {
		return APIResponse{}, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return APIResponse{}, err
	}

	var apiResponse APIResponse

	_, err = toml.Decode(string(responseBody), &apiResponse)
	if err != nil {
		return APIResponse{}, err
	}

	return apiResponse, nil
}

func (client *RepoClient) doRequest() (*http.Response, error) {
	var (
		err error
		url = client.repodURL + client.resource
	)

	request, err := http.NewRequest(client.method, url, client.body)
	if err != nil {
		return &http.Response{}, err
	}

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)
	if err != nil {
		return &http.Response{}, err
	}

	return response, nil
}

func (client *RepoClient) appendURLParts(
	repo string,
	epoch string,
	db string,
	arch string,
) {
	resources := []string{repo, epoch, db, arch}
	for _, resource := range resources {
		if resource == "" {
			break
		}

		client.resource = client.resource + resource + "/"
	}
}
