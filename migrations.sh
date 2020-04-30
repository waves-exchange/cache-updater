#!/bin/bash

init_migration () {
    go run src/migrations/*.go init
}

reset_go_migrations () {
    go run src/migrations/*.go reset
}

run_go_migrations () {
    init_migration
    go run src/migrations/*.go up
}

main () {
    source ~/.bash_profile

    while [ -n "$1" ]
    do
        case "$1" in
            --reset-migration) reset_go_migrations ;;
            --run-migration) run_go_migrations ;;
            --init-migration) init_migration ;;
        esac
        shift;
    done
}

main "$@"