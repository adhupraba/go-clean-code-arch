package cart

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/adhupraba/ecom/service/auth"
	"github.com/adhupraba/ecom/types"
	"github.com/adhupraba/ecom/utils"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store, productStore, userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUserIdFromContext(r.Context())

	var payload types.CartCheckoutPayload

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, err)
		return
	}

	productIds, err := getCartItemIDs(payload.Items)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	products, err := h.productStore.GetProductsByIDs(productIds)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderId, totalPrice, err := h.createOrder(products, payload.Items, userId)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]any{
		"totalPrice": totalPrice,
		"orderId":    orderId,
	})
}
