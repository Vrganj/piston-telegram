package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	telebot "gopkg.in/tucnak/telebot.v2"
)

// ExecuteRequest request
type ExecuteRequest struct {
	Language string   `json:"language"`
	Source   string   `json:"source"`
	Args     []string `json:"args"`
}

// OkExecuteResponse - Success response
type OkExecuteResponse struct {
	Ran      bool   `json:"ran"`
	Language string `json:"language"`
	Version  string `json:"version"`
	Output   string `json:"output"`
}

// ErrExecuteResponse - Error response
type ErrExecuteResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func execute(message *telebot.Message, bot *telebot.Bot, language string, source string) {
	requestBody := ExecuteRequest{
		Language: language,
		Source:   source,
		Args:     []string{},
	}

	requestBodyBuffer, err := json.Marshal(requestBody)

	if err != nil {
		bot.Send(message.Sender, "Something had gone wrong serializing the request body")
		return
	}

	response, err := http.Post("https://emkc.org/api/v1/piston/execute", "application/json", bytes.NewBuffer(requestBodyBuffer))

	if err != nil {
		bot.Send(message.Sender, "Something had gone wrong sending the request")
		return
	}

	responseBodyBuffer, err := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		bot.Send(message.Sender, "Something had gone wrong reading the response body")
		return
	}

	if response.StatusCode != 200 {
		responseBody := ErrExecuteResponse{}

		if json.Unmarshal(responseBodyBuffer, &responseBody) != nil {
			bot.Send(message.Sender, "Something had gone wrong deserializing the execute error response")
			return
		}

		bot.Send(message.Sender, fmt.Sprintf("Error %s: %s", responseBody.Code, responseBody.Message))
		return
	}

	responseBody := OkExecuteResponse{}

	if json.Unmarshal(responseBodyBuffer, &responseBody) != nil {
		bot.Send(message.Sender, "Something had gone wrong deserializing the execute response")
	}

	bot.Send(message.Sender, "```" + responseBody.Output + "```")
}
