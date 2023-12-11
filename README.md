# Lifetime and Daily Leaderboard for streaming application

## Overview

This repository contains a backend service built in Go for video streaming leaderboard. It implements a leaderboard functionality to track both daily and lifetime views of videos. The architecture follows the Onion pattern and draws inspiration from Go-Kit for its design.

## Features

- **Daily and Lifetime Views**: Tracks views of videos on a daily and overall (lifetime) basis.
- **Scalable & Distributed**: Utilizes Redis for caching and Consul for service discovery, making it scalable and easy to manage.
- **Onion Architecture**: Organized in layers for modularity and ease of maintenance.

## Tech Stack

- **Go (Golang)**: Main programming language used for development.
- **Redis**: Used for caching video views and leaderboard scores.
- **Docker**: Enables containerization for easy deployment and portability.
- **Consul**: Configuration loading for distributed architecture

## Architecture

The project is structured based on the Onion architecture, promoting clean and maintainable code. It's divided into layers:
- **Core**: Contains the business logic and domain models.
- **Repositories**: Handles data access and interactions with Redis.
- **Services**: Implements service-level logic.
- **Transport**: Manages HTTP endpoints.

## Setup

1. **Clone the Repository**:
```bash
git clone https://github.com/SurajKadam7/leaderboard_microservice.git
```

2. **Build & Run Using Docker**:
```bash
docker compose up -d
```

**Consul Configurations**
```json
{
      "key":"",
      "port":":",
      "address" : "",
      "poolSize":0,
      "username":"",
      "password": ""
}
```
