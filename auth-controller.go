package controllers

import (
	"database/sql"
	"fmt"
	"intership/models"
	"intership/utils"
	"net/http"
	_ "os"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	user := models.User{
		ID:         uuid.New(),
		Name:       r.FormValue("name"),
		Phone:      r.FormValue("phone"),
		Email:      r.FormValue("email"),
		Password:   r.FormValue("password"),
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}
	if user.Password == "" {
		utils.HandelError(w, http.StatusBadRequest, "Password is required")
		return
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

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	user.Password = hashedPassword

	query, args, err := QB.
		Insert("users").Columns("id", "img", "name", "phone", "email", "password").
		Values(user.ID, user.Img, user.Name, user.Phone, user.Email, user.Password).
		Suffix(fmt.Sprintf("RETURNING %s", strings.Join(user_columns, ", "))).ToSql()
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error generate query ")
		return
	}
	// fmt.Println("query", query)
	// fmt.Println("args", args)

	if err := db.QueryRowx(query, args...).StructScan(&user); err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error creating user"+err.Error())
		return
	}
	utils.SendJSONResponse(w, http.StatusCreated, user)

}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	var credentials struct{
		Email string 
		Password string 
	}
	credentials.Email = r.FormValue("email")
	credentials.Password = r.FormValue("password")

	var user models.User
	query , args , err := QB.Select("id","password").From("users").Where(squirrel.Eq{"email": credentials.Email}).ToSql()

	if err!= nil {
        utils.HandelError(w, http.StatusInternalServerError, "Error bulding query")
        return
    }
	if err := db.Get(&user,query,args...); err != nil {
		if err == sql.ErrNoRows {
			utils.HandelError(w, http.StatusNotFound, "User not found")
            return
		}
		utils.HandelError(w, http.StatusInternalServerError, "Error featching user")
		return
}
 if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
	    utils.HandelError(w, http.StatusInternalServerError, "Invalid credentials")
        return
    }


	tokenRsponse, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error generating JWT")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    tokenRsponse.Token,
		Path:     "/",
		Expires:  time.Now().UTC().Add(time.Hour * 24),
		HttpOnly: true,
	})

	utils.SendJSONResponse(w, http.StatusOK, tokenRsponse)
}

