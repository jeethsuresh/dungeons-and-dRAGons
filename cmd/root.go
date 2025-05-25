package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jeethsuresh/ai/pkg/combat"
	"github.com/jeethsuresh/ai/pkg/llm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dragons",
	Short: "DND RAG MUD",
	Long:  "Text-based DND game that uses local Llama for plots and a structured approach to combat",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Args: %+v\n", args)
		initialPrompt := strings.Join(args, " ")
		fmt.Printf("%+v\n", l.NextPrompt(initialPrompt).Content)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inp := scanner.Text()
			s := strings.TrimSpace(inp)
			respPrompt := l.NextPrompt(s)
			fmt.Printf("----------------------\n")
			if respPrompt.Content == "" && respPrompt.Type == "exploration" {
				fmt.Printf("RespPrompt is NOTHING")
				break
			}
			fmt.Printf("TYPE: %+v\n", respPrompt.Type)
			fmt.Printf("DM SAYS: %+v\n", respPrompt.Content)
			fmt.Printf("----------------------\n")
			if respPrompt.Type == "combat" {
				StartCombat(respPrompt.Content)
			}
		}
	},
}

var l *llm.LLM

func Execute(currLLM *llm.LLM) {
	l = currLLM
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
}

func StartCombat(desc string) {
	fmt.Printf("Started combat\n")
	encounter := &combat.Encounter{}
	encounter.CombatLoop(desc)

}
