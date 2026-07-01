package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"calculator-backend/internal/calculator"
)

type Server struct{}

type calculateRequest struct {
	Expression string `json:"expression"`
}
type operationRequest struct {
	Values []float64 `json:"values"`
}
type successResponse struct {
	Result     float64 `json:"result"`
	Expression string  `json:"expression,omitempty"`
	Operation  string  `json:"operation,omitempty"`
}
type errorResponse struct {
	Error string `json:"error"`
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	s := &Server{}
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/api/calculate", s.calculate)
	mux.HandleFunc("/api/operations/", s.operation)
	return cors(mux)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req calculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	result, err := calculator.Calculate(req.Expression)
	if err != nil {
		writeCalculatorError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, successResponse{Result: result, Expression: req.Expression})
}

func (s *Server) operation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	op := strings.TrimPrefix(r.URL.Path, "/api/operations/")
	var req operationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	result, err := calculator.Operation(op, req.Values)
	if err != nil {
		writeCalculatorError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, successResponse{Result: result, Operation: op})
}

func writeCalculatorError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, calculator.ErrDivisionByZero) || errors.Is(err, calculator.ErrNegativeSqrt) {
		status = http.StatusUnprocessableEntity
	}
	writeError(w, status, err.Error())
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
