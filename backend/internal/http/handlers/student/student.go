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
			err := response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			if err != nil {
				return
			}
			return
		}

		if err != nil {
			err := response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			if err != nil {
				return
			}
			return
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {

			var validationErrs validator.ValidationErrors
			errors.As(err, &validationErrs)
			_ = response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validationErrs))
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

		_ = response.WriteJSON(w, http.StatusCreated, map[string]int64{"success": lastId})

	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student for", slog.String("user id", fmt.Sprint(id)))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("error parsing request", slog.String("id", id))
			_ = response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("error getting student", slog.String("id", id))
			_ = response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		_ = response.WriteJSON(w, http.StatusOK, student)
	}
}

func GetAllStudents(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")

		students, err := storage.GetAllStudents()
		if err != nil {
			if err := response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err)); err != nil {
				return
			}
		}

		if err := response.WriteJSON(w, http.StatusOK, students); err != nil {
			return
		}

	}
}

//func UpdateStudentById(storage storage.Storage) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		id := r.PathValue("id")
//		slog.Info("updating a student", slog.String("user id", fmt.Sprint(id)))
//
//		var student types.Student
//		err := json.NewDecoder(r.Body).Decode(&student)
//		if err != nil {
//			if resperr := response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err)); resperr != nil {
//				return
//			}
//			return
//		}
//
//	}
//}
