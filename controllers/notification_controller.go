package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	middlewares "sodality/handlers"
	"sodality/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var NotificationSetting = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)

	var body models.NotificationSetting

	var newSetting models.NotificationSetting
	err := json.NewDecoder(r.Body).Decode(&newSetting)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	collection := client.Database("sodality").Collection("notificationSetting")
	err = collection.FindOne(r.Context(), bson.D{primitive.E{Key: "user_id", Value: props["user_id"].(string)}}).Decode(&body)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			newSetting.UserID = props["user_id"].(string)
			newSetting.CreatedAt = time.Now().UTC()
			newSetting.UpdatedAt = time.Now().UTC()
			_, err := collection.InsertOne(context.TODO(), newSetting)
			if err != nil {
				middlewares.ServerErrResponse(err.Error(), rw)
				return
			}
			middlewares.SuccessResponse("notification setting updated successfully", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if len(newSetting.Email) <= 0 {
		newSetting.Email = body.Email
	}
	newSetting.UpdatedAt = time.Now().UTC()
	_, err = collection.UpdateOne(r.Context(), bson.D{primitive.E{Key: "_id", Value: body.ID}}, bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "email", Value: newSetting.Email},
				primitive.E{Key: "new_supporters_alerts", Value: newSetting.NewSupporterAlerts},
				primitive.E{Key: "weekly_tips", Value: newSetting.WeeklyTips},
				primitive.E{Key: "weekly_supporter_summary", Value: newSetting.WeeklySupporterSummary},
				primitive.E{Key: "new_crypto_support", Value: newSetting.NewCryptoSupport},
				primitive.E{Key: "updated_at", Value: newSetting.UpdatedAt},
			},
		},
	})

	if err != nil {
		middlewares.ErrorResponse(err.Error(), rw)
		return
	}
	middlewares.SuccessResponse("notification setting updated successfully", rw)
})

var GetNotificationSetting = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	var resp models.NotificationSetting

	collection := client.Database("sodality").Collection("notificationSetting")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "user_id", Value: props["user_id"].(string)}}).Decode(&resp)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.SuccessArrRespond("notification setting does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	middlewares.SuccessArrRespond(resp, rw)
})
