# Test and Coverage Report

## Backend

- Test command: `cd backend && go test ./... -cover`
- Result:
  - `internal/calculator`: 77.8% statement coverage
  - `internal/httpapi`: 58.3% statement coverage

## Frontend

- Test command: `cd frontend && npx vitest run --coverage`
- Result:
  - 2 tests passed
  - Coverage:
    - Statements: 61.66%
    - Branches: 61.53%
    - Functions: 66.66%
    - Lines: 63.63%

## Summary

- Backend tests passed and generated coverage output.
- Frontend tests passed and generated coverage output.
- The application includes unit tests and coverage reporting for both layers.
