package controllers

import (
	"fmt"
	"intership/models"
	"intership/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var cartColumns = []string{
	"id",
	"total_price",
	"quantity",
	"vendor_id",
	"created_at",
	"updated_at",
}

// IndexCartHandler handles GET requests to fetch all carts
func IndexCartHandler(w http.ResponseWriter, r *http.Request) {
	var carts []models.Cart
	query, args, err := QB.Select(strings.Join(cartColumns, ", ")).From("carts").ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Select(&carts, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, carts)
}

// ShowCartHandler handles GET requests to fetch a single cart by ID
func ShowCartHandler(w http.ResponseWriter, r *http.Request) {
	var cart models.Cart
	id := r.PathValue("id")
	query, args, err := QB.Select(strings.Join(cartColumns, ", ")).From("carts").Where("id = ?", id).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&cart, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, cart)
}

// CreateCartHandler handles POST requests to create a new cart
func CreateCartHandler(w http.ResponseWriter, r *http.Request) {
	var cart models.Cart
	if r.FormValue("id") == "" || r.FormValue("total_price") == "" || r.FormValue("quantity") == "" {
		utils.HandelError(w, http.StatusBadRequest, "User ID, total price, and quantity are required")
		return
	}

	cartID, err := uuid.Parse(r.FormValue("id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid id format")
		return
	}

	totalPrice, err := strconv.ParseFloat(r.FormValue("total_price"), 64)
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid total_price format")
		return
	}

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid quantity format")
		return
	}

	cart.ID = cartID
	cart.TotalPrice = totalPrice
	cart.Quantity = quantity

	// Optional vendor ID
	if r.FormValue("vendor_id") != "" {
		vendorID, err := uuid.Parse(r.FormValue("vendor_id"))
		if err != nil {
			utils.HandelError(w, http.StatusBadRequest, "Invalid vendor_id format")
			return
		}
		cart.VendorID = vendorID
	}

	query, args, err := QB.Insert("carts").
		Columns("id", "total_price", "quantity", "vendor_id", "created_at", "updated_at").
		Values(cart.ID, cart.TotalPrice, cart.Quantity, cart.VendorID, "NOW()", "NOW()").
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(cartColumns, ", "))).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&cart); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusCreated, cart)
}

// UpdateCartHandler handles PUT requests to update an existing cart
func UpdateCartHandler(w http.ResponseWriter, r *http.Request) {
	var cart models.Cart
	id := r.PathValue("id")

	query, args, err := QB.Select(strings.Join(cartColumns, ", ")).From("carts").Where("id = ?", id).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&cart, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if r.FormValue("total_price") != "" {
		totalPrice, err := strconv.ParseFloat(r.FormValue("total_price"), 64)
		if err != nil {
			utils.HandelError(w, http.StatusBadRequest, "Invalid total_price format")
			return
		}
		cart.TotalPrice = totalPrice
	}
	if r.FormValue("quantity") != "" {
		quantity, err := strconv.Atoi(r.FormValue("quantity"))
		if err != nil {
			utils.HandelError(w, http.StatusBadRequest, "Invalid quantity format")
			return
		}
		cart.Quantity = quantity
	}
	if r.FormValue("vendor_id") != "" {
		vendorID, err := uuid.Parse(r.FormValue("vendor_id"))
		if err != nil {
			utils.HandelError(w, http.StatusBadRequest, "Invalid vendor_id format")
			return
		}
		cart.VendorID = vendorID
	}

	query, args, err = QB.Update("carts").
		Set("total_price", cart.TotalPrice).
		Set("quantity", cart.Quantity).
		Set("vendor_id", cart.VendorID).
		Where(squirrel.Eq{"id": cart.ID}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(cartColumns, ", "))).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&cart); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, cart)
}

// DeleteCartHandler handles DELETE requests to remove a cart
func DeleteCartHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	query, args, err := QB.Delete("carts").Where("id=?", id).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if _, err := db.Exec(query, args...); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "Cart deleted")
}
