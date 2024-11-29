# Music Library 🎶

This project implements an online music library with RESTful APIs to manage songs. It allows users to view, add, modify, and delete songs, with additional functionality for pagination and filtering.

## Features

- Get library data with filtering and pagination.
- Retrieve song lyrics with pagination by verses.
- Add a new song with the following JSON format:

```json
{
	"group": "Muse",
	"song": "Supermassive Black Hole"
}
```

- Delete songs from the library.
- Edit song details.

### Requirements

- Go 1.22 and above.

### Compilation

First clone the repository and compile the project:

```bash
git clone <repository_url>
cd ./music-library
go build -o music-library ./cmd
```

## Launch

To start the server, use the following command:

```bash
./music-library
```

### Database

The application stores enriched song details in a PostgreSQL database. The database schema is automatically created via migrations on service startup.

### Logging

Debug and info logs are included throughout the application.

### Configuration

All configuration settings (e.g., database credentials) are loaded from a .env file.

### Swagger API

Swagger documentation is generated for the implemented API.

## Project structure

```bash
music-library/
├── cmd/
│ ├── doc/
│ │   ├── docs.go
│ │   ├── swagger.json
│ │   └── swagger.yaml
│ └── main.go
├── internal/
│ ├── config/
│ │   └── config.go
│ ├── db/
│ │   └── db.go
│ ├── handlers/
│ │   ├── response.go
│ │   └── song_handler.go
│ ├── models/
│ │   └── song.go
│ ├── router/
│ │   └── router.go
│ └── services/
│     └── song_services.go
├── migration/
│   └── 001_create_song_table.sql
├── .env
├── go.mod
├── go.sum
└── README.md
```

### Author

Damir Usetov
