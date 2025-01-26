#!/bin/bash 
SERVICE_URL_DRAGON="http://localhost:5051"
SERVICE_URL_HOBBIT="http://localhost:5052"
SERVICE_URL_WEREWOLF="http://localhost:5053"
SERVICE_URL_DWARF="http://localhost:5054"


read -r -d '' DATA <<- EOM
{
  "question":"What is your name and your kind?"
}
EOM

echo ""
echo "Sending question: ${DATA} on ${SERVICE_URL_DRAGON}"
echo ""

curl --no-buffer ${SERVICE_URL_DRAGON}/api/chat \
    -H "Content-Type: application/json" \
    -d "${DATA}" 
echo ""

echo ""
echo "Sending question: ${DATA} on ${SERVICE_URL_HOBBIT}"
echo ""

curl --no-buffer ${SERVICE_URL_HOBBIT}/api/chat \
    -H "Content-Type: application/json" \
    -d "${DATA}" 
echo ""

echo ""
echo "Sending question: ${DATA} on ${SERVICE_URL_WEREWOLF}"
echo ""

curl --no-buffer ${SERVICE_URL_WEREWOLF}/api/chat \
    -H "Content-Type: application/json" \
    -d "${DATA}" 
echo ""

echo ""
echo "Sending question: ${DATA} on ${SERVICE_URL_DWARF}"
echo ""

curl --no-buffer ${SERVICE_URL_DWARF}/api/chat \
    -H "Content-Type: application/json" \
    -d "${DATA}" 
echo ""