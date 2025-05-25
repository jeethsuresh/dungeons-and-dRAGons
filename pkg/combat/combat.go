package combat

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jeethsuresh/ai/pkg/llm"
	"github.com/mitchellh/mapstructure"
)

type Encounter struct {
	llm        *llm.LLM
	Combatants []Combatant
}

type CombatContent struct {
	Content    string `mapstructure:"content"`
	CombatTurn struct {
		Move CombatTurn `mapstructure:"move"`
	} `mapstructure:"combat_turn"`
	Combatants []Combatant `mapstructure:"combatants"`
	Type       string      `mapstructure:"type"`
}

type Combatant struct {
	Health  int64    `mapstructure:"health"`
	Name    string   `mapstructure:"name"`
	Weapons []string `mapstructure:"weapons"`
	Armor   int64    `mapstructure:"armor"`
	Spells  []string `mapstructure:"spells"`
}

type CombatTurn struct {
	Type   string
	Actor  string `mapstructure:"actor"`
	Target string `mapstructure:"target"`
	Damage int64  `mapstructure:"damage"`
}

var systemPrompt = `
	You are a fair combat system, built to give players a challenging combat encounter while still being winnable. 
	You will be given a list of combatants, and a series of combat turns made by each combatant against the others. 
	Your job is to simulate the encounter's next combat turn, and then return it in JSON format. 
	Do not ever respond in anything other than valid JSON. Double-check the JSON if necessary to ensure compliance. 

	Always respond with a valid combat turn. If the combat turn results in a victory or defeat, the TYPE should be "victory" or "defeat" respectively.
`

var responseFormat = map[string]any{
	"name":   "Generic",
	"strict": "true",
	"schema": map[string]any{
		"type": "object",
		"properties": map[string]any{
			"type": map[string]any{
				"type": "string",
				"enum": []string{"combat", "victory", "defeat"},
			},
			"content": map[string]any{
				"type": "string",
			},
			"combatants": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"name":   map[string]any{"type": "string"},
						"health": map[string]any{"type": "integer", "minimum": 0},
						"weapons": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "string",
							},
						},
						"armor": map[string]any{
							"type":    "integer",
							"minimum": 0,
						},
						"spells": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": "string",
							},
						},
					},
					"required": []string{"name", "health", "weapons", "armor", "spells"},
				},
			},
			"combat_turn": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"move": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"type": map[string]any{
								"type": "string",
							},
							"damage": map[string]any{
								"type": "integer",
							},
							"actor": map[string]any{
								"type": "string",
							},
							"target": map[string]any{
								"type": "string",
							},
						},
						"required": []string{"actor", "target", "damage", "type"},
					},
				},
				"required": []string{"move"},
			},
		},
		"required": []string{"combat_turn", "combatants", "content", "type"},
	},
}

func (e *Encounter) CombatLoop(desc string) {
	e.llm = llm.NewLLM(systemPrompt, responseFormat)
	fmt.Printf("%+v\n", e.ProcessLLMMap(e.llm.NextCombatPrompt(desc)))

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		inp := scanner.Text()
		s := strings.TrimSpace(inp)
		respPrompt := e.ProcessLLMMap(e.llm.NextCombatPrompt(s))
		fmt.Printf("----------------------\n")
		if respPrompt.Content == "" && respPrompt.Type == "combat" {
			fmt.Printf("RespPrompt is NOTHING")
			break
		}
		fmt.Printf("TYPE: %+v\n", respPrompt.Type)
		fmt.Printf("DM SAYS: %+v\n", respPrompt.Content)
		if respPrompt.Type == "victory" {
			fmt.Printf("VICTORY")
			break
		} else if respPrompt.Type == "defeat" {
			fmt.Printf("DEFEAT. YOUR JOURNEY ENDS HERE")
			break
		}
		fmt.Printf("----------------------\n")
	}

}

func (e *Encounter) ProcessLLMMap(llmMap map[string]any) *CombatContent {
	respString := &CombatContent{}
	errDecode := mapstructure.Decode(llmMap, respString)
	if errDecode != nil {
		panic(errDecode)
	}
	e.llm.AddToContext([]llm.Message{{Role: "assistant", Content: respString.Content, Length: len(respString.Content)}})
	return respString
}
