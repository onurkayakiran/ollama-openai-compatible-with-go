# OpenAI Compatible API

This project provides a REST API that is compatible with the OpenAI API format, powered by Ollama. Built using the Go Fiber framework.

## Features

- ✅ Full compatibility with OpenAI API format
- ✅ Chat Completions endpoint (`/v1/chat/completions`)
- ✅ Text Completions endpoint (`/v1/completions`)
- ✅ Models endpoint (`/v1/models`)
- ✅ Streaming response support (Server-Sent Events)
- ✅ API Key authentication
- ✅ CORS support
- ✅ Error handling and logging
- ✅ Environment variables configuration via .env file

## Prerequisites

- Go 1.21 or higher
- Ollama installed and running

## Installation

1. Clone the repository:
```bash
git clone <repo-url>
cd openai-compatible
```

2. Install dependencies:
```bash
go mod tidy
```

3. Create your environment configuration:
```bash
cp .env.example .env
```

4. Edit the `.env` file with your configuration:
```bash
# Port configuration
PORT=8080

# API Key for authentication
API_KEY=your-secret-api-key-here

# Ollama server configuration
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=llama3.2:latest
```

5. Make sure Ollama is running:
```bash
ollama serve
```

6. Run the application:
```bash
go run main.go
```

## Configuration

The application uses environment variables for configuration. You can set these variables in the `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | 8080 |
| `API_KEY` | API key for authentication | sk-your-secret-api-key-here |
| `OLLAMA_URL` | Ollama server URL | http://localhost:11434 |
| `OLLAMA_MODEL` | Model to use | llama3.2:latest |

**Note:** The `.env` file is excluded from version control via `.gitignore` for security reasons. Always use `.env.example` as a template.

## API Usage

### Chat Completions

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-secret-api-key-here" \
  -d '{
    "model": "llama3.2:latest",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  }'
```

### Streaming Chat Completions

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-secret-api-key-here" \
  -d '{
    "model": "llama3.2:latest",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ],
    "stream": true
  }'
```

### Text Completions

```bash
curl -X POST http://localhost:8080/v1/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-secret-api-key-here" \
  -d '{
    "model": "llama3.2:latest",
    "prompt": "Hello world",
    "max_tokens": 100
  }'
```

### List Models

```bash
curl -X GET http://localhost:8080/v1/models \
  -H "Authorization: Bearer your-secret-api-key-here"
```

### Health Check

```bash
curl -X GET http://localhost:8080/health
```

## Project Structure

```
openai-compatible/
├── main.go                 # Main application
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies checksums
├── .env                    # Environment variables (not in git)
├── .env.example            # Example environment variables
├── .gitignore             # Git ignore file
├── config/
│   └── config.go          # Configuration management
├── handlers/
│   ├── chat.go            # Chat completions handler
│   ├── completions.go     # Text completions handler
│   └── models.go          # Models handler
├── middleware/
│   └── auth.go            # Authentication middleware
├── models/
│   ├── openai.go          # OpenAI API structures
│   └── ollama.go          # Ollama API structures
└── services/
    └── ollama.go          # Ollama service integration
```

## Supported Parameters

### Chat Completions
- `model`: Model name
- `messages`: Array of messages
- `max_tokens`: Maximum number of tokens
- `temperature`: Creativity level (0.0-2.0)
- `top_p`: Nucleus sampling
- `stream`: Enable/disable streaming
- `stop`: Stop sequences
- `presence_penalty`: Presence penalty
- `frequency_penalty`: Frequency penalty

### Text Completions
- `model`: Model name
- `prompt`: Text prompt
- `max_tokens`: Maximum number of tokens
- `temperature`: Creativity level
- `top_p`: Nucleus sampling
- `stop`: Stop sequences

## Security

- API keys are validated through the authentication middleware
- The `.env` file containing sensitive information is excluded from version control
- Always use strong, unique API keys in production environments

## Building for Production

```bash
# Build the binary
go build -o openai-compatible

# Run the binary
./openai-compatible
```

## License

MIT License

---

[🇹🇷 Türkçe README için tıklayın](README.tr.md)
