# Gemma CLI

A command-line interface for interacting with Google's Gemini AI models. This tool allows you to send structured prompts with input data and receive JSON-formatted responses based on custom schemas.

## Features

- Integration with Google Gemini AI models
- Custom prompt support from files
- Configurable JSON response schemas
- Structured JSON output with custom formatting
- Default schema for simple use cases
- File-based input and output support
- Stdout output option

## Prerequisites

- Go 1.23.6 or later
- Google Cloud Project with Gemini API enabled
- Google Gemini API key

## Installation

1. Clone the repository:
```bash
git clone https://github.com/17twenty/gemma-cli.git
cd gemma-cli
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the application:
```bash
go build -o gemma-cli
```

## Configuration

Set your Google Gemini API key as an environment variable:

```bash
export GEMINI_API_KEY="your_api_key_here"
```

Alternatively, create a `.env` file in the project root:
```
GEMINI_API_KEY=your_api_key_here
```

## Usage

### Basic Syntax

```bash
./gemma-cli -prompt=<prompt_file> -input=<input_file> [options]
```

### Required Flags

- `-prompt=<file>`: Path to the prompt file containing instructions for the AI
- `-input=<file>`: Path to the input file containing data to be processed

### Optional Flags

- `-model=<model>`: Gemini model to use (default: `gemini-1.5-flash`)
- `-schema=<file>`: Path to JSON schema file for structured output
- `-output=<file>`: Output file path (default: stdout)

### Available Models

- `gemini-1.5-flash` (default)
- `gemini-1.5-pro`
- `gemini-1.0-pro`

## Examples

### Example 1: Basic Text Analysis

Analyze a text document and get a simple response:

```bash
./gemma-cli -prompt=examples/prompt.txt -input=examples/input.txt
```

**Output:**
```json
{
  "message": "The text discusses TechCorp Inc.'s strong Q3 2024 earnings with 35% revenue growth, reaching $2.5B. Key highlights include AI investment success, 28% customer growth, European expansion, and positive market outlook with stock price increase to $128."
}
```

### Example 2: Structured Analysis with Custom Schema

Use a custom schema for detailed analysis:

```bash
./gemma-cli -prompt=examples/prompt.txt -input=examples/input.txt -schema=examples/schema.json
```

**Output:**
```json
{
  "summary": "TechCorp Inc. reports strong Q3 2024 earnings with significant growth in revenue, customer base, and market expansion, driven by AI investments and cloud computing services.",
  "main_topics": [
    "quarterly earnings",
    "revenue growth",
    "artificial intelligence",
    "cloud computing",
    "market expansion",
    "stock performance"
  ],
  "entities": {
    "people": [
      "Sarah Johnson",
      "Michael Chen"
    ],
    "organizations": [
      "TechCorp Inc.",
      "Deutsche Bank"
    ],
    "locations": [
      "San Francisco",
      "California",
      "European markets"
    ]
  },
  "key_facts": [
    {
      "fact": "Q3 2024 revenue",
      "value": "$2.5 billion"
    },
    {
      "fact": "Revenue growth",
      "value": "35% increase"
    },
    {
      "fact": "Customer base growth",
      "value": "28% to 50,000 users"
    },
    {
      "fact": "R&D spending",
      "value": "$300 million"
    },
    {
      "fact": "Employee headcount",
      "value": "5,500"
    },
    {
      "fact": "Stock price",
      "value": "$128"
    }
  ],
  "sentiment": "positive",
  "confidence_score": 0.95
}
```

### Example 3: Save Output to File

Save the analysis to a file:

```bash
./gemma-cli -prompt=examples/prompt.txt -input=examples/input.txt -schema=examples/schema.json -output=analysis_result.json
```

### Example 4: Using Different Models

Use the more powerful Gemini Pro model:

```bash
./gemma-cli -prompt=examples/prompt.txt -input=examples/input.txt -model=gemini-1.5-pro -schema=examples/schema.json
```

## File Examples

### Prompt File (`examples/prompt.txt`)
```
You are a helpful AI assistant that analyzes text and provides structured responses.

Please analyze the provided input text and extract key information. Focus on identifying:
- Main topics or themes
- Important entities (people, places, organizations)
- Key facts or data points
- Overall sentiment or tone

Provide your analysis in the requested JSON format.
```

### Input File (`examples/input.txt`)
```
The quarterly earnings report for TechCorp Inc. shows remarkable growth in the third quarter of 2024. The company, headquartered in San Francisco, California, reported revenue of $2.5 billion, representing a 35% increase compared to the same quarter last year.

CEO Sarah Johnson expressed optimism about the company's future, stating, "Our investment in artificial intelligence and cloud computing services has paid off significantly. We've seen strong adoption of our new AI-powered analytics platform across enterprise clients."

The report highlights several key achievements:
- Customer base grew by 28% to reach 50,000 active users
- International expansion into European markets contributed 15% of total revenue
- R&D spending increased to $300 million, focusing on machine learning capabilities
- Employee headcount reached 5,500 with plans to hire 1,000 more by year-end

Market analysts remain bullish on TechCorp's prospects, with Deutsche Bank upgrading their rating to "Buy" and setting a price target of $150 per share. The stock closed at $128 after the earnings announcement, up 8% from the previous day.

Looking ahead, the company plans to launch three new products in Q4 2024, including an advanced cybersecurity suite and two cloud-based collaboration tools. CFO Michael Chen noted that the company maintains a strong cash position of $1.2 billion to fund these initiatives.
```

### Schema File (`examples/schema.json`)
```json
{
  "type": "object",
  "properties": {
    "summary": {
      "type": "string",
      "description": "Brief summary of the analyzed content"
    },
    "main_topics": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "description": "List of main topics or themes identified"
    },
    "entities": {
      "type": "object",
      "properties": {
        "people": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "description": "Names of people mentioned"
        },
        "organizations": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "description": "Names of organizations mentioned"
        },
        "locations": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "description": "Names of locations mentioned"
        }
      },
      "required": ["people", "organizations", "locations"]
    },
    "key_facts": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "fact": {
            "type": "string",
            "description": "The key fact or data point"
          },
          "value": {
            "type": "string",
            "description": "Associated value or metric if applicable"
          }
        },
        "required": ["fact"]
      },
      "description": "Important facts or data points extracted"
    },
    "sentiment": {
      "type": "string",
      "enum": ["positive", "negative", "neutral"],
      "description": "Overall sentiment or tone of the content"
    },
    "confidence_score": {
      "type": "number",
      "minimum": 0,
      "maximum": 1,
      "description": "Confidence score for the analysis (0-1)"
    }
  },
  "required": ["summary", "main_topics", "entities", "key_facts", "sentiment", "confidence_score"]
}
```

## Default Schema

When no schema file is provided, the tool uses a default schema:

```json
{
  "type": "object",
  "properties": {
    "message": {
      "type": "string",
      "description": "Response message from Gemini"
    }
  },
  "required": ["message"]
}
```

## Error Handling

The tool will exit with an error message if:
- Required flags (`-prompt` or `-input`) are missing
- The `GEMINI_API_KEY` environment variable is not set
- Input files cannot be read
- The JSON schema is invalid
- The Gemini API returns an error
- The response cannot be parsed as JSON

## Common Use Cases

1. **Document Analysis**: Extract key information from reports, articles, or documents
2. **Data Processing**: Structure unstructured text data into JSON format
3. **Content Classification**: Categorize and analyze text content
4. **Information Extraction**: Pull specific entities and facts from text
5. **Sentiment Analysis**: Determine the tone and sentiment of text content

## Troubleshooting

### API Key Issues
- Ensure your `GEMINI_API_KEY` is valid and has the necessary permissions
- Check that the Gemini API is enabled in your Google Cloud project

### Model Errors
- Verify the model name is correct (case-sensitive)
- Some models may have different availability or rate limits

### Schema Validation
- Ensure your JSON schema is valid JSON
- Check that the schema follows the JSON Schema specification

### File Permissions
- Verify that input files exist and are readable
- Ensure the output directory is writable

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:
- Create an issue on GitHub
- Check the troubleshooting section above
- Review the Google Gemini API documentation
