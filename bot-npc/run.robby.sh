#!/bin/bash
OLLAMA_HOST=http://robby.local:4000 \
LLM_CHAT=qwen2.5:0.5b \
ADDITIONAL_NPC_DATA="Magic keyword: yolo" \
go run main.go ./data/character-werewolf-thornwood.json
