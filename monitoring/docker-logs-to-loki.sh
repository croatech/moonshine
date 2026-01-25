#!/bin/sh

LOKI_URL="http://loki:3100/loki/api/v1/push"
CONTAINERS="moonshine-postgres-1 moonshine-clickhouse-1 moonshine-prometheus-1"

for container in $CONTAINERS; do
  docker logs --follow --tail 100 $container 2>&1 | while read line; do
    timestamp=$(date -u +"%s%N")
    json=$(cat <<EOF
{
  "streams": [
    {
      "stream": {
        "container": "$container",
        "source": "docker"
      },
      "values": [
        ["$timestamp", "$line"]
      ]
    }
  ]
}
EOF
)
    curl -X POST "$LOKI_URL" -H "Content-Type: application/json" -d "$json" > /dev/null 2>&1
  done &
done

wait
