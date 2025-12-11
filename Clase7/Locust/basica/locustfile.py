from locust import HttpUser, task,between

class myUser(HttpUser):
    wait_time = between(1,3)

    @task
    def home_page(self):
        self.client.get("https://www.google.com")