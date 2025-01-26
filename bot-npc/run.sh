#!/bin/bash
OLLAMA_HOST=http://localhost:11434 \
LLM_CHAT=qwen2.5:1.5b \
ADDITIONAL_NPC_DATA="Magic keyword: yolo" \
go run main.go ./data/character-werewolf-thornwood.json
