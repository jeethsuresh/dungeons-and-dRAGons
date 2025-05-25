package main

import (
	"github.com/jeethsuresh/ai/cmd"
	"github.com/jeethsuresh/ai/pkg/llm"
)

func main() {
	llm := llm.NewLLM(llm.InitialPrompt, llm.ResponseFormat)
	cmd.Execute(llm)
}
