package controllers

import (
	"fmt"
	"intership/models"
	"intership/utils"
	"net/http"
	"strconv"
	_ "time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func CreateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	// Parse user_id as UUID
	userID, err := uuid.Parse(r.FormValue("user_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}
	fmt.Print(r.FormValue("user_id"))
	// Parse role_id as integer
	roleID, err := strconv.Atoi(r.FormValue("role_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid role_id format")
		return
	}

	userRole := models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	// Insert user_role into the database
	query, args, err := QB.Insert("user_roles").
		Columns("user_id", "role_id").
		Values(userRole.UserID, userRole.RoleID).
		Suffix("RETURNING user_id, role_id").
		ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	err = db.Get(&userRole, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error inserting user role")
		return

	}

	utils.SendJSONResponse(w, http.StatusCreated, userRole)
}

func IndexUserRolesHandler(w http.ResponseWriter, r *http.Request) {
	var userRoles []models.UserRole
	query, args, err := QB.Select("user_id", "role_id").
		From("user_roles").ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	err = db.Select(&userRoles, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error fetching user roles")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, userRoles)
}

func ShowUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	var userRole models.UserRole

	// Parse user_id from request
	userID, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	// Parse role_id from request
	roleID, err := strconv.Atoi(r.PathValue("role_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid role_id format")
		return
	}

	query, args, err := QB.Select("user_id", "role_id").
		From("user_roles").
		Where(squirrel.Eq{"user_id": userID, "role_id": roleID}).ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building query")
		return
	}

	err = db.Get(&userRole, query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusNotFound, "User role not found")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, userRole)
}

func DeleteUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	// Parse user_id from request
	userID, err := uuid.Parse(r.PathValue("user_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	// Parse role_id from request
	roleID, err := strconv.Atoi(r.PathValue("role_id"))
	if err != nil {
		utils.HandelError(w, http.StatusBadRequest, "Invalid role_id format")
		return
	}

	query, args, err := QB.Delete("user_roles").
		Where(squirrel.Eq{"user_id": userID, "role_id": roleID}).
		ToSql()

	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error building delete query")
		return
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		utils.HandelError(w, http.StatusInternalServerError, "Error deleting user role")
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "User role deleted successfully")
}
func UpdateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
    // Parse user_id from request
    userID, err := uuid.Parse(r.FormValue("user_id"))
    if err != nil {
        utils.HandelError(w, http.StatusBadRequest, "Invalid user_id format")
        return
    }

    // Parse role_id from request
    roleID, err := strconv.Atoi(r.FormValue("role_id"))
    if err != nil {
        utils.HandelError(w, http.StatusBadRequest, "Invalid role_id format")
        return
    }

    // Create a UserRole instance with the updated values
    userRole := models.UserRole{
        UserID: userID,
        RoleID: roleID,
    }

    // Prepare the update query
    query, args, err := QB.Update("user_roles").
        Set("role_id", userRole.RoleID).
        Where(squirrel.Eq{"user_id": userRole.UserID}).
        ToSql()

    if err != nil {
        utils.HandelError(w, http.StatusInternalServerError, "Error building update query")
        return
    }

    // Execute the update query
    _, err = db.Exec(query, args...)
    if err != nil {
        utils.HandelError(w, http.StatusInternalServerError, "Error updating user role")
        return
    }

    // Respond with a success message
    utils.SendJSONResponse(w, http.StatusOK, map[string]string{"message": "User role updated successfully"})
}

