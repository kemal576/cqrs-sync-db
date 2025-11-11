.PHONY: deploy-elasticsearch deploy-redis deploy-product-read deploy-product-read-all port-forward-product-read delete-elasticsearch delete-redis delete-product-read delete-product-read-all

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

# ------------------------------------------------------------------------------------------------------------
delete-elasticsearch:
	kubectl delete -f Infrastructure/Kubernetes/elasticsearch-deployment.yaml -n elastic

delete-redis:
	kubectl delete -f Infrastructure/Kubernetes/Cache/redis-*.yaml -n cache
	kubectl delete -f Infrastructure/Kubernetes/Cache/redis-sentinel-*.yaml -n cache

delete-product-read:
	kubectl delete -f Infrastructure/Kubernetes/ProductRead/productread-deployment.yaml -n product-read

delete-product-read-all: delete-product-read delete-redis delete-elasticsearch