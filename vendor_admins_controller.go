package controllers

import (
	"intership/models"
	"intership/utils"
	"net/http"
	_"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

// CreateVendorAdminHandler handles the creation of a vendor admin
func CreateVendorAdminHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.FormValue("user_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	vendorID, err := uuid.Parse(r.FormValue("vendor_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid vendor_id format")
		return
	}

	vendorAdmin := models.VendorAdmin{
		UserID:   userID,
		VendorID: vendorID,
	}

	query, args, err := QB.Insert("vendor_admins").
		Columns("user_id", "vendor_id").
		Values(vendorAdmin.UserID, vendorAdmin.VendorID).
		Suffix("RETURNING user_id, vendor_id").
		ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	err = db.Get(&vendorAdmin, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error inserting vendor admin")
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, vendorAdmin)
}

// IndexVendorAdminsHandler handles the listing of vendor admins
func IndexVendorAdminsHandler(w http.ResponseWriter, r *http.Request) {
	var vendorAdmins []models.VendorAdmin
	query, args, err := QB.Select("user_id", "vendor_id").
		From("vendor_admins").ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	err = db.Select(&vendorAdmins, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error fetching vendor admins")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, vendorAdmins)
}

// ShowVendorAdminHandler handles the retrieval of a vendor admin by ID
func ShowVendorAdminHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	vendorID, err := uuid.Parse(r.PathValue("vendor_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid vendor_id format")
		return
	}

	var vendorAdmin models.VendorAdmin
	query, args, err := QB.Select("user_id", "vendor_id").
		From("vendor_admins").
		Where(squirrel.Eq{"user_id": userID, "vendor_id": vendorID}).ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	err = db.Get(&vendorAdmin, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusNotFound, "Vendor admin not found")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, vendorAdmin)
}

// UpdateVendorAdminHandler handles the updating of a vendor admin
func UpdateVendorAdminHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	vendorID, err := uuid.Parse(r.PathValue("vendor_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid vendor_id format")
		return
	}

	var vendorAdmin models.VendorAdmin
	query, args, err := QB.Select("user_id", "vendor_id").
		From("vendor_admins").
		Where(squirrel.Eq{"user_id": userID, "vendor_id": vendorID}).ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	err = db.Get(&vendorAdmin, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusNotFound, "Vendor admin not found")
		return
	}

	// Update only if provided
	if r.FormValue("vendor_id") != "" {
		vendorID, err = uuid.Parse(r.FormValue("vendor_id"))
		if err != nil {
			utils.HandelError(w, http.StatusBadRequest, "Invalid vendor_id format")
			return
		}
		vendorAdmin.VendorID = vendorID
	}

	// You can add more fields here to update, if necessary

	query, args, err = QB.Update("vendor_admins").
		Set("vendor_id", vendorAdmin.VendorID).
		Where(squirrel.Eq{"user_id": userID, "vendor_id": vendorID}).
		ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building update query")
		return
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error updating vendor admin")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "Vendor admin updated successfully")
}

// DeleteVendorAdminHandler handles the deletion of a vendor admin
func DeleteVendorAdminHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	vendorID, err := uuid.Parse(r.PathValue("vendor_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid vendor_id format")
		return
	}

	query, args, err := QB.Delete("vendor_admins").
		Where(squirrel.Eq{"user_id": userID, "vendor_id": vendorID}).
		ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building delete query")
		return
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error deleting vendor admin")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "Vendor admin deleted successfully")
}
