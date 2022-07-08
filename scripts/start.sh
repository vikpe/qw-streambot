echo "Starting all components"

(
    trap 'kill 0' SIGINT

    # proxy
    bash scripts/controllers/proxy.sh &

    # controllers
    sleep 0.2 # wait for proxy to start
    bash scripts/controllers/channel_manager.sh &
    bash scripts/controllers/chatbot.sh &
    bash scripts/controllers/brain.sh &

    sleep 1 # wait for controllers to start
    bash scripts/controllers/ezquake.sh

    wait
)

echo "Stopped all components"
