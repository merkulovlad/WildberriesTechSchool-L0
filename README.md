# Demo Task to Wildberries Tech School

This project is a demonstration service with a simple frontend interface that displays **order data**.  
It is based on a **Go microservice**, using **Kafka**, **PostgreSQL**, and in-memory caching.

---

##  Description

The service implements the following flow:

1. **Order Data Ingestion**
    - Orders are produced into a Kafka topic.
    - A Go consumer service listens to Kafka, receives order data, and processes it.

2. **Data Storage**
    - Orders are added into a PostgreSQL database. 

3. **In-Memory Cache**
    - Orders are cached in memory for fast access.

4. **Frontend**
    - A minimal frontend
    - It fetches order data from the Go backend.
---

## 🛠️ Tech Stack

- **Go** — backend service
- **Kafka** — message queue
- **PostgreSQL** — database
- **Docker & Docker Compose** — containerization
- **Fiber** — web framework
- **Python** — for frontend server

---

## How to get started 

### 1. Clone the repository
```bash
git clone https://github.com/merkulovlad/WildberriesTechSchool-L0.git
cd WildberriesTechSchool-L0
```
### 2. Run services with Docker
```bash
docker-compose up --build
```
### 3. Run frontend(with Python)
```bash
cd frontend
python3 -m http.server 3001
```