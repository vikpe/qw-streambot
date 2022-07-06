echo "Starting all components"

(
    trap 'kill 0' SIGINT
    bash scripts/conrollers/proxy.sh &
    sleep 0.2 # wait for proxy

    bash scripts/conrollers/brain.sh &
    bash scripts/conrollers/chatbot.sh &
    bash scripts/conrollers/ezquake.sh

    wait
)

echo "Stopped all components"
