package controllers

import (
	"fmt"
	"intership/models"
	"intership/utils"
	"net/http"
	"strings"

	_"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var cartItemColumns = []string{
	"cart_id",
	"item_id",
	"quantity",
}

// IndexCartItemsHandler handles GET requests to fetch all cart items
func IndexCartItemsHandler(w http.ResponseWriter, r *http.Request) {
	var cartItems []models.CartItem
	query, args, err := QB.Select(strings.Join(cartItemColumns, ", ")).From("cart_items").ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Select(&cartItems, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, cartItems)
}

// ShowCartItemHandler handles GET requests to fetch a single cart item by cart_id and item_id
func ShowCartItemHandler(w http.ResponseWriter, r *http.Request) {
	var cartItem models.CartItem
	cartID := r.PathValue("cart_id")
	itemID := r.PathValue("item_id")

	query, args, err := QB.Select(strings.Join(cartItemColumns, ", ")).
		From("cart_items").
		Where("cart_id = ? AND item_id = ?", cartID, itemID).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&cartItem, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, cartItem)
}

// CreateCartItemHandler handles POST requests to create a new cart item
func CreateCartItemHandler(w http.ResponseWriter, r *http.Request) {
	var cartItem models.CartItem
	if r.FormValue("cart_id") == "" || r.FormValue("item_id") == "" || r.FormValue("quantity") == "" {
		utils.HandelError(w, http.StatusBadRequest, "Cart ID, Item ID, and Quantity are required")
		return
	}

	cartID, err := uuid.Parse(r.FormValue("cart_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid cart ID format")
		return
	}

	itemID, err := uuid.Parse(r.FormValue("item_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid item ID format")
		return
	}

	quantity := r.FormValue("quantity")
	cartItem.CartID = cartID
	cartItem.ItemID = itemID
	cartItem.Quantity = utils.ParseQuantity(quantity)

	query, args, err := QB.Insert("cart_items").
		Columns("cart_id", "item_id", "quantity").
		Values(cartItem.CartID, cartItem.ItemID, cartItem.Quantity).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(cartItemColumns, ", "))).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := db.QueryRowx(query, args...).StructScan(&cartItem); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusCreated, cartItem)
}

// UpdateCartItemHandler handles PUT requests to update an existing cart item
func UpdateCartItemHandler(w http.ResponseWriter, r *http.Request) {
	var cartItem models.CartItem
	cartID := r.PathValue("cart_id")
	itemID := r.PathValue("item_id")

	query, args, err := QB.Select(strings.Join(cartItemColumns, ", ")).
		From("cart_items").
		Where("cart_id = ? AND item_id = ?", cartID, itemID).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&cartItem, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update quantity if provided
	if r.FormValue("quantity") != "" {
		quantity := r.FormValue("quantity")
		cartItem.Quantity = utils.ParseQuantity(quantity)
	}

	query, args, err = QB.Update("cart_items").
		Set("quantity", cartItem.Quantity).
		Where("cart_id = ? AND item_id = ?", cartItem.CartID, cartItem.ItemID).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(cartItemColumns, ", "))).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query: "+err.Error())
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&cartItem); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error updating cart item: "+err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, cartItem)
}

// DeleteCartItemHandler handles DELETE requests to remove a cart item
func DeleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cart_id")
	itemID := r.PathValue("item_id")

	query, args, err := QB.Delete("cart_items").
		Where("cart_id = ? AND item_id = ?", cartID, itemID).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building delete query: "+err.Error())
		return
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error deleting cart item: "+err.Error())
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "Cart item deleted")
}
