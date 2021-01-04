package main

import (
	"bytes"
	"encoding/json"
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

// ExecutionRequest request
type ExecutionRequest struct {
	Language string   `json:"language"`
	Source   string   `json:"source"`
	Args     []string `json:"args"`
}

// OkExecutionResponse response
type OkExecutionResponse struct {
	Ran      bool   `json:"ran"`
	Language string `json:"language"`
	Version  string `json:"version"`
	Output   string `json:"output"`
}

// ErrExecutionResponse TODO
type ErrExecutionResponse struct{}

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

func main() {
	languages := getLanguages()

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &telebot.LongPoller{Timeout: 5 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/run", func(message *telebot.Message) {
		// Really fragile, needs fixing
		// and needs to tell the user if
		// the input was wrong

		text := message.Text[4:]

		if len(text) == 0 {
			return
		}

		text = text[1:]

		language := strings.ToLower(text[:strings.Index(text, " ")])
		matched := false

		for _, l := range languages {
			if l.Name == language {
				matched = true
				break
			}
		}

		if !matched {
			return
		}

		requestBody := ExecutionRequest{
			Language: language,
			Source:   text[strings.Index(text, " ")+1:],
			Args:     []string{},
		}

		requestBodyBytes, err := json.Marshal(requestBody)

		if err != nil {
			log.Fatal(err)
			return
		}

		response, err := http.Post("https://emkc.org/api/v1/piston/execute", "application/json", bytes.NewBuffer(requestBodyBytes))

		if err != nil {
			log.Fatal(err)
			return
		}

		responseBytes, err := ioutil.ReadAll(response.Body)

		if err != nil {
			log.Fatal(err)
			return
		}

		var executionResponse OkExecutionResponse

		if json.Unmarshal(responseBytes, &executionResponse) != nil {
			log.Fatal(err)
			return
		}

		bot.Send(message.Sender, executionResponse.Output)

		response.Body.Close()
	})

	bot.Start()
}
