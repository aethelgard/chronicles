#!/bin/bash 

: <<'COMMENT'
This script sends a question to all the bot services in the list.
Pre-requisite: All the bot services should be running:
```bash
docker-compose --file ./compose-bots.yml up
```
COMMENT


SERVICES[1]="http://localhost:5051"
SERVICES[2]="http://localhost:5052"
SERVICES[3]="http://localhost:5053"
SERVICES[4]="http://localhost:5054"

read -r -d '' DATA <<- EOM
{
  "question":"What is your name and your kind?"
}
EOM


for key in "${!SERVICES[@]}"; do
    echo ""
    echo "${key}- Sending question: ${DATA} on ${SERVICES[$key]}"
    echo ""

    curl --no-buffer ${SERVICES[$key]}/api/chat \
    -H "Content-Type: application/json" \
    -d "${DATA}" 

    echo ""
done

