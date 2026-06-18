package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"medicare-gateway/internal/gateway"
)

type Router struct {
	service *gateway.Service
	mux     *http.ServeMux
}

func New(service *gateway.Service) http.Handler {
	r := &Router{service: service, mux: http.NewServeMux()}
	r.routes()
	return r.withMiddleware(r.mux)
}

func (r *Router) routes() {
	r.mux.HandleFunc("GET /healthz", r.health)
	r.mux.HandleFunc("POST /api/sign-in", r.signIn)
	r.mux.HandleFunc("POST /api/person", r.person)
	r.mux.HandleFunc("POST /api/{infno}", r.call)
}

func (r *Router) withMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, req)
	})
}

func (r *Router) health(w http.ResponseWriter, req *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (r *Router) signIn(w http.ResponseWriter, req *http.Request) {
	resp, err := r.service.SignIn(req.Context())
	writeGatewayResponse(w, resp, err)
}

func (r *Router) person(w http.ResponseWriter, req *http.Request) {
	resp, err := r.service.PersonInfo(req.Context())
	writeGatewayResponse(w, resp, err)
}

func (r *Router) call(w http.ResponseWriter, req *http.Request) {
	infno := strings.TrimSpace(req.PathValue("infno"))
	if !allowedInfno(infno) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "unsupported medicare infno"})
		return
	}
	var erpReq gateway.ERPRequest
	if req.Body != nil {
		defer req.Body.Close()
		if err := json.NewDecoder(req.Body).Decode(&erpReq); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json: " + err.Error()})
			return
		}
	}
	resp, err := r.service.Call(req.Context(), infno, erpReq)
	writeGatewayResponse(w, resp, err)
}

func writeGatewayResponse(w http.ResponseWriter, resp gateway.Response, err error) {
	if err != nil {
		status := http.StatusBadGateway
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			status = http.StatusRequestTimeout
		}
		writeJSON(w, status, resp)
		return
	}
	if resp.Status == "queued" {
		writeJSON(w, http.StatusAccepted, resp)
		return
	}
	if resp.Status == "failed" {
		writeJSON(w, http.StatusBadGateway, resp)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func allowedInfno(infno string) bool {
	switch infno {
	case "9001", "1101", "2101", "2102", "2103", "3505", "3201", "3202":
		return true
	default:
		return false
	}
}
