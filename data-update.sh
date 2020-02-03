# !/usr/bin/bash

address="3PG2vMhK5CPqsCDodvLGzQ84QkoHXCJ3oNP"
node_url="https://nodes.wavesplatform.com/addresses/data"
save_endpoint="scriptdata.test.json"

update_data () {
    curl -X GET --header 'Accept: application/json' "$node_url/$address" > "$save_endpoint"
}

main () {
    
    # update_data
    # update_data
}

main $@