# CQRS Sync DB Project

This project demonstrates a CQRS (Command Query Responsibility Segregation) architecture with synchronized databases. It includes multiple services for handling commands, queries, and caching using Redis.

---

## Prerequisites
- Docker and Docker Compose installed on your system.
- Ensure the following Docker networks are created:
  - `debezium-net`: Used for communication between services.
  - `redisnet`: Used for Redis caching services.

### Create Docker Networks
Run the following commands to create the required networks:
```bash
# Create debezium-net
docker network create --driver bridge debezium-net

# Create redisnet with a specific subnet
docker network create --driver bridge --subnet=172.30.0.0/24 redisnet
```
---

## Why Use Separate Networks?
- **`debezium-net`**: This network is used for communication between the main services (e.g., ProductRead, ProductWrite, Elasticsearch). It ensures isolation and avoids conflicts with other networks.
- **`redisnet`**: This network is dedicated to Redis services. It allows fine-grained control over Redis communication and ensures that Redis services are isolated from other parts of the system.

By using separate networks, we achieve better isolation, security, and control over service communication.

