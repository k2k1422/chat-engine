# Chat engine

## Introduction
The Chat Engine project is designed to facilitate multiple WebSocket connections for users, whether they are the same user or different users. Built with scalability in mind, it can be deployed across multiple instances to handle varying loads.


## Outline of Capabilities
1. WebSocket Connection Management: Maintains multiple WebSocket connections for users, allowing for real-time communication.
2. Scalability: Designed to scale horizontally, enabling deployment across multiple instances to handle increased traffic and users.
3. JWT Authentication: Implements basic JWT (JSON Web Token) authentication for user verification. Additionally, includes a standalone Python script for token generation, ensuring secure communication with a shared secret key.
4. User Management API: Provides an API endpoint for creating users, offering a straightforward way to manage user accounts.
5. Microservices Architecture: Utilizes Kafka for inter-service communication within the microservices architecture. Each service creates and listens to topics, enabling efficient communication while maintaining a list of user WebSocket connections.


## Overview
The Chat Engine project employs a microservices architecture to ensure modularity, scalability, and resilience. Here's an overview:

WebSocket Connection Manager: Manages WebSocket connections for users, facilitating real-time communication.
1. Authentication Service: Handles user authentication using JWT tokens generated with a shared secret key.
2. User Management Service: Provides functionality to create and manage user accounts via a RESTful API.
3. Messaging Service: Facilitates communication between users through WebSocket connections, utilizing Kafka for inter-service communication.
4. Kafka Message Broker: Serves as a centralized message broker, enabling seamless communication between microservices. Each service creates and listens to topics to efficiently transmit messages.

## How to run (Standlone)

### Set up enviroment

1. Need Go lang installed
2. Need python installed
3. Install docker and docker compose

Run this to export the env to local
```
source ./update_env.sh
```

Make the Infra Up for the service
```
docker-compose u -d
```

Run the go application
```
go run main.go
```

Create a user with username
```
curl --location --request GET 'http://localhost:8000/api/user/v1/create' \
--header 'Content-Type: application/json' \
--data '{
    "username": "kiran"
}'
```

Create another user with username
```
curl --location --request GET 'http://localhost:8000/api/user/v1/create' \
--header 'Content-Type: application/json' \
--data '{
    "username": "kumar"
}'
```

Create a jwt token from python
```
python3 jwt.py kiran
```

output
```
JWT Token: eyJhbGciOiAiSFMyNTYiLCAidHlwIjogIkpXVCJ9.eyJ1c2VybmFtZSI6ICJzYW0iLCAiZXhwIjogMTcxMzc1OTk2NX0.shOgOIpgFFqcxrwDLrzho_zLPtqffY1lQS7Up_mP-4w
```

### In post man create a webscoket connection 

#### Add token to the request
![Alt text](image/image3.png)

#### Create a api call
![Alt text](image/image1.png)

The architecture ensures that the Chat Engine is robust, scalable, and capable of handling large volumes of concurrent users while maintaining security and performance.






