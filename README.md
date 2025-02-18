# KianQuiz

## Overview

KianQuiz is a quiz game designed for 1v1 player interactions. Players compete through rounds, and a winner is determined based on scores. The system leverages modern technologies such as Golang for backend services, Redis for message brokering, and MongoDB for database storage.

## Key Features

- **User Management:** Registration, login, JWT-based authentication, and profile management.
- **Matchmaking:** Matches players 1v1 based on predefined criteria (e.g., categories like 'Football').
- **Real-Time Updates:** Utilizes Redis for presence and match notifications.
- **Leaderboards:** Tracks and displays player rankings.

---

## High-Level Architecture

### 1. **Components**

#### Backend Services

- **Authentication Service**: Manages user login, registration, and JWT token generation.
- **Matchmaking Service**: Matches players based on criteria and publishes match events to Redis.
- **Game Service**: Handles gameplay logic, including rounds, questions, and scoring.
- **Leaderboard Service**: Computes and updates player rankings.

#### Database

- **MongoDB**: Stores user data, game data, and leaderboard information.

#### Message Broker

- **Redis**: Handles real-time communication between services (e.g., matchmaking events, game state updates).

#### Frontend

- **Web Application**: Player-facing interface for login, matchmaking, gameplay, and leaderboards.

---

### 2. **Architecture Diagram**

```plaintext
+-------------------+       +---------------------+
|                   |       |                     |
|   Web Frontend    | <-->  |    API Gateway      |
|                   |       |                     |
+-------------------+       +---------------------+
                                   |
                                   v
                        +-------------------+
                        | Backend Services  |
                        |                   |
                        |  Auth Service     |
                        |  Matchmaking      |
                        |  Game Logic       |
                        |  Leaderboards     |
                        +-------------------+
                                   |
                                   v
        +-------------------+       +-------------------+
        |                   |       |                   |
        |     MongoDB       | <-->  |      Redis        |
        |                   |       |                   |
        +-------------------+       +-------------------+
```

---

## Detailed Design

### 1. **Authentication Service**

- **Technologies**: Golang, JWT
- **Responsibilities**:
  - User registration and login
  - Token generation and validation

#### Endpoints:

- `POST /auth/register`
- `POST /auth/login`
- `GET /auth/profile`

---

### 2. **Matchmaking Service**

- **Technologies**: Golang, Redis
- **Responsibilities**:
  - Match players based on categories
  - Publish match events to Redis

#### Flow:

1. Receive player matchmaking request.
2. Query Redis for available players waiting lisat in the requested category.
3. Pair players and publish a match event.
4. Create game and publish game is ready event.
5. Fetch questions one by one using game id and auth key
6. After the game, scores are calculated and leaderboard is updated

## Future Enhancements

- **Chat Service**: Allow players to chat during matches.
- **Question Pool Optimization**: Implement object pooling or dynamic generation for scalable question handling.
