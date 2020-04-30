#!/bin/bash

# script params
update_freq=5s
redo_migrate=0
log_file=/var/www/pg_update_logs.txt

init_migration () {
    go run src/migrations/*.go init
}

reset_go_migrations () {
    go run src/migrations/*.go reset
}

run_go_migrations () {
    reset_go_migrations
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

main () {
    source ~/.bash_profile

    while [ -n "$1" ]
    do
        case "$1" in
            --frequency) update_freq=$2 ;;
            --reset-migration) reset_go_migrations ;;
            --run-migration) run_go_migrations ;;
            --init-migration) init_migration ;;
            --redo-migration) redo_migrate=1 ;;
            --log-file) log_file=$2 ;;
        esac
        shift;
    done
}

main "$@"