package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Length  int    `json:"-"`
}

type LLM struct {
	Client         http.Client
	Context        []Message
	ContextLength  int
	Model          string
	ResponseFormat map[string]any
	SystemPrompt   string
}

var InitialPrompt = `You are an expert Dungeon Master. 
	You are running a D&D game for your best friends. You are well versed in running campaigns of any length, and you are
	highly motivated to ensure that your friends have a good, but challenging time with the adventure you lay out for them. 
	
	You are highly skilled in the improv art of "yes, and", and you will give your friends the ability to direct the adventure
	where required while still giving them structure and answering questions when they feel confused. You understand that 
	D&D campaigns are a back and forth conversation, and you will not attempt to force your friends into a path they clearly 
	want to avoid. 

	Your responses will either be EXPLORATION or COMBAT responses. Both of these will be expressed as JSON objects. Always
	respond with JSON objects only - never respond with unstructured text or structured text of any type other than JSON.

	If you respond with a COMBAT type response, you must also give the first turn of the combat encounter. 
	Never respond with an EXPLORATION type response if your content contains a battle scene, or an impending battle. 
	
	`

var ResponseFormat = map[string]any{
	"name":   "Generic",
	"strict": "true",
	"schema": map[string]any{
		"type": "object",
		"properties": map[string]any{
			"content": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"type": map[string]any{
						"type": "string",
						"enum": []string{"combat", "exploration"},
					},
					"content": map[string]any{
						"type": "string",
					},
				},
				"required": []any{"content", "type"},
			},
		},
		"required": []any{"content"},
	},
}

func NewLLM(systemPrompt string, responseFormat map[string]any) *LLM {
	return &LLM{
		Client:         http.Client{Timeout: 300 * time.Second},
		Context:        []Message{{Role: "system", Content: systemPrompt, Length: len(InitialPrompt)}},
		ContextLength:  len(InitialPrompt),
		Model:          "hermes-3-llama-3.1-8b",
		ResponseFormat: responseFormat,
	}
}

func (l *LLM) NextPrompt(prompt string) MessageContent {
	endpoint := "http://localhost:1234/v1/chat/completions"
	l.Context = append(l.Context, Message{Role: "user", Content: prompt, Length: len(prompt)})
	reqBody := map[string]any{
		"model":       l.Model,
		"stream":      false,
		"messages":    l.Context,
		"temperature": 0.7,
		"max_tokens":  -1,
		"response_format": map[string]any{
			"type":        "json_schema",
			"json_schema": l.ResponseFormat,
		},
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(reqBody)
	resp, err := l.Client.Post(endpoint, "application/json", buf)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		errUnmarshal := json.NewDecoder(resp.Body).Decode(&errResp)
		if errUnmarshal != nil {
			panic(errUnmarshal)
		}
		panic(fmt.Sprintf("HTTP error: %s", resp.Status))
	}
	respBodyMap := Resp{}
	errBody := json.NewDecoder(resp.Body).Decode(&respBodyMap)
	if errBody != nil {
		panic(errBody)
	}
	// fmt.Printf("RESPBODYMAP: %+v\n", respBodyMap)
	jsonMap := map[string]any{}
	errUnmarshal := json.Unmarshal([]byte(respBodyMap.Choices[0].Message.Content), &jsonMap)
	if errUnmarshal != nil {
		fmt.Printf("ERROR UNMARSHALLING: %+v\n", string(respBodyMap.Choices[0].Message.Content))
		panic(errUnmarshal)
	}
	// fmt.Printf("JSONMAP: %+v\n", jsonMap)
	respString := &MessageContentContainer{}
	errDecode := mapstructure.Decode(jsonMap, respString)
	if errDecode != nil {
		panic(errDecode)
	}
	// fmt.Printf("RESPSTRING: %+v\n", respString)
	if respString.Content.Type == "exploration" {
		l.Context = append(l.Context, Message{Role: "assistant", Content: respString.Content.Content, Length: len(respString.Content.Content)})
	} else {
		fmt.Printf("%+v\n", respString)
	}

	return respString.Content
}

type Resp struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		}
	} `json:"choices"`
}

type MessageContentContainer struct {
	Content MessageContent `mapstructure:"content"`
}

type MessageContent struct {
	Type       string    `mapstructure:"type"`
	Content    string    `mapstructure:"content"`
	Combatants []string  `mapstructure:"combatants"`
	FirstTurn  FirstTurn `mapstructure:"first_turn"`
}

type FirstTurn struct {
	Actor  string `mapstructure:"actor"`
	Name   string `mapstructure:"name"`
	Target string `mapstructure:"target"`
	Damage int    `mapstructure:"damage"`
}

func (l *LLM) NextCombatPrompt(prompt string) map[string]interface{} {
	endpoint := "http://localhost:1234/v1/chat/completions"
	l.Context = append(l.Context, Message{Role: "user", Content: prompt, Length: len(prompt)})
	reqBody := map[string]any{
		"model":       l.Model,
		"stream":      false,
		"messages":    l.Context,
		"temperature": 0.7,
		"max_tokens":  -1,
		"response_format": map[string]any{
			"type":        "json_schema",
			"json_schema": l.ResponseFormat,
		},
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(reqBody)
	resp, err := l.Client.Post(endpoint, "application/json", buf)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		errUnmarshal := json.NewDecoder(resp.Body).Decode(&errResp)
		if errUnmarshal != nil {
			panic(errUnmarshal)
		}
		panic(fmt.Sprintf("HTTP error: %s", resp.Status))
	}
	respBodyMap := Resp{}
	errBody := json.NewDecoder(resp.Body).Decode(&respBodyMap)
	if errBody != nil {
		panic(errBody)
	}
	// fmt.Printf("RESPBODYMAP: %+v\n", respBodyMap)
	jsonMap := map[string]any{}
	errUnmarshal := json.Unmarshal([]byte(respBodyMap.Choices[0].Message.Content), &jsonMap)
	if errUnmarshal != nil {
		fmt.Printf("ERROR UNMARSHALLING: %+v\n", string(respBodyMap.Choices[0].Message.Content))
		panic(errUnmarshal)
	}

	return jsonMap
}

func (l *LLM) AddToContext(ctx []Message) {
	l.Context = append(l.Context, ctx...)
}
