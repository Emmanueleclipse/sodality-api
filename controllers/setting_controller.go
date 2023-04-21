package controllers

import (
	"encoding/json"
	"net/http"
	middlewares "sodality/handlers"
	"sodality/models"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CreatorSetting = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
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
	if len(newUser.Avatar) <= 0 {
		newUser.Avatar = user.Avatar
	}
	if len(newUser.HeaderImage) <= 0 {
		newUser.HeaderImage = user.HeaderImage
	}
	if len(newUser.Title) <= 0 {
		newUser.Title = user.Title
	}
	if len(newUser.SubTitle) <= 0 {
		newUser.SubTitle = user.SubTitle
	}
	if len(newUser.Description) <= 0 {
		newUser.Description = user.Description
	}
	if len(newUser.Facebook) <= 0 {
		newUser.Facebook = user.Facebook
	}
	if len(newUser.Twitter) <= 0 {
		newUser.Twitter = user.Twitter
	}
	if len(newUser.Odysee) <= 0 {
		newUser.Odysee = user.Odysee
	}

	res, err := collection.UpdateOne(r.Context(), bson.D{primitive.E{Key: "_id", Value: user.ID}}, bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "username", Value: newUser.Username},
				primitive.E{Key: "avatar", Value: newUser.Avatar},
				primitive.E{Key: "header_image", Value: newUser.HeaderImage},
				primitive.E{Key: "title", Value: newUser.Title},
				primitive.E{Key: "subtitle", Value: newUser.SubTitle},
				primitive.E{Key: "description", Value: newUser.Description},
				primitive.E{Key: "facebook", Value: newUser.Facebook},
				primitive.E{Key: "twitter", Value: newUser.Twitter},
				primitive.E{Key: "youtube", Value: newUser.Youtube},
				primitive.E{Key: "odysee", Value: newUser.Odysee},
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
	middlewares.SuccessResponse("setting updated successfully", rw)
})
