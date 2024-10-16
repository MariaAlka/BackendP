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
	_ "github.com/jmoiron/sqlx"
)

var item_columns = []string{
	"id",
	"vendor_id",
	"name",
	"price",
	"img",
	"created_at",
	"updated_at",
	fmt.Sprintf("CASE WHEN NULLIF(img, '') IS NOT NULL THEN FORMAT('%s/%%s', img) ELSE NULL END AS img", Domain),
}

// IndexItemHandler handles GET requests to fetch all items
func IndexItemHandler(w http.ResponseWriter, r *http.Request) {
	var items []models.Item
	query, args, err := QB.Select(strings.Join(item_columns, ", ")).From("items").ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Select(&items, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, items)
}

// ShowItemHandler handles GET requests to fetch a single item by ID
func ShowItemHandler(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	id := r.PathValue("id")
	query, args, err := QB.Select(strings.Join(item_columns, ", ")).From("items").Where("id = ?", id).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&item, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, item)
}

// CreateItemHandler handles POST requests to create a new item
func CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if r.FormValue("name") == "" || r.FormValue("price") == "" || r.FormValue("vendor_id") == "" {
		utils.HandelError(w, http.StatusBadRequest, "Name, price, and vendor_id are required")
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid price format")
		return
	}

	vendorID, err := uuid.Parse(r.FormValue("vendor_id")) // Convert string to uuid.UUID
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid vendor_id format")
		return
	}

	item.ID = uuid.New() // Generate new UUID
	item.Name = r.FormValue("name")
	item.Price = price
	item.VendorID = vendorID // Set vendor_id from request

	// Handle image upload
	file, fileHeader, err := r.FormFile("img")
	if err != nil && err != http.ErrMissingFile {
		utils.HandelError(w, http.StatusBadRequest, "Invalid file")
		return
	} else if err == nil {
		defer file.Close()
		imageName, err := utils.SaveImageFile(file, "items", fileHeader.Filename) // Save image in the "items" directory
		if err != nil {
			utils.HandelError(w, http.StatusInternalServerError, "Error saving image")
			return
		}
		item.Img = &imageName // Store image path in the item
	}

	// Build SQL query for inserting item
	query, args, err := QB.Insert("items").Columns("id", "vendor_id", "name", "price", "img").Values(item.ID, item.VendorID, item.Name, item.Price, item.Img).Suffix(fmt.Sprintf("RETURNING %s", strings.Join(item_columns, ", "))).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query: "+err.Error())
		return
	}

	// Execute query and scan the result into the item struct
	if err := db.QueryRowx(query, args...).StructScan(&item); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error creating item: "+err.Error())
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, item)
}

// UpdateItemHandler handles PUT requests to update an existing item
func UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	id := r.PathValue("id")

	query, args, err := QB.Select(strings.Join(item_columns, ", ")).From("items").Where("id = ?", id).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&item, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Update fields if provided
	if r.FormValue("name") != "" {
		item.Name = r.FormValue("name")
	}
	if r.FormValue("price") != "" {
		price, err := strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			utils.HandelError(w, http.StatusBadRequest, "Invalid price format")
			return
		}
		item.Price = price
	}
	if r.FormValue("vendor_id") != "" {
		vendorID, err := uuid.Parse(r.FormValue("vendor_id")) // Convert string to uuid.UUID
		if err != nil {
			utils.HandelError(w, http.StatusBadRequest, "Invalid vendor_id format")
			return
		}
		item.VendorID = vendorID // Update vendor_id as necessary
	}
	if r.FormValue("img") != "" {
		img := r.FormValue("img") // Update image as necessary
		item.Img = &img           // Update image path as necessary
	}

	query, args, err = QB.Update("items").
		Set("name", item.Name).
		Set("price", item.Price).
		Set("vendor_id", item.VendorID). // Ensure vendor_id is updated
		Set("img", item.Img).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": item.ID}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(item_columns, ", "))).
		ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query: "+err.Error())
		return
	}

	if err := db.QueryRowx(query, args...).StructScan(&item); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error updating item: "+err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, item)
}

// DeleteItemHandler handles DELETE requests to remove an item
func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	query, args, err := QB.Delete("items").Where("id=?", id).Suffix("RETURNING img").ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error deleting item: "+err.Error())
		return
	}

	var img *string
	if err := db.QueryRowx(query, args...).Scan(&img); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error getting image: "+err.Error())
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "Item deleted")
}
