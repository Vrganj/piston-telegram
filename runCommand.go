package main

import (
	"strings"

	"gopkg.in/tucnak/telebot.v2"
)

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
	languages := []string{
		"awk",
		"bash",
		"brainfuck", "bf",
		"c",
		"cpp", "c++",
		"csharp", "cs", "c#",
		"deno", "denojs", "denots",
		"elixir", "exs",
		"emacs", "elisp", "el",
		"go",
		"haskell", "hs",
		"java",
		"jelly",
		"julia", "jl",
		"kotlin",
		"lua",
		"nasm", "asm",
		"nasm64", "asm64",
		"nim",
		"node", "javascript", "js",
		"perl", "pl",
		"php",
		"python2",
		"python3", "python",
		"paradoc",
		"ruby",
		"rust",
		"swift",
		"typescript", "ts",
	}

	for _, l := range languages {
		if l == language {
			return true
		}
	}

	return false
}
func runCommand(message *telebot.Message, bot *telebot.Bot) {
	text := strings.TrimSpace(message.Text[4:])
	language := getLanguage(text)
	source := getSource(text)

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

	execute(message, bot, language, source)
}
