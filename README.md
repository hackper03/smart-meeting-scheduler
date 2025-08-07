# Smart Meeting Scheduler API

A Go-based RESTful API that intelligently schedules meetings by finding optimal time slots for multiple participants using a smart scoring algorithm. Built with Gin (HTTP web framework) and GORM (ORM library), and uses PostgreSQL (or SQLite optionally) as the database.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Setup & Installation](#setup--installation)
- [Environment Configuration](#environment-configuration)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Data Models](#data-models)
- [Scheduling Algorithm Details](#scheduling-algorithm-details)
- [Database Seeding](#database-seeding)
- [Logging](#logging)
- [Notes](#notes)

## Features

- Schedule meetings by finding optimal common available time slots across multiple users
- Smart scoring algorithm considering:
  - Preference for earlier time slots
  - Minimizing awkward calendar gaps
  - Buffer time between meetings
  - Respecting working hours (9 AM - 5 PM)
- User and meeting management with events and participant tracking
- REST API with well-structured JSON responses and error handling
- PostgreSQL support via GORM with environment-based configuration (can switch back to SQLite easily)
- Middleware with logging and graceful error recovery

## Architecture

The project has a clean modular structure:

```bash
smart-meeting-scheduler/
├── main.go # Application entry point
├── models/ # Database models
│ ├── user.go
│ ├── event.go
│ └── meeting.go
├── handlers/ # HTTP route handlers
│ ├── schedule.go
│ └── calendar.go
├── services/ # Business logic and scheduling algorithms
│ └── scheduler.go
├── database/ # Database connection and migrations
│ └── database.go
├── utils/ # Utility functions (e.g., JSON response helpers)
│ └── response.go
├── logger/ # Custom logger with configurable log levels and colors
│ └── log.go
├── go.mod # Go module dependencies
├── go.sum
└── README.md
```

## Setup & Installation

### Prerequisites

- Go 1.21 or later installed ([download](https://golang.org/dl/))
- PostgreSQL server
- `git` and command line tools

### Clone the repository

```bash
git clone <your-repository-url>
cd smart-meeting-scheduler
```

### Install dependencies

```bash
go mod tidy
```

## Environment Configuration

The application uses environment variables for configuration, loaded from a `.env` file. Example `.env` file:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=yourusername
DB_PASSWORD=yourpassword
DB_NAME=scheduler
SSL_MODE=disable
```

- To use **PostgreSQL**, set these values accordingly.
- If you want **SQLite**, you can adjust the database connection in `database/database.go` and skip the `.env` or set minimal config.

## Running the Application

```bash
go run main.go
```

The server starts at:

```bash
http://localhost:8080
```

Health endpoint for quick status check:

```bash
GET /health
Response: {"status":"healthy"}
```

## API Endpoints

### 1. Schedule a Meeting

**Request**

```bash
POST /api/v1/schedule
Content-Type: application/json
{
"participantIds": ["user1", "user2", "user3"],
"durationMinutes": 60,
"timeRange": {
"start": "2024-09-01T09:00:00Z",
"end": "2024-09-05T17:00:00Z"
}
```

- `participantIds`: Array of user IDs to include in the meeting
- `durationMinutes`: Length of the meeting in minutes (minimum 1)
- `timeRange`: Object specifying ISO8601/RFC3339 start and end times during which to schedule the meeting

**Response (Success - 201 Created)**

```bash
{
"meetingId": "meeting_xyz123",
"title": "New Meeting",
"participantIds": ["user1", "user2", "user3"],
"startTime": "2024-09-02T10:00:00Z",
"endTime": "2024-09-02T11:00:00Z"
}
```

**Possible Errors**

- `400 Bad Request`: Invalid JSON or missing required fields
- `404 Not Found`: No available slot found for all participants
- `500 Internal Server Error`: On server-side or database error

### 2. Get User's Calendar Events

**Request**

```bash
GET /api/v1/users/{userId}/calendar?start={startTime}&end={endTime}
```

- `userId`: Identifier of the user
- `start` & `end`: Date-time range in RFC3339 format

**Example**

```bash
GET /api/v1/users/user1/calendar?start=2024-09-01T00:00:00Z&end=2024-09-05T23:59:59Z
```

**Response (Success - 200 OK)**

```bash
[
{
"title": "Morning Standup",
"startTime": "2024-09-02T09:00:00Z",
"endTime": "2024-09-02T09:30:00Z"
},
{
"title": "Client Call",
"startTime": "2024-09-02T11:00:00Z",
"endTime": "2024-09-02T12:00:00Z"
}
]
```

## Data Models

- **User**

  - `ID` (string, primary key)
  - `Name` (string)
  - `Events` (list of events associated)

- **Event**

  - `ID` (string, primary key)
  - `UserID` (foreign key)
  - `Title` (event title)
  - `StartTime` (timestamp)
  - `EndTime` (timestamp)

- **Meeting**

  - `ID` (string, primary key)
  - `Title` (string)
  - `ParticipantIDs` (JSON array stored as string)
  - `StartTime` (timestamp)
  - `EndTime` (timestamp)

## Scheduling Algorithm Details

The scheduler:

- Collects all events for participants within the requested time window
- Finds all available gaps that can accommodate the desired meeting duration
- Assigns a **score** to each available slot via heuristics:
  - **Preference for earlier slots:** Higher score for slots closer to 9 AM
  - **Working hours bonus:** Slots during 9 AM to 5 PM favored
  - **Buffer Time:** Penalizes slots adjacent without enough buffer (e.g., 15 minutes)
  - **Minimize calendar gaps:** Avoids creating awkward 30-minute free gaps
  - **Back-to-back meeting bonus:** Rewards slots that tightly attach to existing meetings
- Chooses the highest scoring slot and creates meeting + events for participants

## Database Seeding

On startup, if the database contains no users, the app seeds initial test data:

- Users: Alice (user1), Bob (user2), Charlie (user3), Diana (user4)
- Sample calendar events for realistic testing of scheduling conflicts and availability

This helps you test and develop without manual data entry.

## Logging

- Custom logger with support for log levels:
  - FATAL, ERROR, WARN, INFO, DEBUG, TRACE
- Color-coded log outputs for easier reading in terminals
- Configurable log level via environment variables (`LOG_LEVEL`)
- Logs include timestamps and severity labels

## Notes

- **Database choice:** PostgreSQL is strongly recommended for production but the app can be adapted to SQLite for simple setups.
- **Time formats:** All API requests/responses use RFC3339 (ISO 8601) time format.
- **Error handling:** Uses consistent JSON error response structures with meaningful message and HTTP status codes.
- **Extensibility:** The architecture supports adding authentication, notifications, recurring meetings, or external calendar integrations.
- **Testing:** While not included here, the project supports unit and integration testing with mocks and table-driven tests using `testify`.
- **Deployment:** Can be containerized with Docker and connected to managed PostgreSQL instances for scalable deployment.

Thank you for using **Smart Meeting Scheduler**!  
If you encounter any issues or wish to contribute, please open an issue or submit a pull request.

## Future Scopes

Certainly! Here are the future scopes summarized into 3 concise bullet points in markdown format:

- Add User Authentication and Profile Management: Implement secure login, user roles, and preferences so users can manage their schedules, working hours, and availability more flexibly.
- Improve Scheduling Intelligence: Enhance the algorithm with features like support for recurring meetings, handling timezone differences, smart notifications, and conflict explanation to provide more user-friendly and efficient scheduling.
- Additional User Managements

**Happy Scheduling!**
