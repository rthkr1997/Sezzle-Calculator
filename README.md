# Full-Stack Calculator

A full-stack calculator application built with a React + TypeScript frontend and a Go backend microservice.

## What this project includes

- A browser calculator UI built with React and TypeScript
- A Go backend REST API that computes mathematical expressions
- Support for parentheses, powers, square root, and percentages
- Friendly error handling for invalid math and division by zero
- Unit tests for both frontend and backend
- Docker Compose for easy local deployment

## Project structure


```
calculator-app/
  backend/
    cmd/server/              # Go server entrypoint
    internal/calculator/     # Expression parser and calculator logic
    internal/httpapi/        # REST API handlers
    go.mod                   # Go module file
  frontend/
    src/api/                 # API service for backend calls
    src/components/          # Calculator UI component
    src/styles/              # CSS and design styles
    src/__tests__/           # Frontend unit tests
    package.json             # Node dependencies and scripts
  docker-compose.yml         # Full-stack Docker deployment
  README.md                  # Project documentation
```

## Setup for beginners

### 1. Install required software

You need these tools installed on your computer:

- Go 1.22 or newer
- Node.js 20 or newer
- Docker and Docker Compose (optional, only if you want containers)

### 2. Install frontend packages

From the project root, run:

```bash
cd frontend
npm install
```

This downloads the libraries needed for the React app.

### 3. Run the backend server

Open a new terminal window and run these commands from the project root:

```bash
cd backend
go run ./cmd/server
```

This starts the backend API on:

- `http://localhost:8080`

### 4. Run the frontend app

Open another terminal window and run these commands from the project root:

```bash
cd frontend
npm run dev
```

Then open the app in your browser at:

- `http://localhost:5173`

The frontend will send calculations to the backend automatically.

## Run everything with Docker (from project root)

If you want to run both frontend and backend together with Docker, use this command in the project root:

```bash
docker compose up --build
```

When Docker is ready:

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`

> Docker uses the frontend container on internal port `80`, then maps it to host port `3000`.

## Quick how-to use

1. Start the backend server from `backend/`.
2. Start the frontend app from `frontend/`.
3. Open the browser at `http://localhost:5173`.
4. Type an expression like `12 + 4 * (8 - 3) / 2`.
5. Press the calculate button and see the result.

## API examples

### Health check

Run this from any terminal:

```bash
curl http://localhost:8080/health
```

Response:

```json
{"status":"ok"}
```

### Calculate an expression

Run this from any terminal:

```bash
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression":"12 + 4 * (8 - 3) / 2"}'
```

Response:

```json
{"result":22,"expression":"12 + 4 * (8 - 3) / 2"}
```

### Operation endpoints

Run this from any terminal:

```bash
curl -X POST http://localhost:8080/api/operations/add \
  -H "Content-Type: application/json" \
  -d '{"values":[1,2,3,4]}'
```

Response:

```json
{"result":10,"operation":"add"}
```

Supported endpoints:

- `/api/operations/add`
- `/api/operations/subtract`
- `/api/operations/multiply`
- `/api/operations/divide`
- `/api/operations/power`
- `/api/operations/sqrt`
- `/api/operations/percentage`

## Error examples

### Division by zero

```bash
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression":"10 / 0"}'
```

Response:

```json
{"error":"division by zero"}
```

### Invalid input

```bash
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression":"1 +"}'
```

Response:

```json
{"error":"invalid expression: missing operand"}
```

## Running tests

### Backend tests (run from project root)

```bash
cd backend
go test ./...
```

### Frontend tests (run from project root)

```bash
cd frontend
npm test
```

## Design notes

- The backend parses and validates calculator expressions.
- The frontend sends expressions to the backend and shows results.
- The parser uses a standard algorithm for operator precedence.
- Percent (`%`) is treated as a postfix operator.
- `sqrt(...)` is supported.

## AI prompt used to build this project

This README and the application were built using the following detailed prompt:

> Build a full-stack calculator application with a React TypeScript frontend and a Go backend microservice. The app should support expression parsing with operator precedence, parentheses, exponentiation, square root, and percentage. Include both frontend and backend unit tests, API endpoints for expressions and named operations, error handling for invalid input and division by zero, responsive UI, and Docker Compose deployment. Write a clear README for beginners that explains where to run commands from the project root, how to start each service, how to use the API, and the exact prompt used to produce the application.
