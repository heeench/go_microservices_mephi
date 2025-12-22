package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"go-microservice/models"
	"go-microservice/services"
	"go-microservice/utils"
)

// UserHandler bundles dependencies for user endpoints.
type UserHandler struct {
	Service *services.UserService
	Logger  utils.AuditLogger
	Notify  utils.Notifier
}

func (h *UserHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/users", h.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/api/users/{id}", h.GetByID).Methods(http.MethodGet)
	r.HandleFunc("/api/users", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/api/users/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/api/users/{id}", h.Delete).Methods(http.MethodDelete)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users := h.Service.GetAll()
	writeJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.Service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created := h.Service.Create(user)
	go h.Logger.Log("CREATE", created.ID)
	go h.Notify.Send(created.ID, "created")
	writeJSON(w, http.StatusCreated, created)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload models.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.Service.Update(id, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	go h.Logger.Log("UPDATE", id)
	go h.Notify.Send(id, "updated")
	writeJSON(w, http.StatusOK, updated)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	go h.Logger.Log("DELETE", id)
	go h.Notify.Send(id, "deleted")
	w.WriteHeader(http.StatusNoContent)
}

func parseID(raw string) (int, error) {
	return strconv.Atoi(raw)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
