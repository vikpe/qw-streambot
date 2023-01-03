COLOR_GRAY="0;90"    # timestamp
COLOR_GREEN="0;32"   # proxy
COLOR_YELLOW="1;33"  # ezquake
COLOR_MAGENTA="0;35" # chatbot
COLOR_RED="0;31"     # channel manager
COLOR_CYAN="0;36"    # quake_manager
COLOR_BLUE="0;34"    # team viewer
SEQ_RESET="\e[0m"

function text_in_color() {
  local COLOR_CODE=${1}
  local TEXT=${@:2}
  local SEQ_COLOR="\e[${COLOR_CODE}m"
  echo ${SEQ_COLOR}${TEXT}${SEQ_RESET}
}

function pretty_print() {
  local PREFIX_TEXT=$1
  local COLOR_CODE=${2}
  local TEXT=${@:3}
  local TIMESTAMP=$(text_in_color ${COLOR_GRAY} $(date +%T))
  local PREFIX=$(text_in_color ${COLOR_CODE} ${PREFIX_TEXT})
  echo -e "${TIMESTAMP}  ${PREFIX}  ${TEXT}"
}

function run_forever() {
  NAME=${1}
  COLOR=${2}
  CMD=${3}
  RETRY_TIMEOUT=${4}

  while true; do
    pretty_print ${COLOR} ${NAME} "start"
    ${CMD}
    pretty_print ${COLOR} ${NAME} "stopped, restarting in ${RETRY_TIMEOUT} seconds.."
    sleep ${RETRY_TIMEOUT}
  done
}
