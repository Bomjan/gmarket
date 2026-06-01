package products

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Bomjan/gmarket/backend/internal/storage"
	"github.com/Bomjan/gmarket/backend/internal/types"
	"github.com/Bomjan/gmarket/backend/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating new product")
		// Get Product data
		var product types.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			if e := response.WriteJSON(w, http.StatusInternalServerError, err.Error()); e != nil {
				return
			}
			return
		}
		// Create Product
		lastId, err := storage.CreateProduct(product.Name, product.Price)
		if err != nil {
			if err := response.WriteJSON(w, http.StatusInternalServerError, err); err != nil {
				return
			}
			return
		}

		slog.Info("new product created", slog.String("id", fmt.Sprint(lastId)))

		// Return the data to the api
		if err := response.WriteJSON(w, http.StatusOK, map[string]int64{"success": lastId}); err != nil {
			return
		}
		return

	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log that we are here
		slog.Info("getting product by id")

		// get requested id
		id := r.PathValue("id")
		prodId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			if err := response.WriteJSON(w, http.StatusBadRequest, err.Error()); err != nil {
				return
			}
			return
		}
		// get the product by id
		slog.Info("getting product by id", slog.String("id", fmt.Sprint(prodId)))
		//var product types.Product
		product, err := storage.GetProductById(prodId)
		if err != nil {
			if err := response.WriteJSON(w, http.StatusInternalServerError, err.Error()); err != nil {
				return
			}
			return
		}
		// write to the response
		if err := response.WriteJSON(w, http.StatusOK, product); err != nil {
			return
		}
	}
}

func GetProducts(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// logging that this function is called
		slog.Info("getting all products")

		// get all products now
		var products []types.Product
		products, err := storage.GetAllProducts()
		if err != nil {
			if err := response.WriteJSON(w, http.StatusInternalServerError, err.Error()); err != nil {
				return
			}
			return
		}

		// send these products
		if err := response.WriteJSON(w, http.StatusOK, products); err != nil {
			if err := response.WriteJSON(w, http.StatusInternalServerError, err.Error()); err != nil {
				return
			}
			slog.Info("failed getting all products", slog.String("products", fmt.Sprint(products)))
			return
		}
	}
}

func DeleteProductById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("deleting product by id", slog.String("id", fmt.Sprint(id)))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			if err := response.WriteJSON(w, http.StatusBadRequest, err.Error()); err != nil {
				return
			}
			return
		}
		if err := storage.DeleteProductById(intId); err != nil {
			if err := response.WriteJSON(w, http.StatusInternalServerError, err.Error()); err != nil {
				return
			}
			return
		}

		if err := response.WriteJSON(w, http.StatusOK, map[string]string{"success": "success"}); err != nil {
			return
		}
	}
}

func UpdateProductById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			if err := response.WriteJSON(w, http.StatusBadRequest, err.Error()); err != nil {
				return
			}
			return
		}
		var product types.Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			if e := response.WriteJSON(w, http.StatusBadRequest, err.Error()); e != nil {
				return
			}
			return
		}

		if product.Id != intId {
			slog.Info("failed to update product", slog.String("urlId", fmt.Sprint(intId)), slog.String("bodyId", fmt.Sprint(product.Id)))
			if err := response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "the passed id should be same as provided in the body"}); err != nil {
				return
			}
			return
		}

		rowsAffected, err := storage.UpdateProductById(product.Id, product.Name, product.Price)
		if err != nil {
			if err := response.WriteJSON(w, http.StatusInternalServerError, err.Error()); err != nil {
				return
			}
			return
		}

		if err := response.WriteJSON(w, http.StatusOK, map[string]int64{"rows affected": rowsAffected}); err != nil {
			return
		}
	}
}
