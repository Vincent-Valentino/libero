#!/bin/bash
echo "Starting Libero backend with hot reloading..."
# Change to the directory containing this script
cd "$(dirname "$0")"
air -c .air.toml 