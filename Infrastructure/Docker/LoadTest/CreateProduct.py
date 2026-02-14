from locust import HttpUser, task, between

class ProductWriteUser(HttpUser):
    wait_time = between(1, 3)
    
    @task
    def create_product(self):
        """Create a new product"""
        product_data = {
            "name": "LoadTest Product",
            "description": "Auto-generated product for load testing",
            "price": 99.99
        }
        
        with self.client.post(
            "/products",
            json=product_data,
            catch_response=True,
            name="Create Product"
        ) as response:
            if response.status_code == 201 or response.status_code == 200:
                response.success()
            else:
                response.failure(f"Got status code {response.status_code}")
    
    def on_start(self):
        """Called when a simulated user starts"""
        print("Starting load test user...")

    # docker run -p 8089:8089 -v "${PWD}:/mnt/locust" locustio/locust -f /mnt/locust/CreateProduct.py --host=http://host.docker.internal:PORT

