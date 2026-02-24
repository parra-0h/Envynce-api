# Config Service

A robust REST API built with Go to manage application configurations across different environments.

## Features
- CRUD for Applications, Environments, and Configurations.
- Automatic versioning of configurations.
- Audit logging of all changes.
- API Key authentication.
- SQLite database with GORM auto-migrations.
- Structured JSON responses.

## Setup

### Local Run
1. Ensure Go 1.21+ is installed.
2. Clone the repository.
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```
   The server will start on `http://localhost:8080`.

### Environment Variables
Configure `.env` file:
- `PORT`: Server port (default: 8080)
- `DB_PATH`: SQLite database path (default: config.db)
- `API_KEY`: Key for authentication (default: my-super-secret-key)

### Docker
```bash
docker build -t config-service .
docker run -p 8080:8080 config-service
```

## API Usage

### Authentication
Include the header `X-API-KEY` with the value defined in your `.env`.

### Endpoints

- **Health Check**: `GET /health`
- **Applications**:
  - `POST /applications`: Create a new application (`{"name": "App1"}`)
  - `GET /applications`: List all applications
- **Environments**:
  - `POST /environments`: Create a new environment (`{"name": "production"}`)
  - `GET /environments`: List all environments
- **Configurations**:
  - `POST /configs`: Create or update a config
    ```json
    {
      "key": "DATABASE_URL",
      "value": "postgres://localhost:5432",
      "application_id": 1,
      "environment_id": 1
    }
    ```
  - `GET /configs?application_id=1&environment_id=1`: List active configs
- **Audit Logs**:
  - `GET /audit-logs`: List all change history
