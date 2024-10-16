package controllers

import (
	_ "context"
	"fmt"
	"intership/models"
	"intership/utils"
	_ "log"
	"net/http"
	"os"
	_ "strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

var (
	db *sqlx.DB
	QB = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

func SetDB(database *sqlx.DB) {
	db = database

}

var (
	Domain = os.Getenv("DOMAIN")

	user_columns = []string{
		"id",
		"name",
		"email",
		"phone",
		"created_at",
		"updated_at",
		fmt.Sprintf("CASE WHEN NULLIF(img, '') IS NOT NULL THEN FORMAT('%s/%%s', img) ELSE NULL END AS img", Domain),
	}
	vendor_columns = []string{
		"id", "name", "description", "created_at", "updated_at",
		fmt.Sprintf("CASE WHEN NULLIF(img, '') IS NOT NULL THEN FORMAT('%s/%%s', img) ELSE NULL END AS img", Domain),
	}
)

func IndexUserHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	query, args, err := QB.Select(strings.Join(user_columns, ", ")).
		From("users").ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Select(&users, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, users)
}

func ShowUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	id := r.PathValue("id")
	query, args, err := QB.Select(strings.Join(user_columns, ", ")).
		From("users").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&user, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusOK, user)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	id := r.PathValue("id")
	query, args, err := QB.Select(strings.Join(user_columns, ", ")).
		From("users").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = db.Get(&user, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//update user
	if r.FormValue("name") != "" {
		user.Name = r.FormValue("name")
	}
	if r.FormValue("phone") != "" {
		user.Phone = r.FormValue("phone")
	}
	if r.FormValue("email") != "" {
		user.Email = r.FormValue("email")
	}
	if r.FormValue("password") != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			utils.HandelError(w, http.StatusInternalServerError, "Error hashing password")
			return
		}
		user.Password = hashedPassword
	}
	file, fileHeader, err := r.FormFile("img")
	if err != nil && err != http.ErrMissingFile {
		utils.HandelError(w, http.StatusBadRequest, "Invalid file")
		return
	} else if err == nil {
		defer file.Close()
		imageName, err := utils.SaveImageFile(file, "users", fileHeader.Filename)
		if err != nil {
			utils.HandelError(w, http.StatusInternalServerError, "Error saving image")

		}
		user.Img = &imageName
	}
	query, args, err = QB.Update("users").
		Set("img", user.Img).
		Set("name", user.Name).
		Set("email", user.Email).
		Set("phone", user.Phone).
		Set("password", user.Password).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": user.ID}).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(user_columns, ", "))).ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error bulding query")
		return
	}
	if err := db.QueryRowx(query, args...).StructScan(&user); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error creating user"+err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusCreated, user)

	utils.SendJSONResponse(w, http.StatusOK, user)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	query, args, err := QB.Delete("users").
		Where("id=?", id).
		Suffix("RETURNING img").
		ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error deleting user:"+err.Error())
		return
	}

	var img *string
	if err := db.QueryRowx(query, args...).Scan(&img); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error getting image"+err.Error())
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "user deleted")
}
