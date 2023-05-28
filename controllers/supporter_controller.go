package controllers

import (
	"context"
	"net/http"
	middlewares "sodality/handlers"
	"sodality/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var GetContentForSpecificSupporterByID = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	params := mux.Vars(r)
	var content models.Content

	contentID, _ := primitive.ObjectIDFromHex(params["id"])

	collection := client.Database("sodality").Collection("content")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: contentID}}).Decode(&content)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("content id does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	filter := bson.M{"$and": []interface{}{
		bson.M{"user_id": props["user_id"].(string)},
		bson.M{"creator_id": content.UserID},
		bson.M{"expired_at": bson.M{"$gte": time.Now().UTC()}},
	}}
	// opts := options.Find().SetSort(bson.D{primitive.E{Key: "fund", Value: -1}})

	var allDonationForCreator []*models.Donate
	donationCollection := client.Database("sodality").Collection("donations")
	cursor, err := donationCollection.Find(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("content does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var donations models.Donate
		err := cursor.Decode(&donations)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		allDonationForCreator = append(allDonationForCreator, &donations)
	}

	var totalDonation float64 = 0

	for _, v := range allDonationForCreator {
		totalDonation += v.Donate
	}

	// if (totalDonation < 5 && totalDonation >= 1) && content.ContentType == "Supporter" {
	// 	content.Locked = false
	// } else if (totalDonation < 10 && totalDonation >= 5) && content.ContentType != "Super Fan" {
	// 	content.Locked = false
	// 	if (totalDonation < 5) && content.ContentType == "Fan" {
	// 		content.Locked = true
	// 	}
	// } else if totalDonation >= 10 {
	// 	content.Locked = false
	// }
	middlewares.SuccessArrRespond(content, rw)
})

var GetAllCreatorsContentForSpecificSupporter = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	var allContent []*models.GetAllContentWithCreatorResp

	opts := options.Find().SetSort(bson.D{primitive.E{Key: "fund", Value: -1}})

	collection := client.Database("sodality").Collection("content")
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("content does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var content models.GetAllContentWithCreatorResp
		err := cursor.Decode(&content)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		allContent = append(allContent, &content)
	}
	for _, v := range allContent {
		var user models.User

		userID, _ := primitive.ObjectIDFromHex(v.UserID)
		collection := client.Database("sodality").Collection("users")
		collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&user)
		user.Password = ""
		user.OTPEnabled = false
		user.OTPSecret = ""
		user.OTPAuthURL = ""
		v.User = user
	}

	for _, cont := range allContent {

		filter := bson.M{"$and": []interface{}{
			bson.M{"user_id": props["user_id"].(string)},
			bson.M{"creator_id": cont.UserID},
			bson.M{"expired_at": bson.M{"$gte": time.Now().UTC()}},
		}}

		var allDonationForCreator []*models.Donate
		donationCollection := client.Database("sodality").Collection("donations")
		cursor, err := donationCollection.Find(context.TODO(), filter)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				middlewares.ErrorResponse("content does not exist", rw)
				return
			}
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}
		for cursor.Next(context.TODO()) {
			var donations models.Donate
			err := cursor.Decode(&donations)
			if err != nil {
				middlewares.ServerErrResponse(err.Error(), rw)
				return
			}

			allDonationForCreator = append(allDonationForCreator, &donations)
		}

		var totalDonation float64 = 0

		for _, v := range allDonationForCreator {
			totalDonation += v.Donate
		}

		if (totalDonation < 5 && totalDonation >= 1) && cont.ContentType == "Supporter" {
			cont.Locked = false
		} else if (totalDonation < 10 && totalDonation >= 5) && cont.ContentType != "Super Fan" {
			cont.Locked = false
			if (totalDonation < 5) && cont.ContentType == "Fan" {
				cont.Locked = true
			}
		} else if totalDonation >= 10 {
			cont.Locked = false
		}
	}

	middlewares.SuccessArrRespond(allContent, rw)
})

var GetCreatorContentsForSpecificSupporter = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	var creatorContents []*models.Content

	var user models.User

	userCollection := client.Database("sodality").Collection("users")
	err := userCollection.FindOne(context.TODO(), bson.D{primitive.E{Key: "username", Value: params["username"]}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("user does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	collection := client.Database("sodality").Collection("content")
	cursor, err := collection.Find(context.TODO(), bson.D{primitive.E{Key: "user_id", Value: user.ID.Hex()}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("content id does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var content models.Content
		err := cursor.Decode(&content)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		creatorContents = append(creatorContents, &content)
	}

	filter := bson.M{"$and": []interface{}{
		bson.M{"username": props["username"].(string)},
		bson.M{"creator_username": params["username"]},
		bson.M{"expired_at": bson.M{"$gte": time.Now().UTC()}},
	}}

	var allDonationForCreator []*models.Donate
	donationCollection := client.Database("sodality").Collection("donations")
	cursor1, err := donationCollection.Find(context.TODO(), filter)
	if err != nil && err != mongo.ErrNoDocuments {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor1.Next(context.TODO()) {
		var donations models.Donate
		err := cursor1.Decode(&donations)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		allDonationForCreator = append(allDonationForCreator, &donations)
	}

	for _, cont := range creatorContents {
		cont.Locked = true

		if allDonationForCreator == nil {
			cont.Locked = true
			continue
		}

		for _, donate := range allDonationForCreator {
			if donate.BuyTier == TierOne && cont.TierType == TierOne {
				cont.Locked = false
			} else if donate.BuyTier == TierTwo {
				if cont.TierType == TierOne || cont.TierType == TierTwo {
					cont.Locked = false
				}
			} else if donate.BuyTier == TierThree {
				if cont.TierType == TierOne || cont.TierType == TierTwo || cont.TierType == TierThree {
					cont.Locked = false
				}
			}

			if cont.TierType == AllTier {
				if cont.TierType == TierOne || cont.TierType == TierTwo || cont.TierType == TierThree || cont.TierType == AllTier {
					cont.Locked = false
				}
			}
		}
	}

	middlewares.SuccessArrRespond(creatorContents, rw)
})

var GetCreatorDirectoryByDirectoryNameForSpecificSupporter = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	params := mux.Vars(r)
	var creatorContents []*models.Content

	opts := options.Find().SetSort(bson.D{primitive.E{Key: "fund", Value: -1}})

	collection := client.Database("sodality").Collection("content")
	cursor, err := collection.Find(context.TODO(), bson.D{primitive.E{Key: "category_name", Value: params["category_name"]}}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("content does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var content models.Content
		err := cursor.Decode(&content)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		creatorContents = append(creatorContents, &content)
	}

	filter := bson.M{"$and": []interface{}{
		bson.M{"username": props["username"].(string)},
		bson.M{"creator_username": params["username"]},
		bson.M{"expired_at": bson.M{"$gte": time.Now().UTC()}},
	}}

	var allDonationForCreator []*models.Donate
	donationCollection := client.Database("sodality").Collection("donations")
	cursor1, err := donationCollection.Find(context.TODO(), filter)
	if err != nil && err != mongo.ErrNoDocuments {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor1.Next(context.TODO()) {
		var donations models.Donate
		err := cursor1.Decode(&donations)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		allDonationForCreator = append(allDonationForCreator, &donations)
	}

	for _, cont := range creatorContents {
		cont.Locked = true

		if allDonationForCreator == nil {
			cont.Locked = true
			continue
		}

		for _, donate := range allDonationForCreator {
			if donate.BuyTier == TierOne && cont.TierType == TierOne {
				cont.Locked = false
			} else if donate.BuyTier == TierTwo {
				if cont.TierType == TierOne || cont.TierType == TierTwo {
					cont.Locked = false
				}
			} else if donate.BuyTier == TierThree {
				if cont.TierType == TierOne || cont.TierType == TierTwo || cont.TierType == TierThree {
					cont.Locked = false
				}
			}

			if cont.TierType == AllTier {
				if cont.TierType == TierOne || cont.TierType == TierTwo || cont.TierType == TierThree || cont.TierType == AllTier {
					cont.Locked = false
				}
			}
		}
	}

	middlewares.SuccessArrRespond(creatorContents, rw)
})

var GetCreatorSupportersRecord = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var supporterRecord []*models.DonateResp

	var user models.User

	userCollection := client.Database("sodality").Collection("users")
	err := userCollection.FindOne(context.TODO(), bson.D{primitive.E{Key: "username", Value: params["username"]}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("user does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	collection := client.Database("sodality").Collection("donations")

	cursor, err := collection.Find(context.TODO(), bson.D{primitive.E{Key: "creator_id", Value: user.ID.Hex()}}, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("content does not exist", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var record models.DonateResp
		err := cursor.Decode(&record)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		supporterRecord = append(supporterRecord, &record)
	}

	for _, v := range supporterRecord {
		var user models.User

		userID, _ := primitive.ObjectIDFromHex(v.UserID)
		collection := client.Database("sodality").Collection("users")
		collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&user)
		user.Password = ""
		user.OTPEnabled = false
		user.OTPSecret = ""
		user.OTPAuthURL = ""
		v.User = user
	}

	middlewares.SuccessRespond(supporterRecord, rw)
})

var GetRecentSubscription = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	props, _ := r.Context().Value("props").(jwt.MapClaims)

	var supporterRecord []*models.LastDonationResp

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	collection := client.Database("sodality").Collection("donations")

	filter := bson.M{"$and": []interface{}{
		bson.M{"username": props["username"].(string)},
		bson.M{"creator_username": params["username"]},
	}}

	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.SuccessRespond(nil, rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var record models.LastDonationResp
		err := cursor.Decode(&record)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		supporterRecord = append(supporterRecord, &record)
	}
	var totalDonation float64 = 0

	for _, v := range supporterRecord {
		totalDonation += v.Donate
	}
	if len(supporterRecord) <= 0 {
		middlewares.SuccessRespond(nil, rw)
		return
	}

	supporterRecord[0].AverageDonation = totalDonation / float64(len(supporterRecord))

	middlewares.SuccessRespond(supporterRecord[0], rw)
})
