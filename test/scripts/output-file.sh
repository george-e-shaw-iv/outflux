#!/usr/bin/env bash

# Either dumps or tails the telegraf output file that telegraf dual-writes to
# along with the InfluxDB client instance. The manner in which it determines
# whether to dump or tail is by the first positional argument given to this
# script, either "dump" or "tail". If no argument is passed, it defaults to
# "dump".

action="dump"
if [[ "$#" -eq 1 ]]; then
    action="$1"
    if [[ "$action" != "tail" ]] && [[ "$action" != "dump" ]]; then
        echo "Invalid first positional argument passed, must be either 'tail' or 'dump'"
        exit 1
    fi
fi

case "$action" in
    "dump")
        docker exec $(docker ps -aqf "name=test-telegraf-client-1") cat /etc/telegraf/metrics.out
        ;;
    "tail")
        docker exec $(docker ps -aqf "name=test-telegraf-client-1") tail -f /etc/telegraf/metrics.out
        ;;
    *)
        docker exec $(docker ps -aqf "name=test-telegraf-client-1") cat /etc/telegraf/metrics.out
        ;;
esac
