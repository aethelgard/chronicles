#!/bin/bash
OLLAMA_HOST=http://localhost:11434 \
LLM_NAME_GENERATION=qwen2.5:1.5b \
LLM_SHEET_GENERATION=qwen2.5:3b \
go run main.go werewolf
