package handler

import (
	"OrgAPI/internal/dto"
	"OrgAPI/internal/service"
	"OrgAPI/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

type DepartmentHandler struct {
	service service.DepartmentService
}

func NewDepartmentHandler(
	service service.DepartmentService,
) *DepartmentHandler {
	return &DepartmentHandler{
		service: service,
	}
}

func (h *DepartmentHandler) CreateDepartment(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.CreateDepartmentRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			"invalid json",
		)
		return
	}

	response, err := h.service.CreateDepartment(
		r.Context(),
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	utils.WriteJSON(
		w,
		http.StatusOK,
		response,
	)
}

func (h *DepartmentHandler) GetDepartment(
	w http.ResponseWriter,
	r *http.Request,
) {
	idParam := r.PathValue("id")

	id64, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			"invalid id",
		)
	}

	depth := 1

	depthQuery := r.URL.Query().Get("depth")

	if depthQuery != "" {
		parsedDepth, err := strconv.Atoi(depthQuery)

		if err != nil {
			utils.WriteError(
				w,
				http.StatusBadRequest,
				"invalid depth",
			)
		}

		depth = parsedDepth
	}

	includeEmployees := true

	includeQuery := r.URL.Query().Get(
		"include_employees",
	)

	if includeQuery == "false" {
		includeEmployees = false
	}

	response, err := h.service.GetDepartmentTree(
		r.Context(),
		uint(id64),
		depth,
		includeEmployees,
	)

	if err != nil {
		utils.WriteError(
			w,
			http.StatusNotFound,
			err.Error(),
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	utils.WriteJSON(
		w,
		http.StatusOK,
		response,
	)
}

func (h *DepartmentHandler) UpdateDepartment(
	w http.ResponseWriter,
	r *http.Request,
) {

	idParam := r.PathValue("id")

	id64, err := strconv.ParseUint(
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

	var req dto.UpdateDepartmentRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			"invalid json",
		)
		return
	}

	response, err := h.service.UpdateDepartment(
		r.Context(),
		uint(id64),
		req,
	)

	if err != nil {

		if err.Error() == "cycle detected" {

			utils.WriteError(
				w,
				http.StatusConflict,
				err.Error(),
			)
			return
		}

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

	utils.WriteJSON(
		w,
		http.StatusOK,
		response,
	)
}

func (h *DepartmentHandler) DeleteDepartment(
	w http.ResponseWriter,
	r *http.Request,
) {

	idParam := r.PathValue("id")

	id64, err := strconv.ParseUint(
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

	mode := r.URL.Query().Get("mode")

	reassignIDQuery := r.URL.Query().Get(
		"reassign_to_department_id",
	)

	req := dto.DeleteDepartmentRequest{
		Mode: mode,
	}
	if reassignIDQuery != "" {

		reassignID64, err := strconv.ParseUint(
			reassignIDQuery,
			10,
			64,
		)

		if err != nil {
			utils.WriteError(
				w,
				http.StatusBadRequest,
				"invalid reassign department id",
			)
			return
		}

		reassignID := uint(reassignID64)

		req.ReassignToDepartmentID = &reassignID
	}
	err = h.service.DeleteDepartment(
		r.Context(),
		uint(id64),
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

	w.WriteHeader(http.StatusNoContent)
}
