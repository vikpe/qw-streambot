echo "Starting all components"

(
    trap 'kill 0' SIGINT

    # proxy
    bash scripts/controllers/proxy.sh &

    # controllers
    sleep 0.2 # wait for proxy to start
    bash scripts/controllers/quake_manager.sh &
    bash scripts/controllers/twitch_manager.sh &
    bash scripts/controllers/twitchbot.sh &

    sleep 1 # wait for services to start
    bash scripts/controllers/ezquake.sh

    wait
)

echo "Stopped all components"
