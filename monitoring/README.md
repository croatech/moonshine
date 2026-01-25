# Monitoring Stack

Full observability stack: Prometheus + Grafana + Loki + cAdvisor

## Services

- **Prometheus** (port 9090): Metrics collection
- **Grafana** (port 3000): Dashboards and visualization
- **Loki** (port 3100): Log aggregation
- **Promtail**: Log shipper
- **cAdvisor** (port 8088): Container metrics

## Quick Start

```bash
# Start all services
docker-compose up -d

# Access Grafana
open http://localhost:3000
# Login: admin / admin

# Access Prometheus
open http://localhost:9090

# Check API metrics
curl http://localhost:1666/metrics
```

## Available Metrics

### HTTP Metrics:
- `http_requests_total` - Total requests by method, path, status
- `http_request_duration_seconds` - Request latency histogram

### Application Metrics:
- `moonshine_fights_total` - Total fights started
- `moonshine_fight_duration_seconds` - Fight duration histogram
- `moonshine_players_online` - Current online players
- `websocket_connections_active` - Active WebSocket connections

### System Metrics (cAdvisor):
- `container_cpu_usage_seconds_total` - CPU usage
- `container_memory_usage_bytes` - RAM usage
- `container_network_receive_bytes_total` - Network RX
- `container_network_transmit_bytes_total` - Network TX

## Grafana Dashboards

### 1. API Performance
- Request rate per endpoint
- Latency percentiles (p50, p95, p99)
- Error rate
- Slowest endpoints

### 2. System Resources
- CPU usage by container
- Memory usage by container
- Network I/O
- Disk I/O

### 3. Game Metrics
- Players online
- Fights per minute
- Average fight duration
- WebSocket connections

### 4. Logs (Loki)
- Container logs
- Error logs
- Slow query logs

## PromQL Examples

### Slow endpoints (>1s):
```promql
histogram_quantile(0.95, 
  rate(http_request_duration_seconds_bucket[5m])
) > 1
```

### Request rate:
```promql
rate(http_requests_total[5m])
```

### Error rate:
```promql
rate(http_requests_total{status=~"5.."}[5m])
```

### Top 5 slowest endpoints:
```promql
topk(5, 
  histogram_quantile(0.95, 
    rate(http_request_duration_seconds_bucket[5m])
  )
)
```

### CPU usage:
```promql
rate(container_cpu_usage_seconds_total{name="moonshine-postgres-1"}[5m]) * 100
```

### Memory usage:
```promql
container_memory_usage_bytes{name="moonshine-postgres-1"} / 1024 / 1024 / 1024
```

## Alerts (TODO)

Create `monitoring/alerts.yml`:

```yaml
groups:
  - name: moonshine
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        annotations:
          summary: "High error rate detected"

      - alert: SlowEndpoint
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 2
        for: 10m
        annotations:
          summary: "Endpoint is slow"

      - alert: HighMemoryUsage
        expr: container_memory_usage_bytes / container_spec_memory_limit_bytes > 0.9
        for: 5m
        annotations:
          summary: "Container memory usage > 90%"
```

## Troubleshooting

### Prometheus not scraping API:
Check `http://localhost:9090/targets` - should see `moonshine-api` as UP

### No metrics in Grafana:
1. Check Prometheus is running: `docker-compose logs prometheus`
2. Verify datasource: Grafana → Configuration → Data Sources
3. Test query: `up{job="moonshine-api"}`

### No logs in Loki:
1. Check Promtail: `docker-compose logs promtail`
2. Verify Loki datasource in Grafana
3. Query: `{container="moonshine-postgres-1"}`

## Production Recommendations

1. **Data Retention**: Set in `prometheus.yml`
   ```yaml
   global:
     scrape_interval: 15s
   storage:
     tsdb:
       retention.time: 30d
   ```

2. **Alerting**: Add Alertmanager for Slack/email notifications

3. **Security**: Enable auth for Prometheus/Grafana

4. **Backup**: Backup Grafana dashboards and Prometheus data
