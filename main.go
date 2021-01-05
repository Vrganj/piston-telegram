package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	telebot "gopkg.in/tucnak/telebot.v2"
)

// Language language
type Language struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

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

func getLanguages() []Language {
	response, err := http.Get("https://emkc.org/api/v1/piston/versions")

	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	var languages []Language
	json.Unmarshal(bytes, &languages)

	response.Body.Close()

	return languages
}

func getLanguage(text string) string {
	index := strings.Index(text, " ")

	if index == -1 {
		return strings.ToLower(text)
	}

	return strings.ToLower(text[:index])
}

func getSource(text string) string {
	index := strings.Index(text, " ")

	if index == -1 {
		return ""
	}

	return text[index+1:]
}

func validateLanguage(language string) bool {
	for _, l := range languages {
		if l.Name == language {
			return true
		}
	}

	return false
}

func execute(message *telebot.Message, language string, source string) {
	requestBody := ExecuteRequest{
		Language: language,
		Source:   source,
		Args:     []string{},
	}

	requestBodyBuffer, _ := json.Marshal(requestBody)

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

	bot.Send(message.Sender, responseBody.Output)
}

func runCommand(message *telebot.Message) {
	text := strings.TrimSpace(message.Text[4:])
	language := getLanguage(text)
	source := getSource(text)

	fmt.Println(text)

	if language == "" {
		bot.Send(message.Sender, "Provide a language")
		return
	}

	if !validateLanguage(language) {
		bot.Send(message.Sender, "Invalid language provided")
		return
	}

	if source == "" {
		bot.Send(message.Sender, "Provide source code")
		return
	}

	execute(message, language, source)
}

var bot *telebot.Bot
var languages []Language

func main() {
	languages = getLanguages()

	// TODO: Find a better way to assign bot
	b, err := telebot.NewBot(telebot.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 5 * time.Second},
	})

	bot = b

	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/run", runCommand)

	bot.Start()
}
