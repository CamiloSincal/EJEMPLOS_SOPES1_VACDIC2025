# Locustfile.py
from locust import HttpUser, TaskSet, task, between
import random
import json

class MyTasks(TaskSet):
    
    @task(1)
    def engineering(self):
        # Listado aleatorio de nombres
        names = ["guatemala", "mexico", "panama", "inglaterra", "francia", "italia", "españa", "argentina", "chile", "colombia"]

        climas = ["soleado", "nublado", "lluvioso"]
    
        # Datos de climas
        weather_data = {
            "name": random.choice(names),  # Random name
            "temperatura": random.randint(18, 28),  # Temperatura random entre 18 y 28
            "humedad": random.randint(40, 80),  # Humedad random entre 40 y 80
            "clima": random.choice(climas)  # Clima aleatorio
        }
        
        # Envio de JSON hacia route como POST
        headers = {'Content-Type': 'application/json'}
        self.client.post("/clima", json=weather_data, headers=headers)

class WebsiteUser(HttpUser):
    tasks = [MyTasks]
    wait_time = between(1, 5)  # Tiempo de espera entre tareas entre 1 y 5 segundos