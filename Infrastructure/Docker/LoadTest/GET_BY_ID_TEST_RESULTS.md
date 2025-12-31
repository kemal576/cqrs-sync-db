# Get By ID Load Test Results

This document contains load test results for the ProductRead GetById endpoint in the CQRS CDC project.

## Test Environment

- **Cluster:** Minikube (Docker Driver)
- **Tool:** Locust
- **Test Script:** GetById.py
- **Endpoint:** GET /products/{id}
- **Cache Strategy:** Redis Sentinel with cache-aside pattern
- **Initial Scale:** All services run with **1 replica**

### Pod Resource Specifications

| Service | CPU Request | CPU Limit | Memory Request | Memory Limit | Replicas (Initial) |
|--------|-------------|-----------|----------------|--------------|-------------------|
| ProductRead | 200m | 500m | 256Mi | 512Mi | 1 |
| ProductWrite | 200m | 500m | 256Mi | 512Mi | 2 |
| ProductConsumer | 100m | 500m | 128Mi | 256Mi | 1 |
| MongoDB | - | - | - | - | 1 |
| Kafka | - | - | - | - | 1 |
| Elasticsearch | - | - | - | - | 1 |
| Redis Sentinel | - | - | - | - | 3 |

**Note:** MongoDB, Kafka, Elasticsearch, and Redis have no resource limits defined.

## ProductRead - Get By ID Endpoint

### Test Results

| Date | Users | Spawn Rate | Duration | Total Requests | RPS | Avg Response (ms) | Min (ms) | Max (ms) | 95% (ms) | 99% (ms) | Success Rate | Failure Rate | Cache Hit Rate | Notes |
|-------|-------|------------|----------|----------------|-----|-------------------|----------|----------|----------|----------|--------------|--------------|----------------|-------|
| 31.12.2025 | ~50 | - | 54s | 4617 | 85.54 | 569 | 3 | 9093 | 2000 | 3400 | 100% | 0% | ~2.8% | Cold cache - mostly ES queries, cache expired after test (TTL) |
| 31.12.2025 | ~150 | - | 2 min | 31963 | 266.37 | 45.49 | 2 | 1086 | 210 | 450 | 100% | 0% | ~3% | Warm cache but low hit rate - 12x faster than cold cache! |

### Column Descriptions

- **Date:** Test date (DD.MM.YYYY)
- **Users:** Number of concurrent virtual users
- **Spawn Rate:** Users added per second
- **Duration:** Test duration (minutes)
- **Total Requests:** Total number of requests
- **RPS:** Requests Per Second
- **Avg Response:** Average response time (milliseconds)
- **Min:** Minimum response time
- **Max:** Maximum response time
- **95%:** 95th percentile response time
- **99%:** 99th percentile response time
- **Success Rate:** Percentage of successful requests (200/404)
- **Failure Rate:** Percentage of failed requests
- **Cache Hit Rate:** Percentage of requests served from Redis cache
- **Notes:** Special notes (cache behavior, issues, etc.)

## System Resources (During Test)

System metrics captured from Prometheus/Grafana during test execution:

### Test #1 (31.12.2025 - ~50 users) - Cold Cache

| Service | CPU Usage | Memory Usage | Notes |
|---------|-----------|--------------|-------|
| ProductRead | 1m (~0.1%) | 12 Mi | Very low CPU - I/O bound waiting for ES |
| Redis Sentinel (3 pods) | 2-3m each | 2-4 Mi each | Low usage - cache miss dominant |
| Elasticsearch | 5m (~0.5%) | 510 Mi | Handling most requests - cache cold |
| ProductConsumer | - | - | Not monitored during test |

**Analysis:** Cache was cold (DBSIZE=0 after TTL expiration). Hit rate only 2.8% (135 hits / 4608 misses). High latency due to Elasticsearch queries.

### Test #2 (31.12.2025 - ~150 users) - Warm Cache

| Service | CPU Usage | Memory Usage | Notes |
|---------|-----------|--------------|-------|
| ProductRead | 93m (~9.3%) | 21 Mi | Moderate CPU - handling 266 RPS |
| Redis (primary) | 9m | 4-5 Mi | Serving cache hits efficiently |
| Elasticsearch | 71m (~7%) | 512 Mi | Reduced load compared to cold cache |
| ProductConsumer | - | - | Not monitored during test |

**Analysis:** Cache warming improved performance significantly. Hit rate ~3% (1089/35628). Despite low hit rate, response time improved 12x (569ms → 45ms). Redis + optimized ES queries working well.

## Cache Performance Analysis

### Cache Hit/Miss Ratio

| Test | Total Requests | Cache Hits | Cache Misses | Hit Rate | Avg Response (Cache Hit) | Avg Response (Cache Miss) |
|------|----------------|------------|--------------|----------|--------------------------|---------------------------|
| #1 (Cold) | 4617 | ~135 | ~4482 | 2.8% | ~3-10 ms | ~600+ ms |
| #2 (Warm) | 31963 | ~1089 | ~30874 | 3.0% | ~2-10 ms | ~50 ms |

### Expected Performance

- **Cache Hit Response Time:** 2-5ms (Redis in-memory)
- **Cache Miss Response Time:** 10-30ms (Elasticsearch query + Redis set)
- **Target Cache Hit Rate:** >80% after warmup phase

## Summary & Conclusions

### Key Findings

1. **System Capacity:** TBD concurrent users / TBD RPS with 1 replica
2. **Cache Effectiveness:** TBD% hit rate after warmup
3. **Response Time Comparison:** Cache hit TBD ms vs Cache miss TBD ms
4. **Bottleneck:** TBD (Redis, Elasticsearch, or Network)

### Performance Progression

| Test | Users | Replicas | Success Rate | Cache Hit Rate | Key Change |
|------|-------|----------|--------------|----------------|------------|
| #1 | 150 | 1 | -% | -% | Baseline |

### Recommendations

1. **For cache optimization:** Monitor Redis memory usage and eviction policy
2. **For better performance:** Increase Redis memory limit if needed
3. **For monitoring:** Track cache hit rate and alert if <70%
4. **For production:** Consider read replicas if cache hit rate is low

---

**Last Updated:** 31.12.2025
