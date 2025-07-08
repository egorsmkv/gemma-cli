package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/17twenty/gemma-cli/env"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// DefaultSchema represents the default JSON schema when none is provided
var DefaultSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"message": map[string]interface{}{
			"type":        "string",
			"description": "Response message from Gemini",
		},
	},
	"required": []string{"message"},
}

// Config holds the application configuration
type Config struct {
	APIKey     string
	PromptFile string
	Model      string
	SchemaFile string
	OutputFile string
	InputFile  string
}

func main() {
	// Parse command line flags
	var (
		promptFile = flag.String("prompt", "", "Path to prompt file (required)")
		model      = flag.String("model", "gemini-1.5-flash", "Model to use (default: gemini-1.5-flash)")
		schemaFile = flag.String("schema", "", "Path to JSON schema file (optional)")
		outputFile = flag.String("output", "", "Output file path (default: stdout)")
		inputFile  = flag.String("input", "", "Input file path (required)")
	)
	flag.Parse()

	// Validate required flags
	if *promptFile == "" || *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -prompt=<prompt.txt> -input=<input.txt> [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nRequired flags:\n")
		fmt.Fprintf(os.Stderr, "  -prompt=<file>   Path to prompt file\n")
		fmt.Fprintf(os.Stderr, "  -input=<file>    Path to input file\n")
		fmt.Fprintf(os.Stderr, "\nOptional flags:\n")
		fmt.Fprintf(os.Stderr, "  -model=<model>   Model to use (default: gemini-1.5-flash)\n")
		fmt.Fprintf(os.Stderr, "  -schema=<file>   Path to JSON schema file\n")
		fmt.Fprintf(os.Stderr, "  -output=<file>   Output file path (default: stdout)\n")
		fmt.Fprintf(os.Stderr, "\nEnvironment variables:\n")
		fmt.Fprintf(os.Stderr, "  GEMINI_API_KEY   Google Gemini API key (required)\n")
		os.Exit(1)
	}

	// Get API key from environment
	apiKey := env.GetAsString("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "Error: GEMINI_API_KEY environment variable is required\n")
		os.Exit(1)
	}

	config := Config{
		APIKey:     apiKey,
		PromptFile: *promptFile,
		Model:      *model,
		SchemaFile: *schemaFile,
		OutputFile: *outputFile,
		InputFile:  *inputFile,
	}

	if err := run(config); err != nil {
		slog.Error("Application error", "error", err)
		os.Exit(1)
	}
}

func run(config Config) error {
	// Read prompt file
	promptContent, err := os.ReadFile(config.PromptFile)
	if err != nil {
		return fmt.Errorf("failed to read prompt file: %w", err)
	}

	// Read input file
	inputContent, err := os.ReadFile(config.InputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Load schema
	var schema map[string]interface{}
	if config.SchemaFile != "" {
		schemaContent, err := os.ReadFile(config.SchemaFile)
		if err != nil {
			return fmt.Errorf("failed to read schema file: %w", err)
		}
		if err := json.Unmarshal(schemaContent, &schema); err != nil {
			return fmt.Errorf("failed to parse schema file: %w", err)
		}
	} else {
		schema = DefaultSchema
	}

	// Create Gemini client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}
	defer client.Close()

	// Get the model
	model := client.GenerativeModel(config.Model)

	// Configure the model for JSON output
	model.ResponseMIMEType = "application/json"

	// Set the response schema
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	model.ResponseSchema = &genai.Schema{}
	if err := json.Unmarshal(schemaBytes, model.ResponseSchema); err != nil {
		return fmt.Errorf("failed to set response schema: %w", err)
	}

	// Create the full prompt
	fullPrompt := fmt.Sprintf("%s\n\nInput:\n%s", string(promptContent), string(inputContent))

	// Generate content
	resp, err := model.GenerateContent(ctx, genai.Text(fullPrompt))
	if err != nil {
		return fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract the response
	if len(resp.Candidates) == 0 {
		return fmt.Errorf("no response candidates received")
	}

	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText += string(txt)
		}
	}

	// Parse and format the JSON response
	var jsonResponse interface{}
	if err := json.Unmarshal([]byte(responseText), &jsonResponse); err != nil {
		return fmt.Errorf("failed to parse response as JSON: %w", err)
	}

	// Format with 2-space indentation
	formattedJSON, err := json.MarshalIndent(jsonResponse, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON response: %w", err)
	}

	// Write output
	if config.OutputFile != "" {
		if err := os.WriteFile(config.OutputFile, formattedJSON, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	} else {
		fmt.Println(string(formattedJSON))
	}

	return nil
}
