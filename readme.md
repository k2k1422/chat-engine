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


The architecture ensures that the Chat Engine is robust, scalable, and capable of handling large volumes of concurrent users while maintaining security and performance.






