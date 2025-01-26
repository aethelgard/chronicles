#!/bin/bash 
SERVICE_URL="http://localhost:5051"
#SERVICE_URL="http://localhost:8080"

read -r -d '' DATA <<- EOM
{
  "question":"Describe what you are wearing?"
}
EOM

echo "Sending question: ${DATA} on ${SERVICE_URL}"
echo ""
# --silent

curl --no-buffer ${SERVICE_URL}/api/chat \
    -H "Content-Type: application/json" \
    -d "${DATA}" 

echo ""