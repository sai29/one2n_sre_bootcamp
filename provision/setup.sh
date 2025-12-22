#!/usr/bin/env bash
set -e

echo 'Acquire::ForceIPv4 "true";' > /etc/apt/apt.conf.d/99force-ipv4

install_docker() {
  apt-get update -y
  apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

  if ! command -v docker >/dev/null 2>&1; then
    curl -4 -fsSL https://get.docker.com | sh
  fi
}

install_docker_compose() {
  if ! docker compose version >/dev/null 2>&1; then
    mkdir -p /usr/local/lib/docker/cli-plugins
    curl -SL https://github.com/docker/compose/releases/download/v2.27.0/docker-compose-linux-x86_64 \
      -o /usr/local/lib/docker/cli-plugins/docker-compose
    chmod +x /usr/local/lib/docker/cli-plugins/docker-compose
  fi
}

setup_permissions() {
  usermod -aG docker vagrant
}

main() {
  install_docker
  install_docker_compose
  setup_permissions
}

main
