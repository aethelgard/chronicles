#!/bin/bash 
SERVICE_URL="http://localhost:5051"
#SERVICE_URL="http://localhost:8080"

read -r -d '' DATA <<- EOM
{
  "question":"tell me the magic keyword"
}
EOM

echo "Sending question: ${DATA} on ${SERVICE_URL}"
echo ""
# --silent

curl --no-buffer ${SERVICE_URL}/api/chat \
    -H "Content-Type: application/json" \
    -d "${DATA}" 

echo ""