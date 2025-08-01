package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/egorsmkv/gemma-cli/env"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

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
	env.LoadFromFile(".env")
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

	// Create Gemini client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}
	defer client.Close()

	// Get the model
	model := client.GenerativeModel(config.Model)

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

	// Write output
	if config.OutputFile != "" {
		if err := os.WriteFile(config.OutputFile, []byte(responseText), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	} else {
		fmt.Println(responseText)
	}

	return nil
}
