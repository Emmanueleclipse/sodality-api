package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	middlewares "sodality/handlers"
	"sodality/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var AddCreatorTiers = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)

	username, _ := props["username"].(string)

	var tier models.CreatorTiers
	err := json.NewDecoder(r.Body).Decode(&tier)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	var existedTier models.CreatorTiers

	tier.Username = username

	collection := client.Database("sodality").Collection("creatorTiers")
	err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "username", Value: username}}).Decode(&existedTier)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err := collection.InsertOne(r.Context(), tier)
			if err != nil {
				middlewares.ServerErrResponse(err.Error(), rw)
				return
			}
			middlewares.SuccessResponse("success", rw)
			return
		} else {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}
	}

	_, err = collection.UpdateOne(r.Context(), bson.D{primitive.E{Key: "username", Value: username}}, bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "tier_one_name", Value: tier.TierOneName},
				primitive.E{Key: "tier_one_price", Value: tier.TierOnePrice},
				primitive.E{Key: "tier_two_name", Value: tier.TierTwoName},
				primitive.E{Key: "tier_two_price", Value: tier.TierTwoPrice},
				primitive.E{Key: "tier_three_name", Value: tier.TierThreeName},
				primitive.E{Key: "tier_three_price", Value: tier.TierThreePrice},
			},
		},
	})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	middlewares.SuccessRespond("success", rw)
})

var GetCreatorTierByUserID = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var tier models.CreatorTiers

	collection := client.Database("sodality").Collection("creatorTiers")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "username", Value: params["username"]}}).Decode(&tier)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.SuccessRespond(nil, rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	middlewares.SuccessRespond(tier, rw)
})
