#!/bin/bash

# inner variables
address="3PG2vMhK5CPqsCDodvLGzQ84QkoHXCJ3oNP"
node_url="https://nodes.wavesplatform.com/addresses/data"
save_endpoint="scriptdata.test.json"

# script params
update_freq=5s
redo_migrate=0
log_file=/var/www/pg_update_logs.txt
current_pid=?

raw_update_data () {
    curl -X GET --header 'Accept: application/json' "$node_url/$address" > "$save_endpoint"
}

run_go_migrations () {
    go run src/migrations/*.go reset
    go run src/migrations/*.go up
}

run_go_build () {
    go build
    if [ $redo_migrate -eq 1 ]
    then
        run_go_migrations
    fi

    ./cache-updater
}

run_go_build_recursively () {
    printf "$(date +"%D %T") " >> "$log_file"
    ./cache-updater >> "$log_file"

    sleep $update_freq;

    run_go_build_recursively
}

main () {
    args=$@
    source ~/.bash_profile

    while [ -n "$1" ]
    do
        case "$1" in
            --frequency) update_freq=$2 ;;
            --redo-migration) redo_migrate=1 ;;
            --pwd ) PWD=$2 ;;
        esac
        shift;
    done

    cd $PWD

    run_go_build
    run_go_build_recursively
}

main "$@"