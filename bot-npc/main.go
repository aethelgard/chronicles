package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ollama/ollama/api"
)

type Sheet struct {
	Title              string `json:"title"`
	Age                string `json:"age"`
	Family             string `json:"family"`
	Occupation         string `json:"occupation"`
	PhysicalAppearance string `json:"physical_appearance"`
	Clothing           string `json:"clothing"`
	FoodPreferences    string `json:"food_preferences"`
	BackgroundStory    string `json:"background_story"`
	PersonalityTraits  string `json:"personality_traits"`
	Quote              string `json:"quote"`
}

type Character struct {
	Name  string `json:"name"`
	Kind  string `json:"kind"`
	Sheet Sheet  `json:"sheet"`
}

func GetCharacter() (Character, error) {
	var character Character
	character.Name = os.Getenv("CHARACTER_NAME")
	character.Kind = os.Getenv("CHARACTER_KIND")

	if character.Name == "" || character.Kind == "" {
		return character, fmt.Errorf("😡: character name or kind not set")
	}

	return character, nil
}

/*
GetBytesBody returns the body of an HTTP request as a []byte.
  - It takes a pointer to an http.Request as a parameter.
  - It returns a []byte.
*/
func GetBytesBody(request *http.Request) []byte {
	body := make([]byte, request.ContentLength)
	request.Body.Read(body)
	return body
}

func main() {
	ctx := context.Background()

	var characterDataPath = ""
	if len(os.Args) > 1 {
		characterDataPath = os.Args[1]
	} else {
		log.Fatal("😡: character data file does not exist")
	}

	// Read the character data file
	file, err := os.ReadFile(characterDataPath)
	if err != nil {
		log.Fatal("😡: failed to read character data file")
	}

	// Unmarshal the character data
	var character Character
	if err := json.Unmarshal(file, &character); err != nil {
		log.Fatal("😡: failed to unmarshal character data")
	}

	var httpPort = os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	ollamaUrl := os.Getenv("OLLAMA_HOST")
	if ollamaUrl == "" {
		ollamaUrl = "http://localhost:11434"
	}

	modelForChat := os.Getenv("LLM_CHAT")
	if modelForChat == "" {
		modelForChat = "qwen2.5:0.5b"
	}

	// ADDITIONAL_NPC_DATA
	additionalNpcData := os.Getenv("ADDITIONAL_NPC_DATA")

	fmt.Println("🌍", ollamaUrl, "📕", modelForChat)

	client, errCli := api.ClientFromEnvironment()
	if errCli != nil {
		log.Fatal("😡:", errCli)
	}

	fmt.Println("🧙‍♂️", character.Name, "🧝‍♂️", character.Kind)

	// Context

	contextOfTheCharacter := "CONTEXT: \n"
	contextOfTheCharacter += "Name: " + character.Name + "\n"
	contextOfTheCharacter += "Kind: " + character.Kind + "\n"
	contextOfTheCharacter += "Title: " + character.Sheet.Title + "\n"
	contextOfTheCharacter += "Age: " + character.Sheet.Age + "\n"
	contextOfTheCharacter += "Family: " + character.Sheet.Family + "\n"
	contextOfTheCharacter += "Occupation: " + character.Sheet.Occupation + "\n"
	contextOfTheCharacter += "Physical Appearance: " + character.Sheet.PhysicalAppearance + "\n"
	contextOfTheCharacter += "Clothing: " + character.Sheet.Clothing + "\n"
	contextOfTheCharacter += "Food Preferences: " + character.Sheet.FoodPreferences + "\n"
	contextOfTheCharacter += "Background Story: " + character.Sheet.BackgroundStory + "\n"
	contextOfTheCharacter += "Personality Traits: " + character.Sheet.PersonalityTraits + "\n"
	contextOfTheCharacter += "Quote: " + character.Sheet.Quote + "\n"

	if additionalNpcData != "" {
		contextOfTheCharacter += additionalNpcData + "\n"
		fmt.Println("🧠", additionalNpcData)
	}

	systemContentTpl := `You are a %s, your name is %s,
	expert at interpreting and answering questions based on provided sources.
	Using only the provided context, answer the user's question 
	to the best of your ability using only the resources provided. 
	Be verbose!`

	systemInstructions := fmt.Sprintf(systemContentTpl, character.Kind, character.Name)

	// 🧠 Memory
	memory := []api.Message{
		{Role: "system", Content: "CONVERSATION MEMORY:"},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/chat", func(response http.ResponseWriter, request *http.Request) {

		// add a flusher
		flusher, ok := response.(http.Flusher)
		if !ok {
			response.Write([]byte("😡 Error: expected http.ResponseWriter to be an http.Flusher"))
		}
		body := GetBytesBody(request)
		// unmarshal the json data
		var data map[string]string

		err := json.Unmarshal(body, &data)
		if err != nil {
			response.Write([]byte("😡 Error: " + err.Error()))
		}

		userContent := data["question"]

		// Prompt construction
		messages := []api.Message{
			{Role: "system", Content: contextOfTheCharacter},
			{Role: "system", Content: systemInstructions},
		}

		// 🧠 Add memory
		messages = append(messages, memory...)
		// Add the new user question
		messages = append(messages, api.Message{Role: "user", Content: userContent})

		stream := true
		//noStream  := false

		// Configuration
		options := map[string]interface{}{
			"temperature":      0.8,
			"repeat_last_n":    2,
			"top_k":            10,
			"top_p":            0.9,
			"presence_penalty": 1.5,
		}

		req := &api.ChatRequest{
			Model:     modelForChat,
			Messages:  messages,
			Options:   options,
			KeepAlive: &api.Duration{Duration: 1 * time.Minute},
			Stream:    &stream,
		}

		answer := ""
		respFunc := func(resp api.ChatResponse) error {

			response.Write([]byte(resp.Message.Content))
			fmt.Print(resp.Message.Content)
			answer += resp.Message.Content

			flusher.Flush()

			return nil
		}

		err = client.Chat(ctx, req, respFunc)

		// 🧠 Save the conversation in memory
		memory = append(
			memory,
			api.Message{Role: "user", Content: userContent},
			api.Message{Role: "assistant", Content: answer},
		)

		if err != nil {
			log.Fatal("😡:", err)
		}

	})

	var errListening error
	log.Println("🌍 http server is listening on: " + httpPort)
	errListening = http.ListenAndServe(":"+httpPort, mux)

	log.Fatal(errListening)

}
