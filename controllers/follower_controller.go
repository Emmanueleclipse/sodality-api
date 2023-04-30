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
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var FollowCreator = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	params := mux.Vars(r)

	var follow models.Followers

	follow.CreatorID = params["creator_id"]
	follow.UserID = props["user_id"].(string)

	if follow.CreatorID == follow.UserID {
		middlewares.ErrorResponse("you can't follow yourself", rw)
		return
	}

	follow.CreatedAt = time.Now().UTC()

	var existedFollow models.Followers

	filter := bson.M{"$and": []interface{}{
		bson.M{"creator_id": follow.CreatorID},
		bson.M{"user_id": follow.UserID}}}

	collection := client.Database("sodality").Collection("followers")
	err := collection.FindOne(context.TODO(), filter).Decode(&existedFollow)
	if err != nil && err != mongo.ErrNoDocuments {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	if err != mongo.ErrNoDocuments {
		middlewares.ErrorResponse("you already follow this creator", rw)
		return
	}

	result, err := collection.InsertOne(context.TODO(), follow)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessResponse(`inserted at `+strings.Replace(string(res), `"`, ``, 2), rw)
})

var UnfollowCreator = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	params := mux.Vars(r)

	var follow models.Followers

	follow.CreatorID = params["creator_id"]
	follow.UserID = props["user_id"].(string)

	filter := bson.M{"$and": []interface{}{
		bson.M{"creator_id": follow.CreatorID},
		bson.M{"user_id": follow.UserID}}}

	collection := client.Database("sodality").Collection("followers")
	del, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil && err != mongo.ErrNoDocuments {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	if del.DeletedCount == 0 {
		middlewares.ErrorResponse("you already unfollow this creator", rw)
		return
	}

	middlewares.SuccessResponse(`successfully unfollow creator`+follow.CreatorID, rw)
})

var GetCreatorFollowers = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	filter := bson.M{"creator_id": params["user_id"]}

	var followersCount models.FollowersCount
	collection := client.Database("sodality").Collection("followers")
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil && err != mongo.ErrNoDocuments {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	followersCount.Count = count

	middlewares.SuccessRespond(followersCount, rw)
})

var GetCreatorSupporter = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	pipeline := bson.A{
		bson.M{"$match": bson.M{"creator_id": params["creator_id"], "expired_at": bson.M{"$gte": time.Now().UTC()}}},
		bson.M{"$group": bson.M{"_id": "$user_id", "count": bson.M{"$sum": 1}}},
	}

	var supporterCount []bson.M
	collection := client.Database("sodality").Collection("donations")
	cur, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil && err != mongo.ErrNoDocuments {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}
		supporterCount = append(supporterCount, result)
	}
	var totalCount models.SupporterCount

	totalCount.Count = int64(len(supporterCount))

	middlewares.SuccessRespond(totalCount, rw)
})
