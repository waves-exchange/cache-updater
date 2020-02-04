# !/usr/bin/bash

# inner variables
address="3PG2vMhK5CPqsCDodvLGzQ84QkoHXCJ3oNP"
node_url="https://nodes.wavesplatform.com/addresses/data"
save_endpoint="scriptdata.test.json"

# script params
update_freq=5s
log_file=pg_update_logs.txt
current_pid=?

raw_update_data () {
    curl -X GET --header 'Accept: application/json' "$node_url/$address" > "$save_endpoint"
}

run_go_build () {
    go build && ./cache-updater
}

run_go_build_recursively () {
    printf "$(date +"%D %T") " >> "$log_file"
    ./cache-updater >> "$log_file"

    sleep $update_freq;

    run_go_build_recursively
}

main () {
    # run_go_build
    args=$@

    while [ -n "$1" ]
    do
        case "$1" in
            --frequency) update_freq=$2 ;;
        esac
        shift;
    done

    go build
    run_go_build_recursively
}

main $@