# Task Manager API

A lightweight RESTful API for managing tasks, built with **Go**. This project demonstrates clean CRUD operations, filtering, sorting, pagination, and concurrency‑safe in‑memory storage – all using only the Go standard library.

## Features

- **CRUD + Patch** – Create, read, update, replace, and delete tasks.
- **Partial updates** – Send only the fields you want to change.
- **Filtering** – By title, text, and completion status.
- **Sorting** – By title (asc/desc), completion status, or creation time.
- **Pagination** – Use `offset` and `limit` query parameters.
- **Statistics** – Get total tasks, completed count, completion rate, and the latest task.
- **Concurrency‑safe** – Mutex‑protected in‑memory store.

## Tech Stack

- **Go 1.21+** – Standard library only
- `net/http` – HTTP server and routing
- `sync` – Concurrency control
- `encoding/json` – JSON serialisation

## Getting Started

### Prerequisites

- Go 1.21 or later installed.

### Installation

```bash
git clone https://github.com/Saik0-0/TaskManager.git
cd TaskManager
go run .
```
The server will start on http://localhost:9092

## API Endpoints

### Tasks Collection

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/tasks` | Create a new task |
| GET | `/tasks` | Get all tasks with optional filters, sorting, and pagination |

#### GET `/tasks` Query Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `title` | string | Filter tasks by title (partial match) |
| `text` | string | Filter tasks by text content (partial match) |
| `complete` | boolean | Filter by completion status (`true`/`false`) |
| `sort` | string | Sort tasks: `title`, `-title`, `completed`, `-completed`, `time`, `-time` |
| `offset` | integer | Pagination offset (default: 0) |
| `limit` | integer | Maximum items to return (default: all) |

### Single Task

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tasks/{id}` | Get a specific task by ID |
| PUT | `/tasks/{id}` | Replace an entire task |
| PATCH | `/tasks/{id}` | Partially update a task |
| DELETE | `/tasks/{id}` | Delete a task |

### Statistics

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/stats` | Get task statistics |


## Technical Highlights

- **In-Memory Storage**: Thread-safe task storage using mutex locks

- **Atomic ID Generation**: Uses sync/atomic for thread-safe ID generation

- **RESTful Design**: Follows REST conventions with proper HTTP methods and status codes

- **Comprehensive Filtering**: Supports multiple filter parameters with OR logic

- **Flexible Sorting**: Implements custom sorting for different field types

- **Pagination**: Efficient array slicing for paginated responses

- **Partial Updates**: Supports PATCH with pointer fields to distinguish between zero values and omitted fields