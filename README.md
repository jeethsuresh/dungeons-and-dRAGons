# dungeons n dRAGons 

This is a quick hackathon project that uses a local AI (LM Studio on my end, though any OpenAI compatible API will work - probably including OpenAI's itself with some minor tweaks for the API key and whatnot) to run a campaign. 

Doesn't work GREAT on my machine because I'm only running an 8b llama model, but beefier systems should get quite a lot out of it. 

Uses two conversations: one for combat, one for exploration; may have more specialized ones in the future. 

## Future work

- Better combat loop
- Proper RAG - inventory, spells/weapons/armor, quests, etc.
- Better prompt engineering, including prompt enhancement to minimize memory loss from the rolling context window

## Usage

Run `go run main.go <INITIAL SETTING/PROMPT>` with LM Studio running in the background (localhost:1234)

e.g. I've been running: `go run main.go Im an orc warrior stuck in a forest, trying to find my way home to my village` 

