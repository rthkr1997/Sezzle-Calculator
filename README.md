# Full-Stack Calculator

A maintainable full-stack calculator application with a React + TypeScript frontend and a Go REST API backend. It supports full expressions, not only one- or two-operand requests.

## Features

- Basic operations: addition, subtraction, multiplication, division
- Advanced operations: exponentiation, square root, percentage
- Multi-step expressions with parentheses, for example: `12 + 4 * (8 - 3) / 2`
- REST API validation and JSON responses
- Frontend input validation, API error handling, responsive layout, and recent calculation history
- Unit tests for backend calculator logic, API handlers, and frontend UI behavior
- Dockerfiles and Docker Compose for local full-stack deployment

## Repository Layout

```text
calculator-app/
  backend/
    cmd/server/              # Go HTTP server entrypoint
    internal/calculator/     # Expression parser/evaluator and operation service
    internal/httpapi/        # REST API handlers
  frontend/
    src/api/                 # API client
    src/components/          # React calculator UI
    src/styles/              # CSS
    src/__tests__/           # Frontend tests
  docker-compose.yml
```

## Backend Setup

Requirements: Go 1.22+

```bash
cd backend
go test ./... -cover
PORT=8080 go run ./cmd/server
```

The backend runs on `http://localhost:8080` by default.

## Frontend Setup

Requirements: Node.js 20+ or 22+

```bash
cd frontend
npm install
npm run dev
```

The frontend runs on `http://localhost:5173` and calls `http://localhost:8080` by default.

To point the frontend at another API URL:

```bash
VITE_API_URL=http://localhost:8080 npm run dev
```

## Run Tests and Coverage

Backend:

```bash
cd backend
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Frontend:

```bash
cd frontend
npm install
npm test
```

Vitest writes coverage output under `frontend/coverage/`.

## Docker Compose

```bash
docker compose up --build
```

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`

## API Usage

### Health Check

```bash
curl http://localhost:8080/health
```

Response:

```json
{"status":"ok"}
```

### Calculate an Expression

```bash
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression":"12 + 4 * (8 - 3) / 2"}'
```

Response:

```json
{"result":22,"expression":"12 + 4 * (8 - 3) / 2"}
```

More examples:

```bash
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression":"sqrt(81) + 50%"}'
```

```json
{"result":9.5,"expression":"sqrt(81) + 50%"}
```

### Operation Endpoints

These endpoints are included for direct operation-style API calls. They accept multiple operands where the operation logically supports it.

```bash
curl -X POST http://localhost:8080/api/operations/add \
  -H "Content-Type: application/json" \
  -d '{"values":[1,2,3,4]}'
```

```json
{"result":10,"operation":"add"}
```

Supported operation paths:

- `/api/operations/add`
- `/api/operations/subtract`
- `/api/operations/multiply`
- `/api/operations/divide`
- `/api/operations/power`
- `/api/operations/sqrt`
- `/api/operations/percentage`

## Error Examples

Division by zero returns HTTP `422`:

```bash
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression":"10 / 0"}'
```

```json
{"error":"division by zero"}
```

Invalid syntax returns HTTP `400`:

```json
{"error":"invalid expression: missing operand"}
```

## Design Decisions and Assumptions

- The primary API is expression-based because a full calculator should support chained expressions, precedence, unary minus, and parentheses rather than forcing only two operands.
- The backend owns calculation correctness. The frontend performs lightweight character validation, but backend validation is authoritative.
- The parser uses tokenization + shunting-yard conversion to Reverse Polish Notation, then evaluates the RPN stack. This keeps precedence handling explicit and testable.
- Percent is treated as a postfix unary operator: `50%` becomes `0.5`.
- `sqrt` is supported as a function: `sqrt(81)`.
- Division by zero and square root of negative numbers return `422 Unprocessable Entity`; malformed input returns `400 Bad Request`.
- The frontend is intentionally framework-light: React state, a typed API client, and component-level tests.

## AI Tooling / Prompts Used

The following prompt was used to guide implementation:

> Build a full-stack calculator application with a React TypeScript frontend and a Go backend microservice. Include basic and advanced arithmetic operations, support expressions with more than two operands, validate inputs, handle errors like division by zero, include unit tests, Docker support, and a README with setup instructions, API examples, design decisions, and prompts used.

## Suggested Git Commands

```bash
git init
git add .
git commit -m "Build full-stack calculator application"
git branch -M main
git remote add origin <your-repository-url>
git push -u origin main
```
