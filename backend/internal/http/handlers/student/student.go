package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Bomjan/gmarket/backend/internal/storage"
	"github.com/Bomjan/gmarket/backend/internal/types"
	"github.com/Bomjan/gmarket/backend/internal/utils/response"
	"github.com/go-playground/validator"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a Student")

		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {

			validationErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validationErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created successfully", slog.String("user id", fmt.Sprint(lastId)))
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
		}

		response.WriteJSON(w, http.StatusCreated, map[string]int64{"success": lastId})

	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student for", slog.String("user id", fmt.Sprint(id)))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("error parsing request", slog.String("id", id))
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("error getting student", slog.String("id", id))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, student)
	}
}
