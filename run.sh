#!/bin/bash

SERVICE_NAME="moma-api"
SERVICE_COMMAND="./api"
LOG_FILE="/var/log/${SERVICE_NAME}.log"
PID_FILE="/var/run/${SERVICE_NAME}.pid"

start_service() {
    if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        echo "${SERVICE_NAME} is already running (PID $(cat "$PID_FILE"))"
        exit 0
    fi

    echo "Starting ${SERVICE_NAME}..."
    nohup "$SERVICE_COMMAND" >> "$LOG_FILE" 2>&1 &
    echo $! > "$PID_FILE"
    echo "${SERVICE_NAME} started with PID $(cat "$PID_FILE")"
}

stop_service() {
    if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        echo "Stopping ${SERVICE_NAME}..."
        kill $(cat "$PID_FILE") && rm -f "$PID_FILE"
        echo "${SERVICE_NAME} stopped."
    else
        echo "${SERVICE_NAME} is not running."
    fi
}

status_service() {
    if [ -f "$PID_FILE" ] && kill -0 $(cat "$PID_FILE") 2>/dev/null; then
        echo "${SERVICE_NAME} is running with PID $(cat "$PID_FILE")."
    else
        echo "${SERVICE_NAME} is not running."
    fi
}

restart_service() {
    echo "Restarting ${SERVICE_NAME}..."
    stop_service
    start_service
}

case "$1" in
    start)
        start_service
        ;;
    stop)
        stop_service
        ;;
    restart)
        restart_service
        ;;
    status)
        status_service
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status}"
        exit 1
        ;;
esac