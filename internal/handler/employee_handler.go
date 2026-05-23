package handler

import (
	"OrgAPI/internal/dto"
	"OrgAPI/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

type EmployeeHandler struct {
	service EmployeeService
}

func NewEmployeeHandler(service EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service: service,
	}
}

func (h *EmployeeHandler) CreateEmployee(
	w http.ResponseWriter,
	r *http.Request,
) {

	idParam := r.PathValue("id")

	id64, err := strconv.ParseInt(
		idParam,
		10,
		64,
	)

	if err != nil {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			"invalid id",
		)
		return
	}

	var req dto.CreateEmployeeRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			"invalid json",
		)
		return
	}

	response, err := h.service.CreateEmployee(
		r.Context(),
		int64(id64),
		req,
	)

	if err != nil {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	utils.WriteJSON(
		w,
		http.StatusOK,
		response,
	)
}
