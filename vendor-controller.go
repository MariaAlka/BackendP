package controllers

import (
	_ "context"
	"fmt"
	"intership/models"
	"intership/utils"
	_ "log"
	"net/http"
	_"os"
	_ "strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)



func SetDBv(database *sqlx.DB) {
	db = database

}



func IndexVendorHandler(w http.ResponseWriter, r *http.Request) {
	var vendors []models.Vendor
	query, args, err := QB.Select(strings.Join(vendor_columns, ", ")).
		From("vendors").ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Select(&vendors, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, vendors)
}

func ShowVendorHandler(w http.ResponseWriter, r *http.Request) {
	var vendor models.Vendor
	id := r.PathValue("id")
	query, args, err := QB.Select(strings.Join(vendor_columns, ", ")).
		From("vendors").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&vendor, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, vendor)
}


func UpdateVendorHandler(w http.ResponseWriter, r *http.Request) {
	var vendor models.Vendor
	id := r.PathValue("id")
	query, args, err := QB.Select(strings.Join(vendor_columns, ", ")).
		From("vendors").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&vendor, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//update user
	if r.FormValue("name") != "" {
		vendor.Name = r.FormValue("name")
	}
	if r.FormValue("description") != "" {
		vendor.Description = r.FormValue("description")
	}
	file, fileHeader, err := r.FormFile("img")
	if err != nil && err != http.ErrMissingFile {
		utils.HandelError(w, http.StatusBadRequest, "Invalid file")
		return
	} else if err == nil {
		defer file.Close()
		imageName, err := utils.SaveImageFile(file, "vendors", fileHeader.Filename)
		if err != nil {
			utils.HandelError(w, http.StatusInternalServerError, "Error saving image")

		}
		vendor.Img = &imageName
	}
query , args , err = QB.Update("vendors").
Set("img",vendor.Img).
Set("name",vendor.Name).
Set("description",vendor.Description).
Set("updated_at", time.Now()).
Where(squirrel.Eq{"id":vendor.ID}).
Suffix(fmt.Sprintf("RETURNING %s", strings.Join(vendor_columns, ", "))).ToSql()

if err != nil {
	utils.HandelError(w, http.StatusInternalServerError, "Error bulding query")
    return
}
if err := db.QueryRowx(query, args...).StructScan(&vendor); err != nil {
	utils.HandelError(w, http.StatusInternalServerError, "Error creating user"+err.Error())
	return
}
utils.SendJSONResponse(w, http.StatusCreated, vendor)

	utils.SendJSONResponse(w, http.StatusOK, vendor)
}


func DeleteVendorHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	query , args , err := QB.Delete("vendors").
	Where("id=?",id).
	Suffix("RETURNING img").
	ToSql()
	if err!= nil {
        utils.HandelError(w, http.StatusInternalServerError, "Error deleting user:"+err.Error())
        return
    }

	var img *string
	if err := db.QueryRowx(query, args...).Scan(&img); err!= nil {
        utils.HandelError(w, http.StatusInternalServerError, "Error getting image"+err.Error())
        return
    }

	utils.SendJSONResponse(w, http.StatusOK, "vendor deleted")
}
