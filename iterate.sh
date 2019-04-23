#!/usr/bin/env bash

function noop() {
  echo "" >/dev/null
}

while true; do
  trap noop SIGINT;
  make run
  trap - SIGINT;
  sleep 1s
done
