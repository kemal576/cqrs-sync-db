.PHONY: deploy-elasticsearch deploy-redis deploy-product-read deploy-product-read-all port-forward-product-read delete-elasticsearch delete-redis delete-product-read delete-product-read-all deploy-mongo deploy-lock-redis deploy-product-write deploy-product-write-all port-forward-product-write delete-mongo delete-mongo-full delete-lock-redis delete-product-write delete-product-write-all

# ============================================================================================================
# ProductRead Deployment Commands
# ============================================================================================================
deploy-elasticsearch:
	kubectl apply -f Infrastructure/Kubernetes/elasticsearch-deployment.yaml -n elastic
	@echo "Waiting for Elasticsearch to be ready..."
	@kubectl wait --for=condition=ready pod -l app=elasticsearch -n elastic --timeout=60s

deploy-redis:
	kubectl apply -f Infrastructure/Kubernetes/Cache/redis-*.yaml -n cache
	kubectl apply -f Infrastructure/Kubernetes/Cache/redis-sentinel-*.yaml -n cache
	@echo "Waiting for Redis to be ready..."
	@kubectl wait --for=condition=ready pod -l app=redis -n cache --timeout=60s
	@kubectl wait --for=condition=ready pod -l app=redis-sentinel -n cache --timeout=60s

deploy-product-read:
	kubectl apply -f Infrastructure/Kubernetes/ProductRead/productread-deployment.yaml -n product-read
	@echo "Waiting for ProductRead to be ready..."
	@kubectl wait --for=condition=ready pod -l app=product-read -n product-read --timeout=60s

port-forward-product-read:
	@echo "Port forwarding ProductRead service to localhost:5001..."
	kubectl port-forward -n product-read svc/productread 5001:8080

deploy-product-read-all: deploy-elasticsearch deploy-redis deploy-product-read port-forward-product-read

# ============================================================================================================
# ProductWrite Deployment Commands
# ============================================================================================================
deploy-mongo:
	kubectl apply -f Infrastructure/Kubernetes/Mongo/mongodb-deployment.yaml -n mongo
	@echo "Waiting for MongoDB to be ready..."
	@kubectl wait --for=condition=ready pod -l app=mongodb -n mongo --timeout=120s
	@echo "Initializing MongoDB replica set..."
	kubectl apply -f Infrastructure/Kubernetes/Mongo/mongodb-init-job.yaml -n mongo
	@echo "Waiting for MongoDB initialization to complete..."
	@kubectl wait --for=condition=complete job/mongodb-init-replicaset -n mongo --timeout=60s
	@echo "MongoDB replica set initialized successfully!"

deploy-lock-redis:
	kubectl apply -f Infrastructure/Kubernetes/Lock/lock-redis-deployment.yaml -n lock
	@echo "Waiting for Lock Redis to be ready..."
	@kubectl wait --for=condition=ready pod -l app=redis-lock -n lock --timeout=60s

deploy-product-write:
	kubectl apply -f Infrastructure/Kubernetes/ProductWrite/productwrite-deployment.yaml -n product-write
	@echo "Waiting for ProductWrite to be ready..."
	@kubectl wait --for=condition=ready pod -l app=product-write -n product-write --timeout=60s

port-forward-product-write:
	@echo "Port forwarding ProductWrite service to localhost:5002..."
	kubectl port-forward -n product-write svc/productwrite 5002:8080

deploy-product-write-all: deploy-mongo deploy-lock-redis deploy-product-write port-forward-product-write

# ============================================================================================================
# Delete Commands
# ============================================================================================================
delete-elasticsearch:
	kubectl delete -f Infrastructure/Kubernetes/elasticsearch-deployment.yaml -n elastic

delete-redis:
	kubectl delete -f Infrastructure/Kubernetes/Cache/redis-*.yaml -n cache
	kubectl delete -f Infrastructure/Kubernetes/Cache/redis-sentinel-*.yaml -n cache

delete-product-read:
	kubectl delete -f Infrastructure/Kubernetes/ProductRead/productread-deployment.yaml -n product-read

delete-product-read-all: delete-product-read delete-redis delete-elasticsearch

delete-mongo:
	kubectl delete -f Infrastructure/Kubernetes/Mongo/mongodb-deployment.yaml -n mongo

delete-lock-redis:
	kubectl delete -f Infrastructure/Kubernetes/Lock/lock-redis-deployment.yaml -n lock

delete-product-write:
	kubectl delete -f Infrastructure/Kubernetes/ProductWrite/productwrite-deployment.yaml -n product-write

delete-product-write-all: delete-product-write delete-lock-redis delete-mongo