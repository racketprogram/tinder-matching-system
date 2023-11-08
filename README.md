# README and System Design for Matching System

## Overview

The Matching System simulates a simple matchmaking service, enabling users to find potential matches based on criteria like height and gender. It exposes a RESTful API for user interaction.

## Features

- Add a new user to the system.
- Find potential matches for a user.
- Remove a user from the system.
- Query the system for the most compatible individuals.

## API Endpoints

- `POST /add_single_person_and_match`: Add a user and retrieve potential matches.
- `POST /query_single_people`: Retrieve possible matches for a user.
- `POST /remove_single_person`: Remove a user from the system.

## Components

- **HTTP Server**: Serves as the entry point for incoming HTTP requests, routing them to appropriate handlers.
- **Handler Functions**: Process requests for endpoints, including user addition, match querying, and user removal.
- **In-Memory Data Store**: A thread-safe data structure that stores user data and supports match finding operations.
- **Matching Algorithm**: Implements the logic to determine potential matches based on user attributes.

## Data Structures

- **Users Map**: A `map[int]*Person` that maintains user data, indexed by user ID.

## Algorithms

### Matching

Iterates over potential matches, applying criteria to find suitable pairs. It terminates when the requested number of matches is found or the user list is fully traversed.

### Time Complexity

- **Add User**: O(1) - Insertion into the map is constant time.
- **Find Matches**: O(n) - Iteration over user map, where `n` is the total number of users.
- **Remove User**: O(1) - Deletion from the map using the user ID is constant time.

## Concurrency

A `sync.Mutex` is used to ensure thread-safe operations on the shared `Users` map during reads and writes.

## Scalability

The current implementation uses an in-memory store, which isn't persistent or horizontally scalable. For scaling out and persisting data, integrating with a database and adopting a stateless architecture is advisable.

## Deployment

The service is containerized with Docker, enabling straightforward deployment. The `Dockerfile` produces a secure, lightweight image.

## Running Tests

To run the tests for the Matching System, follow these steps:

```bash
# Navigate to the project directory
cd tinder-matching-system

# Run the tests
go test -v ./...
```

## Running the Application

Build and run the Docker container using:

```bash
# Navigate to the project directory
cd tinder-matching-system

docker build -t matching-system .
docker run -d -p 8080:8080 matching-system
```