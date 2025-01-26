package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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
	BackgroundStory     string `json:"background_story"`
	PersonalityTraits  string `json:"personality_traits"`
	Quote              string `json:"quote"`
}

type Character struct {
	Name  string `json:"name"`
	Kind  string `json:"kind"`
	Sheet Sheet  `json:"sheet"`
}

func GetSheetSchema() ([]byte, error) {

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title": map[string]any{
				"type": "string",
			},
			"age": map[string]any{
				"type": "string",
			},
			"family": map[string]any{
				"type": "string",
			},
			"occupation": map[string]any{
				"type": "string",
			},
			"physical_appearance": map[string]any{
				"type": "string",
			},
			"clothing": map[string]any{
				"type": "string",
			},
			"food_preferences": map[string]any{
				"type": "string",
			},
			"background_story": map[string]any{
				"type": "string",
			},
			"personality_traits": map[string]any{
				"type": "string",
			},
			"quote": map[string]any{
				"type": "string",
			},
		},
		"required": []string{"title", "age", "family", "occupation", "physical_appearance", "clothing", "food_preferences", "background_story", "personality_traits", "quote"},
	}

	jsonSchema, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(jsonSchema), nil
}

// define schema for a structured output
func GetCharacterSchema() ([]byte, error) {
	// define schema for a structured output
	// ref: https://ollama.com/blog/structured-outputs
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type": "string",
			},
			"kind": map[string]any{
				"type": "string",
			},
		},
		"required": []string{"name", "kind"},
	}

	jsonSchema, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(jsonSchema), nil
}

func GetNewCharacter(ctx context.Context, client *api.Client, kind, model string, systemInstructions, generationInstructions []byte) (Character, error) {

	// define schema for a structured output
	// ref: https://ollama.com/blog/structured-outputs
	jsonSchema, err := GetCharacterSchema()
	if err != nil {
		return Character{}, err
	}

	userContent := fmt.Sprintf("Generate a random name for an %s (kind always equals %s).", kind, kind)
	// Prompt construction
	messages := []api.Message{
		{Role: "system", Content: string(systemInstructions)},
		{Role: "system", Content: string(generationInstructions)},
		{Role: "user", Content: userContent},
	}

	//stream := true
	noStream := false

	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
		Options: map[string]interface{}{
			"temperature":    1.7,
			"repeat_last_n":  2,
			"repeat_penalty": 2.2,
			"top_k":          10,
			"top_p":          0.9,
			//"presence_penalty": 1.5,
		},
		Format:    jsonSchema,
		KeepAlive: &api.Duration{Duration: 1 * time.Minute},
		Stream:    &noStream,
	}

	generateName := func() (string, error) {
		jsonResult := ""
		respFunc := func(resp api.ChatResponse) error {
			jsonResult = resp.Message.Content
			return nil
		}
		// Start the chat completion
		err := client.Chat(ctx, req, respFunc)
		if err != nil {
			return jsonResult, err
		}
		return jsonResult, nil
	}

	// Generate a random name
	jsonStr, err := generateName()
	if err != nil {
		return Character{}, err
	}
	character := Character{}

	err = json.Unmarshal([]byte(jsonStr), &character)
	if err != nil {
		return Character{}, err
	}
	//fmt.Println(character.Name, character.Kind)
	character.Kind = kind
	return character, nil
}

func GenerateCharacterSheet(ctx context.Context, client *api.Client, character Character, model string, systemInstructions, generationInstructions []byte) (Sheet, error) {

	jsonSchema, err := GetSheetSchema()
	if err != nil {
		return Sheet{}, err
	}

	userContent := fmt.Sprintf("Using the steps below, create a %s with this name:%s", character.Kind, character.Name)
	// Prompt construction
	messages := []api.Message{
		{Role: "system", Content: string(systemInstructions)},
		{Role: "user", Content: userContent},
		{Role: "user", Content: string(generationInstructions)},
	}

	//stream := true
	noStream := false

	req := &api.ChatRequest{
		Model:    model,
		Messages: messages,
		Options: map[string]interface{}{
			//"temperature":   0.0,
			"temperature":    1.0,
			"repeat_last_n":  2,
			"repeat_penalty": 2.0,
			"top_k":          10,
			"top_p":          0.9,
			//"num_ctx":       4096, // https://github.com/ollama/ollama/blob/main/docs/modelfile.md#valid-parameters-and-values
		},
		Format:    jsonSchema,
		KeepAlive: &api.Duration{Duration: 1 * time.Minute},
		Stream:    &noStream,
	}

	generateSheet := func() (string, error) {
		jsonResult := ""
		respFunc := func(resp api.ChatResponse) error {
			jsonResult = resp.Message.Content
			return nil
		}
		// Start the chat completion
		err := client.Chat(ctx, req, respFunc)
		if err != nil {
			return jsonResult, err
		}
		return jsonResult, nil
	}

	// Generate a random sheet
	jsonStr, err := generateSheet()
	if err != nil {
		return Sheet{}, err
	}
	sheet := Sheet{}

	err = json.Unmarshal([]byte(jsonStr), &sheet)
	if err != nil {
		return Sheet{}, err
	}
	return sheet, nil

}

func main() {

	ctx := context.Background()

	// get the first argument
	// if it is not provided, use the default value
	var kind = "human"
	if len(os.Args) > 1 {
		kind = os.Args[1]
	}

	instructionsPath := os.Getenv("INSTRUCTIONS_PATH")
	if instructionsPath == "" {
		instructionsPath = "./instructions"
	}

	ollamaUrl := os.Getenv("OLLAMA_HOST")
	if ollamaUrl == "" {
		ollamaUrl = "http://localhost:11434"
	}

	modelForNameGeneration := os.Getenv("LLM_NAME_GENERATION")
	if modelForNameGeneration == "" {
		modelForNameGeneration = "qwen2.5:1.5b"
	}

	modelForSheetGeneration := os.Getenv("LLM_SHEET_GENERATION")
	if modelForSheetGeneration == "" {
		modelForSheetGeneration = "qwen2.5:3b"
	}

	fmt.Println("🌍", ollamaUrl, "📕", modelForNameGeneration, "&", modelForSheetGeneration)

	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal("😡:", err)
	}

	systemInstructions, err := os.ReadFile(instructionsPath + "/system.instructions.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	nameGenerationInstructions, err := os.ReadFile(instructionsPath + "/name.generation.instructions.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	sheetGenerationInstructions, err := os.ReadFile(instructionsPath + "/sheet.generation.instructions.md")
	if err != nil {
		log.Fatal("😡:", err)
	}

	// Get the character
	character, err := GetNewCharacter(ctx, client, kind, modelForNameGeneration, systemInstructions, nameGenerationInstructions)
	if err != nil {
		log.Fatal("😡:", err)
	}

	fmt.Println("🧙‍♂️", character.Name, "🧝‍♂️", character.Kind)

	// Generate the character sheet
	character.Sheet, err = GenerateCharacterSheet(ctx, client, character, modelForSheetGeneration, systemInstructions, sheetGenerationInstructions)
	if err != nil {
		log.Fatal("😡:", err)
	}
	fmt.Println("📜", character.Sheet.Title)
	fmt.Println("⏳", character.Sheet.Age)
	fmt.Println("👪", character.Sheet.Family)
	fmt.Println("👨‍🌾", character.Sheet.Occupation)
	fmt.Println("👀", character.Sheet.PhysicalAppearance)
	fmt.Println("👕", character.Sheet.Clothing)
	fmt.Println("🍽", character.Sheet.FoodPreferences)
	fmt.Println("📜", character.Sheet.BackgroundStory)
	fmt.Println("🧠", character.Sheet.PersonalityTraits)
	fmt.Println("💬", character.Sheet.Quote)

	// Marshal the character to JSON
	characterJSON, err := json.MarshalIndent(character, "", "  ")
	if err != nil {
		log.Fatal("😡:", err)
	}

	// Write the JSON to a file
	// Process the character's name
	processedName := strings.Replace(strings.ToLower(character.Name), " ", "_", -1)
	processedName = strings.Replace(processedName, "'", "", -1)

	characterJsonFile := "./data/character-" + character.Kind + "-" + processedName + ".json"
	err = os.WriteFile(characterJsonFile, characterJSON, 0644)
	if err != nil {
		log.Fatal("😡:", err)
	}

	fmt.Println("✅ character persisted to", characterJsonFile)

}
