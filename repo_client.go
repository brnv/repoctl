package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/kovetskiy/toml"
)

type APIResponse struct {
	Success bool                `json:"success"`
	Error   string              `json:"error"`
	Data    map[string][]string `json:"data"`
}

const (
	formPackageFile = "package_file"
)

func (response *APIResponse) String() string {
	if !response.Success {
		return response.Error
	}

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
	repodURL    string
	method      string
	resource    string
	body        bytes.Buffer
	packageFile *os.File
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

func (client *RepoClient) LoadPackageFile(packageFile string) error {
	var err error

	client.packageFile, err = os.Open(packageFile)
	if err != nil {
		return err
	}

	return nil
}

func (client *RepoClient) doRequest() (*http.Response, error) {
	var (
		err error
		url = client.repodURL + client.resource
	)

	request, err := http.NewRequest(client.method, url, &client.body)
	if err != nil {
		return &http.Response{}, err
	}

	err = client.loadForm(request)
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

func (client *RepoClient) loadForm(request *http.Request) error {
	if client.packageFile == nil {
		return nil
	}

	form := multipart.NewWriter(&client.body)
	defer form.Close()

	formWriter, err := form.CreateFormFile(
		formPackageFile, client.packageFile.Name(),
	)
	if err != nil {
		return err
	}

	if _, err = io.Copy(formWriter, client.packageFile); err != nil {
		return err
	}

	request.Header.Set("Content-Type", form.FormDataContentType())

	return nil
}

func (client *RepoClient) appendURLParts(parts []string) {
	for index, part := range parts {
		if part == "" {
			break
		}

		client.resource = client.resource + part

		if index != len(parts)-1 {
			client.resource = client.resource + "/"
		}
	}
}
