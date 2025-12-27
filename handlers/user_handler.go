package handlers

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

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

// MiddlewareWrapper обертка для middleware
type MiddlewareWrapper func(handler httprouter.Handle) httprouter.Handle

func (h *UserHandler) RegisterRoutes(r *httprouter.Router, mw MiddlewareWrapper) {
	// Регистрируем с middleware
	r.GET("/api/users", mw(h.GetAll))
	r.GET("/api/users/:id", mw(h.GetByID))
	r.POST("/api/users", mw(h.Create))
	r.PUT("/api/users/:id", mw(h.Update))
	r.DELETE("/api/users/:id", mw(h.Delete))
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	users := h.Service.GetAll()
	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := parseID(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.Service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var user models.User
	if err := utils.DecodeJSON(r, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created := h.Service.Create(user)
	// Асинхронное логирование и нотификация не блокируют ответ
	go h.Logger.Log("CREATE", created.ID)
	go h.Notify.Send(created.ID, "created")
	utils.WriteJSON(w, http.StatusCreated, created)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := parseID(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var payload models.User
	if err := utils.DecodeJSON(r, &payload); err != nil {
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
	utils.WriteJSON(w, http.StatusOK, updated)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := parseID(ps.ByName("id"))
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
