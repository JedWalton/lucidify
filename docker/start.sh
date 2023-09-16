#!/bin/sh

trap "exit" INT TERM
trap "kill 0" EXIT

# Start the first ssh command in the background
ssh -R lucidify-api.serveo.net:80:localhost:8080 serveo.net &

# Start the second ssh command in the background
ssh -R lucidify.serveo.net:80:localhost:3001 serveo.net &

docker compose up

# Wait for all background processes to finish
wait

