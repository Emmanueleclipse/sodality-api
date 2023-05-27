package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	middlewares "sodality/handlers"
	"sodality/models"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var DonateUser = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)

	var donate models.Donate
	err := json.NewDecoder(r.Body).Decode(&donate)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	username, _ := props["username"].(string)
	if username == donate.CreatorUsername {
		middlewares.ServerErrResponse("you cannot donate yourself", rw)
		return
	}

	var existingCreator models.User
	// creatorID, _ := primitive.ObjectIDFromHex(donate.CreatorID)

	userCollection := client.Database("sodality").Collection("users")
	err = userCollection.FindOne(r.Context(), bson.D{primitive.E{Key: "username", Value: donate.CreatorUsername}, primitive.E{Key: "role", Value: "creator"}}).Decode(&existingCreator)
	if err != nil && err != mongo.ErrNoDocuments {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if err == mongo.ErrNoDocuments {
		middlewares.ErrorResponse("creator does not exists", rw)
		return
	}

	donate.Username = username
	donate.CreatedAt = time.Now().UTC()
	donate.ExpiredAt = time.Now().UTC().AddDate(0, 1, 0)

	donationCollection := client.Database("sodality").Collection("donations")
	result, err := donationCollection.InsertOne(context.TODO(), donate)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	count := existingCreator.TotalDonations
	count += donate.Donate

	_, err = userCollection.UpdateOne(r.Context(), bson.D{primitive.E{Key: "username", Value: donate.CreatorUsername}}, bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "total_donations", Value: count},
			},
		},
	})
	if err != nil {
		middlewares.ErrorResponse(err.Error(), rw)
		return
	}

	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessResponse(`inserted at `+strings.Replace(string(res), `"`, ``, 2), rw)
})

var DonateContent = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)

	var donate models.DonateContent
	err := json.NewDecoder(r.Body).Decode(&donate)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	var existingContent models.Content
	contentID, _ := primitive.ObjectIDFromHex(donate.ContentID)

	contentCollection := client.Database("sodality").Collection("content")
	err = contentCollection.FindOne(r.Context(), bson.D{primitive.E{Key: "_id", Value: contentID}}).Decode(&existingContent)
	if err != nil && err != mongo.ErrNoDocuments {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if existingContent.ID != contentID || err == mongo.ErrNoDocuments {
		middlewares.ErrorResponse("content does not exists", rw)
		return
	}

	userID, _ := props["user_id"].(string)
	donate.UserID = userID
	donate.CreatedAt = time.Now().UTC()
	donate.ExpiredAt = time.Now().UTC().AddDate(0, 1, 0)

	if existingContent.UserID == userID {
		middlewares.ServerErrResponse("you cannot donate your own content", rw)
		return

	}

	donationCollection := client.Database("sodality").Collection("contentDonations")
	result, err := donationCollection.InsertOne(context.TODO(), donate)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	count := existingContent.Fund
	count += donate.Donate

	_, err = contentCollection.UpdateOne(r.Context(), bson.D{primitive.E{Key: "_id", Value: contentID}}, bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "fund", Value: count},
			},
		},
	})
	if err != nil {
		middlewares.ErrorResponse(err.Error(), rw)
		return
	}

	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessResponse(`inserted at `+strings.Replace(string(res), `"`, ``, 2), rw)
})
