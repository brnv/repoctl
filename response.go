package main

import "encoding/json"

type APIResponse struct {
	Success    bool                `json:"success"`
	Error      string              `json:"error"`
	Data       map[string][]string `json:"data"`
	jsonOutput bool
}

func (response *APIResponse) getOutput() string {
	if response.jsonOutput {
		output, err := response.toJSON()
		if err != nil {
			return err.Error()
		}

		return output
	}

	if !response.Success {
		return response.Error
	}

	return response.String()
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
