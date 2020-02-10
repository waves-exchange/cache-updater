# !/usr/sh

# /etc/systemd/system/mydaemon.service

# first arg - service name
service_name="neutrino-cache-daemon"
daemon_fl="data-update.sh"

if [ -n "$1" ]
then
  service_name=$1
fi

build_only () {
  file="$service_name.service"
  echo "
    [Unit]
    Description=Neutrino cache update daemon

    [Service]
    ExecStart=$pwd/$daemon_fl
    Restart=on-failure

    [Install]
    WantedBy=multi-user.target
  "

  mv "$file" /etc/systemd/system/
}

build_n_start () {
  build_only
  start_service
}
start_service () {
  systemctl start "$service_name.service"
}

main () {
  while [ -n "$1" ]
  do
    case "$1" in
      --build-only ) build_only ;;
      --build-n-start ) build_n_start ;;
      --start ) start_service ;;
    esac
    shift
  done
}