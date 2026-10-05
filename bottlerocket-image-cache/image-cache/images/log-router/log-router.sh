#!/bin/sh
# Forward everything under LOG_DIR to stdout, once, then exit.
# A real router would tail and ship to a log backend.
set -eu

LOG_DIR="${LOG_DIR:-/var/log/workload}"
echo "log-router: forwarding ${LOG_DIR}"
if [ -d "${LOG_DIR}" ]; then
  find "${LOG_DIR}" -type f -exec sh -c 'echo "--- $1"; cat "$1"' _ {} \;
else
  echo "log-router: ${LOG_DIR} does not exist, nothing to forward"
fi
