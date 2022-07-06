echo "Starting all components"

(
    trap 'kill 0' SIGINT
    bash scripts/controllers/proxy.sh &
    sleep 0.2 # wait for proxy

    bash scripts/controllers/brain.sh &
    bash scripts/controllers/chatbot.sh &
    bash scripts/controllers/ezquake.sh

    wait
)

echo "Stopped all components"
