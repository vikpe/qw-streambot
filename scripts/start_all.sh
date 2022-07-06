echo "Starting all services"

(
    trap 'kill 0' SIGINT
    bash scripts/proxy_controller.sh &
    sleep 0.2s # wait for proxy

    bash scripts/brain_controller.sh &
    bash scripts/chatbot_controller.sh &
    bash scripts/ezquake_controller.sh

    wait
)

echo "Stopped all services"
