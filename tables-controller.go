package controllers

import (
	"fmt"
	"intership/models"
	"intership/utils"
	"net/http"
	"strconv"
	"strings"
	_"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	_ "github.com/jmoiron/sqlx"
)

var table_columns = []string{
	"id",
	"vendor_id",
	"name",
	"is_available",
	"customer_id",
	"is_needs_service",

}

// IndexTableHandler handles GET requests to fetch all tables
func IndexTableHandler(w http.ResponseWriter, r *http.Request) {
	var tables []models.Table
	query, args, err := QB.Select(strings.Join(table_columns, ", ")).From("tables").ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Select(&tables, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, tables)
}

// ShowTableHandler handles GET requests to fetch a single table by ID
func ShowTableHandler(w http.ResponseWriter, r *http.Request) {
	var table models.Table
	id := r.PathValue("id")
	query, args, err := QB.Select(strings.Join(table_columns, ", ")).From("tables").Where("id = ?", id).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&table, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, table)
}

// CreateTableHandler handles POST requests to create a new table
func CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	var table models.Table
	if r.FormValue("name") == "" || r.FormValue("vendor_id") == "" {
		utils.HandelError(w, http.StatusBadRequest, "Name and vendor_id are required")
		return
	}

	vendorID, err := uuid.Parse(r.FormValue("vendor_id")) // Convert string to uuid.UUID
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid vendor_id format")
		return
	}

	table.ID = uuid.New() // Generate new UUID
	table.Name = r.FormValue("name")
	table.VendorID = vendorID // Set vendor_id from request
	if r.FormValue("is_available") != "" {
		table.IsAvailable, _ = strconv.ParseBool(r.FormValue("is_available"))
	}
	if r.FormValue("customer_id") != "" {
		customerID, err := uuid.Parse(r.FormValue("customer_id"))
		if err == nil {
			table.CustomerID = &customerID // Set customer_id if provided
		}
	}
	if r.FormValue("is_needs_service") != "" {
		table.IsNeedsService, _ = strconv.ParseBool(r.FormValue("is_needs_service"))
	}

	query, args, err := QB.Insert("tables").
		Columns("id", "vendor_id", "name", "is_available", "customer_id", "is_needs_service").
		Values(table.ID, table.VendorID, table.Name, table.IsAvailable, table.CustomerID, table.IsNeedsService).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(table_columns, ", "))).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query: "+err.Error())
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&table); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error creating table: "+err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusCreated, table)
}

// UpdateTableHandler handles PUT requests to update an existing table
func UpdateTableHandler(w http.ResponseWriter, r *http.Request) {
	var table models.Table
	id := r.PathValue("id")

	query, args, err := QB.Select(strings.Join(table_columns, ", ")).From("tables").Where("id = ?", id).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&table, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update fields if provided
	if r.FormValue("name") != "" {
		table.Name = r.FormValue("name")
	}
	if r.FormValue("vendor_id") != "" {
		vendorID, err := uuid.Parse(r.FormValue("vendor_id"))
		if err != nil {
			utils.HandelError(w, http.StatusBadRequest, "Invalid vendor_id format")
			return
		}
		table.VendorID = vendorID // Update vendor_id as necessary
	}
	if r.FormValue("is_available") != "" {
		table.IsAvailable, _ = strconv.ParseBool(r.FormValue("is_available"))
	}
	if r.FormValue("customer_id") != "" {
		customerID, err := uuid.Parse(r.FormValue("customer_id"))
		if err == nil {
			table.CustomerID = &customerID // Update customer_id as necessary
		}
	}
	if r.FormValue("is_needs_service") != "" {
		table.IsNeedsService, _ = strconv.ParseBool(r.FormValue("is_needs_service"))
	}

	query, args, err = QB.Update("tables").
		Set("vendor_id", table.VendorID).
		Set("name", table.Name).
		Set("is_available", table.IsAvailable).
		Set("customer_id", table.CustomerID).
		Set("is_needs_service", table.IsNeedsService).
		Where(squirrel.Eq{"id": table.ID}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(table_columns, ", "))).
		ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query: "+err.Error())
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&table); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error updating table: "+err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, table)
}

// DeleteTableHandler handles DELETE requests to remove a table
func DeleteTableHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	query, args, err := QB.Delete("tables").Where("id=?", id).Suffix("RETURNING id").ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error deleting table: "+err.Error())
		return
	}

	var deletedID string
	if err := db.QueryRowx(query, args...).Scan(&deletedID); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error deleting table: "+err.Error())
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, fmt.Sprintf("Table with ID %s deleted", deletedID))
}
