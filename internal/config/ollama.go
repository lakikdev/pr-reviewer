package config

import "github.com/namsral/flag"

var OllamaBaseURL = flag.String("ollama-base-url", "http://192.168.1.177:11434", "ollama base URL.")
var OllamaModel = flag.String("ollama-model", "qwen2.5-coder:7b", "ollama model name.")
