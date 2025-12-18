#!/usr/bin/env bash
set -e

command_exists() {
  command -v "$1" >/dev/null 2>&1
}

echo "Checking prerequisites..."

missing=0

if ! command_exists docker; then
  echo "❌ docker not found"
  missing=1
else
  echo "✅ docker found"
fi

if ! docker compose version >/dev/null 2>&1; then
  echo "❌ docker compose not found"
  missing=1
else
  echo "✅ docker compose found"
fi

if ! command_exists make; then
  echo "❌ make not found"
  missing=1
else
  echo "✅ make found"
fi

if [ "$missing" -eq 0 ]; then
  echo "All prerequisites are installed."
  exit 0
fi

echo
echo "Some tools are missing."
echo
echo "Please install the following manually, whichever were found missing by this script:"
echo "  - Docker (https://docs.docker.com/get-docker/)"
echo "  - Docker Compose"
echo "  - GNU Make"
echo
echo "On linux after installing Docker, you may need to:"
echo "  sudo usermod -aG docker \$USER"
echo "  newgrp docker"

exit 1