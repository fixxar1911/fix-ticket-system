# Prometheus Metrics Documentation

This document describes the Prometheus metrics available in the ticketing system.

## HTTP Metrics

### `http_requests_total`
Counter of total HTTP requests.

**Labels:**
- `method`: HTTP method (GET, POST, PUT, DELETE)
- `endpoint`: API endpoint path
- `status`: HTTP status code

**Example Query:**
```promql
# Total requests in the last 5 minutes
sum(rate(http_requests_total[5m]))

# Requests by endpoint
sum by (endpoint) (rate(http_requests_total[5m]))

# Error rate (status >= 400)
sum(rate(http_requests_total{status=~"4..|5.."}[5m])) / sum(rate(http_requests_total[5m]))
```

### `http_request_duration_seconds`
Histogram of HTTP request durations.

**Labels:**
- `method`: HTTP method
- `endpoint`: API endpoint path

**Example Query:**
```promql
# 95th percentile latency
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))

# Average latency by endpoint
sum(rate(http_request_duration_seconds_sum[5m])) by (endpoint) / sum(rate(http_request_duration_seconds_count[5m])) by (endpoint)
```

## Ticket Metrics

### `ticket_operations_total`
Counter of ticket operations.

**Labels:**
- `operation`: Operation type (create, get, get_all, update, delete)
- `status`: Operation status (success, error)

**Example Query:**
```promql
# Operation success rate
sum(rate(ticket_operations_total{status="success"}[5m])) / sum(rate(ticket_operations_total[5m]))

# Operations by type
sum by (operation) (rate(ticket_operations_total[5m]))
```

### `ticket_status_total`
Gauge of tickets by status.

**Labels:**
- `status`: Ticket status (open, in_progress, resolved, closed)

**Example Query:**
```promql
# Current ticket distribution
ticket_status_total

# Percentage of tickets by status
ticket_status_total / sum(ticket_status_total)
```

## Error Metrics

### `error_total`
Counter of errors by type.

**Labels:**
- `type`: Error type
  - `create_ticket`: Error creating ticket
  - `get_ticket`: Error retrieving ticket
  - `get_all_tickets`: Error retrieving all tickets
  - `update_ticket`: Error updating ticket
  - `delete_ticket`: Error deleting ticket
  - `invalid_input`: Invalid request input
  - `invalid_id`: Invalid ticket ID

**Example Query:**
```promql
# Error rate by type
sum(rate(error_total[5m])) by (type)

# Total errors in the last hour
sum(increase(error_total[1h]))
```

## Recommended Alerts

Here are some recommended alert rules for Prometheus:

```yaml
groups:
- name: ticket_system
  rules:
  # Error Rate Alerts
  - alert: HighErrorRate
    expr: sum(rate(error_total[5m])) > 0.1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: High error rate detected
      description: Error rate is above 10% for the last 5 minutes

  - alert: CriticalErrorRate
    expr: sum(rate(error_total[5m])) > 0.3
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: Critical error rate detected
      description: Error rate is above 30% for the last 2 minutes

  # Latency Alerts
  - alert: HighLatency
    expr: histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le)) > 1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: High latency detected
      description: 95th percentile latency is above 1 second

  - alert: CriticalLatency
    expr: histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le)) > 3
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: Critical latency detected
      description: 95th percentile latency is above 3 seconds

  # Ticket Operation Alerts
  - alert: HighTicketCreationRate
    expr: rate(ticket_operations_total{operation="create"}[5m]) > 10
    for: 5m
    labels:
      severity: info
    annotations:
      summary: High ticket creation rate
      description: More than 10 tickets created per minute

  - alert: TicketOperationErrors
    expr: rate(ticket_operations_total{status="error"}[5m]) > 0
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: Ticket operation errors detected
      description: There are errors in ticket operations

  # Status Distribution Alerts
  - alert: HighOpenTickets
    expr: ticket_status_total{status="open"} > 100
    for: 15m
    labels:
      severity: warning
    annotations:
      summary: High number of open tickets
      description: More than 100 tickets are in open status

  - alert: StaleTickets
    expr: ticket_status_total{status="in_progress"} > 50
    for: 24h
    labels:
      severity: warning
    annotations:
      summary: Stale tickets detected
      description: More than 50 tickets have been in progress for 24 hours

  # System Health Alerts
  - alert: HighRequestRate
    expr: sum(rate(http_requests_total[5m])) > 100
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: High request rate detected
      description: More than 100 requests per second

  - alert: ErrorSpike
    expr: sum(increase(error_total[5m])) > 50
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: Error spike detected
      description: More than 50 errors in the last 5 minutes

  # Endpoint-specific Alerts
  - alert: EndpointHighErrorRate
    expr: sum(rate(http_requests_total{status=~"5.."}[5m])) by (endpoint) / sum(rate(http_requests_total[5m])) by (endpoint) > 0.05
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: High error rate for endpoint
      description: Endpoint {{ $labels.endpoint }} has more than 5% error rate

  - alert: EndpointHighLatency
    expr: histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le, endpoint)) > 2
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: High latency for endpoint
      description: Endpoint {{ $labels.endpoint }} has 95th percentile latency above 2 seconds
```

### Alert Severity Levels

- **Info**: Non-critical alerts for monitoring trends
- **Warning**: Issues that need attention but don't require immediate action
- **Critical**: Issues that require immediate attention

### Alert Response Guidelines

1. **Critical Alerts**
   - Immediate investigation required
   - On-call engineer should be notified
   - Escalate if not resolved within 15 minutes

2. **Warning Alerts**
   - Investigate during business hours
   - Create ticket if issue persists
   - Review during next team meeting

3. **Info Alerts**
   - Monitor trends
   - Review during regular maintenance
   - Adjust thresholds if needed

### Alert Maintenance

1. **Regular Review**
   - Review alert effectiveness monthly
   - Adjust thresholds based on historical data
   - Remove unused alerts

2. **Documentation**
   - Document alert response procedures
   - Keep runbooks up to date
   - Document known issues and solutions

3. **Testing**
   - Test alert delivery channels regularly
   - Verify alert conditions
   - Simulate alert scenarios

## Grafana Dashboard

Recommended Grafana dashboard panels:

1. **System Overview**
   - Request rate
   - Error rate
   - Average response time
   - Active tickets by status

2. **Ticket Operations**
   - Operation success rate
   - Operations by type
   - Ticket status distribution
   - Ticket creation rate

3. **Error Analysis**
   - Error rate by type
   - Error distribution
   - Error trends

4. **Performance**
   - Response time percentiles
   - Request duration by endpoint
   - Database operation latency

## Best Practices

1. **Monitoring**
   - Set up alerts for critical metrics
   - Monitor error rates and response times
   - Track ticket status distribution
   - Watch for unusual patterns

2. **Maintenance**
   - Regularly review and adjust alert thresholds
   - Clean up old metrics if needed
   - Document any changes to metrics

3. **Performance**
   - Keep metric cardinality low
   - Use appropriate metric types
   - Monitor metric storage usage

## Troubleshooting

Common issues and solutions:

1. **High Error Rate**
   - Check application logs
   - Verify database connectivity
   - Review recent changes

2. **High Latency**
   - Check database performance
   - Review slow queries
   - Monitor system resources

3. **Missing Metrics**
   - Verify Prometheus configuration
   - Check application logs
   - Ensure metrics endpoint is accessible 