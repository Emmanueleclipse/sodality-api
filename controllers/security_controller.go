package controllers

import (
	"encoding/json"
	"net/http"
	middlewares "sodality/handlers"
	"sodality/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var GenerateQR = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "sodality.com",
		AccountName: props["username"].(string),
		SecretSize:  15,
	})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	id, _ := props["user_id"].(string)
	userID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"otp_secret":   key.Secret(),
			"otp_auth_url": key.URL(),
		},
	}

	userCollection := client.Database("sodality").Collection("users")
	_, err = userCollection.UpdateOne(r.Context(), filter, update)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	var resp models.GenerateAuthURL
	resp.OTPSecret = key.Secret()
	resp.URL = key.URL()
	middlewares.SuccessRespond(resp, rw)
})

var VerifyOTP = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	id, _ := props["user_id"].(string)
	userID, _ := primitive.ObjectIDFromHex(id)

	collection := client.Database("sodality").Collection("users")
	var existingUser models.User
	err = collection.FindOne(r.Context(), bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&existingUser)
	if err != nil {
		middlewares.ErrorResponse("user doesn't exist", rw)
		return
	}

	valid := totp.Validate(user.Token, existingUser.OTPSecret)
	if !valid {
		middlewares.ErrorResponse("invalid otp", rw)
		return
	}

	middlewares.SuccessRespond("otp verified successfully", rw)
})

var Update2FA = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	userID, _ := primitive.ObjectIDFromHex(props["user_id"].(string))

	collection := client.Database("sodality").Collection("users")
	_, err = collection.UpdateOne(r.Context(), bson.D{primitive.E{Key: "_id", Value: userID}}, bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "otp_enabled", Value: user.OTPEnabled},
			},
		},
	})

	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	middlewares.SuccessResponse("2FA update successfully", rw)
})
