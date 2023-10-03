package controllers

import (
	"context"
	"echo-app/configs"
	"echo-app/models"
	"echo-app/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func CreateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user models.User
	defer cancel()

	//Validate the request body
	if err := c.Bind(&user); err != nil {
		return Error(c, err, http.StatusBadRequest)
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return Error(c, validationErr, http.StatusBadRequest)
	}

	newUser := models.User{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}

	_, err := userCollection.InsertOne(ctx, newUser)

	if err != nil {
		return Error(c, err, http.StatusInternalServerError)
	}

	return Success(c, newUser)
}

func GetUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

	if err != nil {
		return Error(c, err, http.StatusInternalServerError)
	}

	return Success(c, user)
}

func EditUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	//validate the request body
	if err := c.Bind(&user); err != nil {
		return Error(c, err, http.StatusBadRequest)
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return Error(c, validationErr, http.StatusBadRequest)
	}

	update := bson.M{
		"name":     user.Name,
		"location": user.Location,
		"title":    user.Title,
	}

	result, err := userCollection.UpdateOne(
		ctx,
		bson.M{"id": objId},
		bson.M{"$set": update},
	)

	if err != nil {
		return Error(c, err, http.StatusBadRequest)
	}

	//get updated user details
	var updatedUser models.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(
			ctx,
			bson.M{"id": objId}).Decode(&updatedUser)

		if err != nil {
			return Error(c, err, http.StatusInternalServerError)
		}
	}

	return Success(c, updatedUser)
}

func DeleteUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})

	if err != nil {
		return Error(c, err, http.StatusInternalServerError)
	}

	if result.DeletedCount < 1 {
		return Error(c, err, http.StatusNotFound)
	}

	return Success(c, models.User{Id: objId})
}

func GetAllUsers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})

	if err != nil {
		return Error(c, err, http.StatusNotFound)
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return Error(c, err, http.StatusInternalServerError)
		}

		users = append(users, singleUser)
	}

	return c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": users}})
}

func Error(c echo.Context, validationErr error, httpStatus int) error {
	return c.JSON(
		httpStatus,
		responses.UserResponse{
			Status:  httpStatus,
			Message: "error",
			Data: &echo.Map{
				"data": validationErr.Error(),
			},
		})
}

func Success(c echo.Context, user models.User) error {
	return c.JSON(
		http.StatusOK,
		responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data: &echo.Map{
				"data": user,
			},
		})
}
