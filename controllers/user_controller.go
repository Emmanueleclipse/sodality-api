package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"sodality/db"
	middlewares "sodality/handlers"
	"sodality/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var client = db.Dbconnect()

// RegisterUser -> Register User with email, username and dash
var RegisterUser = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	collection := client.Database("sodality").Collection("users")
	var existingUser models.User
	err = collection.FindOne(r.Context(), bson.D{primitive.E{Key: "username", Value: user.Username}}).Decode(&existingUser)
	if err == nil {
		middlewares.ErrorResponse("username is already taken", rw)
		return
	}
	// err = collection.FindOne(r.Context(), bson.D{primitive.E{Key: "email", Value: user.Email}}).Decode(&existingUser)
	// if err == nil {
	// 	middlewares.ErrorResponse("email is already exists", rw)
	// 	return
	// }
	err = collection.FindOne(r.Context(), bson.D{primitive.E{Key: "dash", Value: user.Dash}}).Decode(&existingUser)
	if err == nil {
		middlewares.ErrorResponse("dash is already exists", rw)
		return
	}
	err = collection.FindOne(r.Context(), bson.D{primitive.E{Key: "mnemonic", Value: user.Mnemonic}}).Decode(&existingUser)
	if err == nil {
		middlewares.ErrorResponse("Mnemonic Invalid", rw)
		return
	}
	passwordHash, err := middlewares.HashPassword(user.Password)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	user.Password = passwordHash
	user.OTPEnabled = false
	result, err := collection.InsertOne(r.Context(), user)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessResponse(`inserted at `+strings.Replace(string(res), `"`, ``, 2), rw)
})

// LoginUser -> Let the user login with identity and password
var LoginUser = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	collection := client.Database("sodality").Collection("users")
	var existingUser models.User
	// if len(user.Username) > 0 {
	err = collection.FindOne(r.Context(), bson.D{primitive.E{Key: "username", Value: user.Username}}).Decode(&existingUser)
	// } else {
	// err = collection.FindOne(r.Context(), bson.D{primitive.E{Key: "email", Value: user.Email}}).Decode(&existingUser)
	// }
	if err != nil {
		middlewares.ErrorResponse("user doesn't exist", rw)
		return
	}
	isPasswordMatch := middlewares.CheckPasswordHash(user.Password, existingUser.Password)
	if !isPasswordMatch {
		middlewares.ErrorResponse("password doesn't match", rw)
		return
	}
	token, err := middlewares.GenerateJWT(existingUser)
	if err != nil {
		middlewares.ErrorResponse("failed to generate token", rw)
		return
	}
	middlewares.SuccessResponse(string(token), rw)
})

// GetUserByID -> Get user details with user id
var GetUserByID = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.User

	userID, _ := primitive.ObjectIDFromHex(params["id"])

	collection := client.Database("sodality").Collection("users")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("user id does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	user.Password = ""
	user.OTPEnabled = false
	user.OTPSecret = ""
	user.OTPAuthURL = ""
	middlewares.SuccessArrRespond(user, rw)
})

// GetProfile -> Get own profile
var GetProfile = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	var user models.User

	userID, _ := primitive.ObjectIDFromHex(props["user_id"].(string))
	collection := client.Database("sodality").Collection("users")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("user id does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	user.Password = ""
	user.OTPSecret = ""
	user.OTPAuthURL = ""

	middlewares.SuccessArrRespond(user, rw)
})

// UpdateUser -> Update user details from username
var UpdateUser = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	var user models.User

	userID, _ := primitive.ObjectIDFromHex(props["user_id"].(string))

	collection := client.Database("sodality").Collection("users")
	err := collection.FindOne(r.Context(), bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&user)
	if err != nil {
		middlewares.AuthorizationResponse("malformed token", rw)
		return
	}

	var newUser models.User
	err = json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	if len(newUser.Username) <= 0 {
		newUser.Username = user.Username
	}
	// if len(newUser.Email) <= 0 {
	// 	newUser.Email = user.Email
	// }
	// if newUser.SubscriberCount <= 0 {
	// 	newUser.SubscriberCount = user.SubscriberCount
	// }
	if len(newUser.HeaderImage) <= 0 {
		newUser.HeaderImage = user.HeaderImage
	}
	if len(newUser.Avatar) <= 0 {
		newUser.Avatar = user.Avatar
	}
	if len(newUser.Dash) <= 0 {
		newUser.Dash = user.Dash
	}
	if len(newUser.Bio) <= 0 {
		newUser.Bio = user.Bio
	}
	if len(newUser.Role) <= 0 {
		newUser.Role = user.Role
	}
	if len(newUser.Description) <= 0 {
		newUser.Description = user.Description
	}
	if len(newUser.Title) <= 0 {
		newUser.Title = user.Title
	}
	if len(newUser.SubTitle) <= 0 {
		newUser.SubTitle = user.SubTitle
	}

	res, err := collection.UpdateOne(r.Context(), bson.D{primitive.E{Key: "_id", Value: user.ID}}, bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "username", Value: newUser.Username},
				// primitive.E{Key: "email", Value: newUser.Email},
				primitive.E{Key: "header_image", Value: newUser.HeaderImage},
				primitive.E{Key: "title", Value: newUser.Title},
				primitive.E{Key: "subtitle", Value: newUser.SubTitle},
				primitive.E{Key: "description", Value: newUser.Description},
				primitive.E{Key: "avatar", Value: newUser.Avatar},
				primitive.E{Key: "dash", Value: newUser.Dash},
				primitive.E{Key: "bio", Value: newUser.Bio},
				primitive.E{Key: "role", Value: newUser.Role},
			},
		},
	})

	if err != nil {
		middlewares.ErrorResponse("username is already taken.", rw)
		return
	}
	if res.MatchedCount == 0 {
		middlewares.ErrorResponse("user doesn't exist", rw)
		return
	}

	token, err := middlewares.GenerateJWT(newUser)
	if err != nil {
		middlewares.ErrorResponse("Failed to generate JWT", rw)
		return
	}
	middlewares.SuccessResponse(string(token), rw)
})
