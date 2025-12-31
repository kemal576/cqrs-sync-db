# Create Product Load Test Results

This document contains load test results for the ProductWrite Create Product endpoint in the CQRS CDC project.

## Test Environment

- **Cluster:** Minikube (Docker Driver)
- **Tool:** Locust
- **Test Script:** CreateProduct.py
- **Endpoint:** POST /products (Create Product)
- **Initial Scale:** All services run with **1 replica**, later scaled to **2 replicas**

### Pod Resource Specifications

| Service | CPU Request | CPU Limit | Memory Request | Memory Limit | Replicas (Test #1-2) | Replicas (Test #3-8) |
|--------|-------------|-----------|----------------|--------------|---------------------|---------------------|
| ProductWrite | 200m | 500m | 256Mi | 512Mi | 1 | 2 |
| ProductRead | 200m | 500m | 256Mi | 512Mi | 1 | 1 |
| ProductConsumer | 100m | 500m | 128Mi | 256Mi | 1 | 1 |
| MongoDB | - | - | - | - | 1 | 1 |
| Kafka | - | - | - | - | 1 | 1 |
| Elasticsearch | - | - | - | - | 1 | 1 |

**Note:** MongoDB, Kafka, and Elasticsearch have no resource limits defined.
**Connection Pool Config:** MaxConnectionPoolSize=200 (Test #1-7), 300 (Test #8), MinConnectionPoolSize=50

## ProductWrite - Create Product Endpoint

### Test Results

| Date | Users | Spawn Rate | Duration | Total Requests | RPS | Avg Response (ms) | Min (ms) | Max (ms) | 95% (ms) | 99% (ms) | Success Rate | Failure Rate | Notes |
|-------|-------|------------|----------|----------------|-----|-------------------|----------|----------|----------|----------|--------------|--------------|-------|
| 31.12.2025 | 150 | 5 | 2 min | 7916 | 74.9 | 13.29 | 3 | 87 | 40 | 50 | 100% | 0% | First test - Median: 8ms |
| 31.12.2025 | 500 | 10 | 2 min | 18141 | 185.5 | 10.67 | 4 | 79 | 28 | 44 | 74.2% | 25.8% | Connection timeouts (status 0) - System overwhelmed |
| 31.12.2025 | 500 | 10 | 2 min | 18217 | 184.7 | 11.49 | 4 | 1218 | 29 | 46 | 74.3% | 25.7% | 2 replicas + pool config - Load balancing issue |
| 31.12.2025 | 500 | 10 | 1 min | 9106 | 129.1 | 13.77 | 4 | 253 | 35 | 61 | 81.3% | 18.7% | Ingress load balancing - Improved but still failing |
| 31.12.2025 | 300 | 5 | 2 min | 11864 | 98.9 | 16.96 | 5 | 1156 | 42 | 73 | 100% | 0% | SUCCESS! 2 replicas + Ingress - Optimal config |
| 31.12.2025 | 380 | 5 | 2 min | 11950 | 99.6 | 13.82 | 5 | 103 | 35 | 49 | 99.79% | 0.21% | Near perfect! System capacity identified at ~380 users |
| 31.12.2025 | 400 | 5 | 2 min | 14318 | 119.3 | 13.41 | 4 | 101 | 34 | 51 | 93.55% | 6.45% | Over capacity - Connection pool exhaustion begins |
| 31.12.2025 | 500 | 10 | 1.7 min | 15005 | 146.1 | 13.03 | 4 | 159 | 34 | 53 | 75.25% | 24.75% | 300 pool size - Still over capacity |

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
- **Success Rate:** Percentage of successful requests
- **Failure Rate:** Percentage of failed requests
- **Notes:** Special notes (system status, issues, etc.)

## System Resources (During Test)

System metrics captured from Grafana during test execution:

### Test #1 (31.12.2025 - 150 users)

| Service | CPU Usage | Memory Usage | Notes |
|---------|-----------|--------------|-------|
| ProductWrite | Max 6.56% | Max 108 Mi | Well under resource limits (500m CPU, 512Mi RAM) |
| MongoDB | Max 19.35% | Max 379 Mi | No resource limits defined |
| Kafka | Max 1.06% | Max 421 Mi | No resource limits defined |
| Elasticsearch | Max 0.47% | Max 524 Mi | No resource limits defined |
| ProductConsumer | Max 0.11% | Max 12 Mi | Well under resource limits (500m CPU, 256Mi RAM) |

### Test #2 (31.12.2025 - 500 users)

| Service | CPU Usage | Memory Usage | Notes |
|---------|-----------|--------------|-------|
| ProductWrite | Max 11.54% | Max 106 Mi | Still under limits - CPU usage doubled |
| MongoDB | Max 26.16% | Max 381 Mi | CPU increased 35% from test #1 |
| Kafka | Max 0.99% | Max 431 Mi | Minimal CPU impact |
| Elasticsearch | Max 0.75% | Max 512 Mi | CPU slightly increased |
| ProductConsumer | Max 0.11% | Max 12 Mi | No change from test #1 |

### Test #3 (31.12.2025 - 500 users, 2 replicas)

| Service | CPU Usage | Memory Usage | Notes |
|---------|-----------|--------------|-------|
| ProductWrite (2 pods) | Max 10.48% / 8.95% | Max 70 Mi / 66 Mi | Load not distributed - NodePort issue |
| MongoDB | Max 25.83% | Max 375 Mi | Similar to test #2 |
| Kafka | Max 1.11% | Max 455 Mi | Memory increased slightly |
| Elasticsearch | Max 0.69% | Max 504 Mi | Stable |
| ProductConsumer | Max 0.1% | Max 10 Mi | Stable |

**Issue Identified:** NodePort service doesn't load balance properly. Traffic hits single pod causing same failure rate.

### Test #4 (31.12.2025 - 500 users, 2 replicas + Ingress)

| Service | CPU Usage | Memory Usage | Notes |
|---------|-----------|--------------|-------|
| ProductWrite (2 pods) | Max 6.49% / 6.79% | Max 81 Mi / 76 Mi | Load balanced via Ingress - Equal distribution |
| MongoDB | Max 30.59% | Max 391 Mi | CPU at 30% - potential bottleneck |
| Kafka | Max 0.48% | Max 211 Mi | Memory dropped significantly |
| Elasticsearch | Max 0.52% | Max 499 Mi | Stable |
| ProductConsumer | Max 0.12% | Max 8 Mi | Stable |

**Improvements:** Ingress successfully load balances traffic. Failure rate improved from 25.8% to 18.7%. MongoDB CPU increasing suggests it may be the bottleneck.

### Test #5 (31.12.2025 - 300 users, 2 replicas + Ingress) - ✅ SUCCESS

| Service | CPU Usage | Memory Usage | Notes |
|---------|-----------|--------------|-------|
| ProductWrite (2 pods) | Max 6.67% / 7.35% | Max 79 Mi / 67 Mi | Perfect load balancing - 0% failure rate! |
| MongoDB | Max 25.34% | Max 529 Mi | Optimal CPU usage |
| Kafka | Max 29.93% | Max 223 Mi | Higher CPU but stable |
| Elasticsearch | Max 50% | Max 579 Mi | High CPU during indexing |
| ProductConsumer | Max 0.17% | Max 11 Mi | Minimal resource usage |

**Result:** 100% success rate achieved! Reducing users from 500 to 300 kept total connections (2 × 200 = 400) below MongoDB's capacity. System can handle ~100 RPS with 2 replicas.

### Test #6 (31.12.2025 - 380 users, 2 replicas + Ingress) - ✅ NEAR PERFECT

| Service | CPU Usage | Memory Usage | Notes |
|---------|-----------|--------------|-------|
| ProductWrite (2 pods) | Max 6.05% / 6.01% | Max 80 Mi / 70 Mi | Excellent load distribution |
| MongoDB | Max 25.07% | Max 377 Mi | Optimal performance |
| Kafka | Max 0.72% | Max 200 Mi | Minimal CPU usage |
| Elasticsearch | Max 0.68% | Max 510 Mi | Low CPU, stable memory |
| ProductConsumer | Max 0.12% | Max 10 Mi | Very efficient |

**Result:** 99.79% success rate (25 failures out of 11,950). System capacity identified at ~380 concurrent users / ~100 RPS. This is the sweet spot before connection exhaustion.

### Test #7 (31.12.2025 - 400 users, 2 replicas + Ingress) - ⚠️ OVER CAPACITY

| Service | CPU Usage | Memory Usage | Notes |
|---------|-----------|--------------|-------|
| ProductWrite (2 pods) | Max 6.35% / 6.51% | Max 81 Mi / 72 Mi | Still balanced but connections exhausted |
| MongoDB | Max 26.38% | Max 382 Mi | Connection limit reached |
| Kafka | Max 0.44% | Max 196 Mi | Very low CPU usage |
| Elasticsearch | Max 0.68% | Max 507 Mi | Stable performance |
| ProductConsumer | Max 0.12% | Max 10 Mi | Minimal resource usage |

**Result:** 93.55% success rate (923 failures). System exceeds capacity at 400 users. Connection pool (2 × 200 = 400) is at its limit. **Maximum capacity confirmed: 380-390 concurrent users / ~120 RPS.**

### Test #8 (31.12.2025 - 500 users, 2 replicas, 300 pool + Ingress) - ⚠️ STILL FAILING

| Service | CPU Usage | Memory Usage | Notes |
|---------|-----------|--------------|-------|
| ProductWrite (2 pods) | Max 8.48% / 9.97% | Max 63 Mi / 78 Mi | Higher CPU but still failing |
| MongoDB | Max 27.92% | Max 385 Mi | CPU stable at ~28% |
| Kafka | Max 1.06% | Max 197 Mi | Low CPU, stable |
| Elasticsearch | Max 0.63% | Max 504 Mi | Minimal CPU usage |
| ProductConsumer | Max 0.11% | Max 10 Mi | Very efficient |

**Result:** 75.25% success rate with 300 pool size. Increasing pool from 200 to 300 improved from 74.2% to 75.25%, but 500 users still exceeds system capacity. Bottleneck is not just connection pool but overall throughput. **System optimal: 380-400 users max with 2 replicas.**

## Summary & Conclusions

### Key Findings

1. **System Capacity:** 380-390 concurrent users / ~120 RPS with 2 replicas
2. **Optimal Configuration:** 2 replicas + Nginx Ingress + 200 connection pool per pod
3. **Bottleneck:** Connection pool exhaustion at ~400 concurrent connections
4. **Resource Utilization:** ProductWrite CPU <10%, MongoDB CPU ~28% (not fully utilized)

### Performance Progression

| Test | Users | Replicas | Success Rate | Key Change |
|------|-------|----------|--------------|------------|
| #1 | 150 | 1 | 100% | Baseline |
| #2 | 500 | 1 | 74.2% | System overwhelmed |
| #3 | 500 | 2 | 74.3% | NodePort - no improvement |
| #4 | 500 | 2 + Ingress | 81.3% | Ingress improved 7% |
| #5 | 300 | 2 + Ingress | 100% | Reduced load - optimal |
| #6 | 380 | 2 + Ingress | 99.79% | System limit identified |
| #7 | 400 | 2 + Ingress | 93.55% | Over capacity |
| #8 | 500 | 2 + Ingress + 300 pool | 75.25% | Pool increase - minimal effect |

### Recommendations

1. **For 500+ users:** Add 3rd replica (600-900 total connections)
2. **For better latency:** Increase `WaitQueueTimeout` from 5s to 10-15s
3. **For monitoring:** Set alerts at 90% capacity (~350 concurrent users)
4. **For production:** Use 2 replicas with HPA (Horizontal Pod Autoscaler) targeting 350 users

---

**Last Updated:** 31.12.2025
