package controllers

import (
	"fmt"
	"intership/models"
	"intership/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// Order columns for SQL queries
var order_columns = []string{
	"id",
	"total_order_cost",
	"customer_id",
	"vendor_id",
	"status",
	"created_at",
	"updated_at",
}

// IndexOrderHandler handles GET requests to fetch all orders
func IndexOrderHandler(w http.ResponseWriter, r *http.Request) {
	var orders []models.Order
	query, args, err := QB.Select(strings.Join(order_columns, ", ")).From("orders").ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Select(&orders, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, orders)
}

// ShowOrderHandler handles GET requests to fetch a single order by ID
func ShowOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	id := r.PathValue("id")
	query, args, err := QB.Select(strings.Join(order_columns, ", ")).From("orders").Where("id = ?", id).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&order, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, order)
}

// CreateOrderHandler handles POST requests to create a new order
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if r.FormValue("total_order_cost") == "" || r.FormValue("customer_id") == "" || r.FormValue("vendor_id") == "" || r.FormValue("status") == "" {
		utils.HandelError(w, http.StatusBadRequest, "Total order cost, customer_id, vendor_id, and status are required")
		return
	}

	totalOrderCost, err := strconv.ParseFloat(r.FormValue("total_order_cost"), 64)
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid total order cost format")
		return
	}

	customerID, err := uuid.Parse(r.FormValue("customer_id")) // Convert string to uuid.UUID
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid customer_id format")
		return
	}

	vendorID, err := uuid.Parse(r.FormValue("vendor_id")) // Convert string to uuid.UUID
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid vendor_id format")
		return
	}

	order.ID = uuid.New() // generate new UUID
	order.TotalOrderCost = totalOrderCost
	order.CustomerID = customerID
	order.VendorID = vendorID
	order.Status = models.OrderStatus(r.FormValue("status")) // Convert status to enum type
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	query, args, err := QB.Insert("orders").Columns("id", "total_order_cost", "customer_id", "vendor_id", "status", "created_at", "updated_at").
		Values(order.ID, order.TotalOrderCost, order.CustomerID, order.VendorID, order.Status, order.CreatedAt, order.UpdatedAt).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(order_columns, ", "))).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query: "+err.Error())
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&order); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error creating order: "+err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusCreated, order)
}

// UpdateOrderHandler handles PUT requests to update an existing order
func UpdateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	id := r.PathValue("id")

	query, args, err := QB.Select(strings.Join(order_columns, ", ")).From("orders").Where("id = ?", id).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&order, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update fields if provided
	if r.FormValue("total_order_cost") != "" {
		totalOrderCost, err := strconv.ParseFloat(r.FormValue("total_order_cost"), 64)
		if err != nil {
			utils.HandelError(w, http.StatusBadRequest, "Invalid total order cost format")
			return
		}
		order.TotalOrderCost = totalOrderCost
	}
	if r.FormValue("customer_id") != "" {
		customerID, err := uuid.Parse(r.FormValue("customer_id")) // Convert string to uuid.UUID
		if err != nil {
			utils.HandelError(w, http.StatusBadRequest, "Invalid customer_id format")
			return
		}
		order.CustomerID = customerID
	}
	if r.FormValue("vendor_id") != "" {
		vendorID, err := uuid.Parse(r.FormValue("vendor_id")) // Convert string to uuid.UUID
		if err != nil {
			utils.HandelError(w, http.StatusBadRequest, "Invalid vendor_id format")
			return
		}
		order.VendorID = vendorID
	}
	if r.FormValue("status") != "" {
		order.Status = models.OrderStatus(r.FormValue("status")) // Update status as necessary
	}

	order.UpdatedAt = time.Now()

	query, args, err = QB.Update("orders").
		Set("total_order_cost", order.TotalOrderCost).
		Set("customer_id", order.CustomerID).
		Set("vendor_id", order.VendorID).
		Set("status", order.Status).
		Set("updated_at", order.UpdatedAt).
		Where(squirrel.Eq{"id": order.ID}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(order_columns, ", "))).
		ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query: "+err.Error())
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&order); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error updating order: "+err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, order)
}

// DeleteOrderHandler handles DELETE requests to remove an order
func DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	query, args, err := QB.Delete("orders").Where("id=?", id).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error deleting order: "+err.Error())
		return
	}

	if _, err := db.Exec(query, args...); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error deleting order: "+err.Error())
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "Order deleted")
}
