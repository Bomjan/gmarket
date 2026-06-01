package products

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

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
