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

var orderItemColumns = []string{
	"id",
	"order_id",
	"item_id",
	"quantity",
	"price",
}

// IndexOrderItemHandler handles GET requests to fetch all order_items
func IndexOrderItemHandler(w http.ResponseWriter, r *http.Request) {
	var orderItems []models.OrderItem
	query, args, err := QB.Select(strings.Join(orderItemColumns, ", ")).From("order_items").ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Select(&orderItems, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, orderItems)
}

// ShowOrderItemHandler handles GET requests to fetch a single order_item by ID
func ShowOrderItemHandler(w http.ResponseWriter, r *http.Request) {
	var orderItem models.OrderItem
	id := r.PathValue("id")
	query, args, err := QB.Select(strings.Join(orderItemColumns, ", ")).From("order_items").Where("id = ?", id).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&orderItem, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, orderItem)
}

// CreateOrderItemHandler handles POST requests to create a new order_item
func CreateOrderItemHandler(w http.ResponseWriter, r *http.Request) {
	var orderItem models.OrderItem
	if r.FormValue("order_id") == "" || r.FormValue("item_id") == "" || r.FormValue("quantity") == "" || r.FormValue("price") == "" {
		utils.HandelError(w, http.StatusBadRequest, "Order ID, Item ID, quantity, and price are required")
		return
	}

	orderID, err := uuid.Parse(r.FormValue("order_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid order_id format")
		return
	}

	itemID, err := uuid.Parse(r.FormValue("item_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid item_id format")
		return
	}

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid quantity format")
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid price format")
		return
	}

	orderItem.ID = uuid.New()
	orderItem.OrderID = orderID
	orderItem.ItemID = itemID
	orderItem.Quantity = quantity
	orderItem.Price = price

	query, args, err := QB.Insert("order_items").
		Columns("id", "order_id", "item_id", "quantity", "price").
		Values(orderItem.ID, orderItem.OrderID, orderItem.ItemID, orderItem.Quantity, orderItem.Price).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(orderItemColumns, ", "))).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&orderItem); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusCreated, orderItem)
}

// UpdateOrderItemHandler handles PUT requests to update an existing order_item
func UpdateOrderItemHandler(w http.ResponseWriter, r *http.Request) {
	var orderItem models.OrderItem
	id := r.PathValue("id")

	query, args, err := QB.Select(strings.Join(orderItemColumns, ", ")).From("order_items").Where("id = ?", id).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&orderItem, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if r.FormValue("quantity") != "" {
		quantity, err := strconv.Atoi(r.FormValue("quantity"))
		if err != nil {
			utils.HandelError(w, http.StatusBadRequest, "Invalid quantity format")
			return
		}
		orderItem.Quantity = quantity
	}
	if r.FormValue("price") != "" {
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			utils.HandelError(w, http.StatusBadRequest, "Invalid price format")
			return
		}
		orderItem.Price = price
	}

	query, args, err = QB.Update("order_items").
		Set("quantity", orderItem.Quantity).
		Set("price", orderItem.Price).
		Where(squirrel.Eq{"id": orderItem.ID}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(orderItemColumns, ", "))).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&orderItem); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, orderItem)
}

// DeleteOrderItemHandler handles DELETE requests to remove an order_item
func DeleteOrderItemHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	query, args, err := QB.Delete("order_items").Where("id=?", id).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if _, err := db.Exec(query, args...); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "Order item deleted")
}
