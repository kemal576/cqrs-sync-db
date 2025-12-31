from locust import HttpUser, task, between, events
import random
import string

class ProductReadUser(HttpUser):
    wait_time = between(0.1, 0.5)
    product_ids = []
    
    def on_start(self):
        """Called when a user starts - warmup phase to populate product IDs"""
        print("Starting load test user...")
        self.product_ids = ["PRODUCT_IDS"]
    
    @task
    def get_product_by_id(self):
        """Get a product by ID - tests cache and database performance"""
        if not self.product_ids:
            return
        
        # Randomly select a product ID
        product_id = random.choice(self.product_ids)
        
        with self.client.get(
            f"/products/{product_id}",
            catch_response=True,
            name="/products/:id"
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 404:
                # 404 is expected for non-existent products
                response.success()
            else:
                response.failure(f"Got status code {response.status_code}")

@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    print("Load test starting...")
    print("Note: Ensure ProductRead service is accessible and products exist in Elasticsearch")
    print("Cache will be populated during the test")

@events.test_stop.add_listener
def on_test_stop(environment, **kwargs):
    print("Load test finished!")
