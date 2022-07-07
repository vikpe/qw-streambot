echo "Starting all components"

(
    trap 'kill 0' SIGINT
    bash scripts/controllers/proxy.sh &
    sleep 0.2 # wait for proxy

    bash scripts/controllers/chatbot.sh &
    bash scripts/controllers/brain.sh &

    sleep 1 # wait brain and chatbot
    bash scripts/controllers/ezquake.sh

    wait
)

echo "Stopped all components"
