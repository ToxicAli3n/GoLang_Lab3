#!/bin/bash

send_command() {
    curl -X POST -d "$1" http://localhost:17000
}

x=1
y=0
step=0.01

send_command "green"
send_command "figure $x $y"

while true; do
    send_command "move $x $y"
    x=$(awk "BEGIN {printf \"%.2f\", $x - $step}")
    y=$(awk "BEGIN {printf \"%.2f\", $y + $step}")

    if (( $(awk "BEGIN {print ($x <= 0 && $y >= 1)}") )); then
        break
    fi

    send_command "update"
    sleep 0.1
done