package controllers

import (
	"fmt"
	"intership/models"
	"intership/utils"
	"net/http"
	_ "os"
	"strings"
	"time"

	_ "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func SignUpVendorHandler(w http.ResponseWriter, r *http.Request) {
    vendor := models.Vendor{
        ID:         uuid.New(),
        Name:       r.FormValue("name"),
        Description:      r.FormValue("description"),
        Created_at: time.Now(),
        Updated_at: time.Now(),
    }
    file, fileHeader, err := r.FormFile("img")
    if err != nil && err != http.ErrMissingFile {
        utils.HandelError(w, http.StatusBadRequest, "Invalid file")
        return
    } else if err == nil {
        defer file.Close()
        imageName, err := utils.SaveImageFile(file, "vendors", fileHeader.Filename)
fmt.Println(imageName)
        if err != nil {
            utils.HandelError(w, http.StatusInternalServerError, "Error saving image")

        }
        vendor.Img = &imageName
    }

    query, args, err := QB.
        Insert("vendors").Columns("id", "img", "name", "description").
        Values(vendor.ID, vendor.Img, vendor.Name, vendor.Description).
        Suffix(fmt.Sprintf("RETURNING %s", strings.Join(vendor_columns, ", "))).ToSql()
    if err != nil {
        utils.HandelError(w, http.StatusInternalServerError, "Error generate query ")
        return
    }
    if err := db.QueryRowx(query, args...).StructScan(&vendor); err != nil {
        utils.HandelError(w, http.StatusInternalServerError, "Error creating user"+err.Error())
        return
    }
    utils.SendJSONResponse(w, http.StatusCreated, vendor)

}

