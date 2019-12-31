package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/dqkcode/notes/internal/app/types"
)

type (
	service interface {
		Register(ctx context.Context, user RegisterRequest) (string, error)
		Login(ctx context.Context, user LoginRequest) (string, error)
	}
	Handler struct {
		srv service
	}
)

func NewHandler(service_input service) *Handler {
	return &Handler{
		srv: service_input,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.srv.Register(r.Context(), req)
	if err == ErrUserAlreadyExist {
		json.NewEncoder(w).Encode(types.Response{
			Code:  types.UserAlreadyExist,
			Error: err.Error(),
		})

	}
	if err != nil {
		json.NewEncoder(w).Encode(types.Response{
			Code: types.ErrorDB,
		})
	}
	json.NewEncoder(w).Encode(types.Response{
		Code: types.CodeSuccess,
		Data: map[string]interface{}{
			"id": id,
		},
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logrus.Errorf("request is invalid type")
		logrus.Infof("request : %s", req)
	}

}
