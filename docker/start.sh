#!/bin/sh

trap "exit" INT TERM
trap "kill 0" EXIT

# Run ngrok in non-interactive mode and redirect output to a log file
nohup ngrok http --domain=proper-immune-sawfish.ngrok-free.app 8080 > ngrok.log 2>&1 &
# nohup ngrok http --domain=proper-immune-sawfish.ngrok-free.app 3002 > ngrok.log 2>&1 &

docker compose up

# Wait for all background processes to finish
wait

