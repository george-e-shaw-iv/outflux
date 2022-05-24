#!/usr/bin/env bash

# This script writes to InfluxDB by way of telegraf which is proxying the request
# and writing the data to a file as well in order for outflux to read it and send
# it to the server during sync events.

# Text formatting directive helpers.
bold=$(tput bold)
normal=$(tput sgr0)

numMetrics=$1
for i in $(seq 1 "$numMetrics"); do
    let server_number="$RANDOM%2"
    let cpu_load_int="$RANDOM%100"
    cpu_load=`bc <<< "scale=2; $cpu_load_int/100"`

    echo "${bold}($i/$numMetrics)${normal} Sending cpu_load metric $cpu_load with tag server0$server_number:"
    curl -s -o /dev/null -w "    -> Status Code: %{http_code}\n" -XPOST 'http://localhost:8086/api/v2/write' --data-binary 'cpu_load,host=server0'"$server_number"' value='"$cpu_load"''
done