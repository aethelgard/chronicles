#!/bin/bash
OLLAMA_HOST=http://robby.local:4000 \
LLM_NAME_GENERATION=qwen2.5:0.5b \
LLM_SHEET_GENERATION=qwen2.5:0.5b \
go run main.go hobbit
