package handler

import (
	"time"

	"github.com/cidmiranda/go-fiber-ws/database"
	"github.com/cidmiranda/go-fiber-ws/model"
	"github.com/cidmiranda/go-fiber-ws/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// LoginUser Login User by username and password
//
//	@Summary		Login User by username and password
//	@Description	Login User by username and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request			body		model.Login		true	"Request of Login User Object"
//	@Failure		401				{object}	model.Message			"Error: Unauthorized"
//	@Failure		400				{object}	model.Message			"Error: Bad Requestt"
//	@Success		200				{object}	model.Token
//	@Router			/api/auth [post]
func LoginUser(c *fiber.Ctx) error {
	db := database.DB.Db
	userLogin := new(model.Login)
	// Store the body in the user and return error if encountered
	err := c.BodyParser(userLogin)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	user := new(model.User)
	// find single user in the database by username
	db.Find(&user, "username = ?", userLogin.Username)
	if user.ID == uuid.Nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": nil})
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Invalid login credentials. Please try again"})
	}

	// Create token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{"token": token})
}

func Restricted(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current book.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}
	return c.SendString("Accessible")
}

// CreateUser Create a new User
//
//		@Summary		Create a new User
//		@Description	Create a new User
//		@Tags			User
//		@Accept			json
//		@Produce		json
//		@Param			request			body		model.User		true	"Request of User Object"
//		@Failure		401				{object}	model.Message			"Error: Unauthorized"
//		@Failure		500				{object}	model.Message			"Error: Bad Requestt"
//		@Success		200				{object}	model.User
//	 	@Security ApiKeyAuth
//		@Router			/api/user [post]
func CreateUser(c *fiber.Ctx) error {
	db := database.DB.Db
	user := new(model.User)
	// Store the body in the user and return error if encountered
	err := c.BodyParser(user)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}
	err = db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create user", "data": err})
	}
	// Return the created user
	user.Password = ""
	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "User has created", "data": user})
}

// GetAllUsers Get All Users from db
//
//		@Summary		Get All Users from db
//		@Description	Get All Users from db
//		@Tags			User
//		@Accept			json
//		@Produce		json
//		@Failure		401				{object}	model.Message				"Error: Unauthorized"
//		@Failure		500				{object}	model.Message				"Error: Bad Requestt"
//		@Success		200				{object}	model.Users
//	 	@Security ApiKeyAuth
//		@Router			/api/user [get]
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DB.Db
	var users []model.User
	// find all users in the database
	db.Find(&users)
	// If no user found, return an error
	if len(users) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Users not found", "data": nil})
	}
	// return users
	return c.Status(200).JSON(fiber.Map{"status": "sucess", "message": "Users Found", "data": users})
}

// GetSingleUser Get a User from db
//
//	@Summary		Get a User from db
//	@Description	Get a User from db
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string			true	"id of User Object"
//	@Failure		401				{object}	model.Message			"Error: Unauthorized"
//	@Failure		500				{object}	model.Message			"Error: Bad Requestt"
//	@Success		200				{object}	model.User
//	@Security ApiKeyAuth
//	@Router			/api/user/{id} [get]
func GetSingleUser(c *fiber.Ctx) error {
	db := database.DB.Db
	// get id params
	id := c.Params("id")
	var user model.User
	// find single user in the database by id
	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User Found", "data": user})
}

// UpdateUser Update a User in db
//
//	@Summary		Update a User in db
//	@Description	Update a User in bd
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			request			body		model.User		true	"Request of User Object"
//	@Param			id				path		string			true	"id of User Object"
//	@Failure		401				{object}	model.Message			"Error: Unauthorized"
//	@Failure		500				{object}	model.Message			"Error: Bad Requestt"
//	@Success		200				{object}	model.User
//	@Security ApiKeyAuth
//	@Router			/api/user{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	type updateUser struct {
		Username string `json:"username"`
	}
	db := database.DB.Db
	var user model.User
	// get id params
	id := c.Params("id")
	// find single user in the database by id
	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}
	var updateUserData updateUser
	err := c.BodyParser(&updateUserData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}
	user.Username = updateUserData.Username
	// Save the Changes
	db.Save(&user)
	// Return the updated user
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "users Found", "data": user})
}

// DeleteUserByID Delete user in db by ID
//
//		@Summary		Delete user in db by ID
//		@Description	Delete user in db by ID
//		@Tags			User
//		@Accept			json
//		@Produce		json
//		@Param			id				path		string			true	"id of User Object"
//		@Failure		401				{object}	model.Message			"Error: Unauthorized"
//		@Failure		500				{object}	model.Message			"Error: Bad Requestt"
//		@Success		200				{object}	model.Message
//	 	@Security ApiKeyAuth
//		@Router			/api/user{id} [delete]
func DeleteUserByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var user model.User
	// get id params
	id := c.Params("id")
	// find single user in the database by id
	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}
	err := db.Delete(&user, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete user", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User deleted"})
}
