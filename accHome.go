package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type postAccRequestBody struct {
	Username string `json:"username"`
}

func postAccHome(c *gin.Context) {
	userid, err := validateToken(c.GetHeader("Authorization"), false, true)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "INVALID TOKEN")
		return
	}

	var requestBody postAccRequestBody
	if err = c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "NEED USERNAME IN BODY")
		return
	}

	result, err := mongoClient.Database("npuzzle").Collection("Akun").
		UpdateOne(context.TODO(),
			bson.D{{Key: "creator_id", Value: userid}},
			bson.D{{Key: "$set", Value: bson.D{{Key: "username", Value: requestBody.Username}}}},
		)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	if result.MatchedCount != 1 {
		c.IndentedJSON(http.StatusBadRequest, "USER NOT FOUND")
		return
	}

	c.IndentedJSON(http.StatusOK, nil)
}

func getAccHome(c *gin.Context) {
	userid, err := validateToken(c.GetHeader("Authorization"), false, true)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	var pipeline mongo.Pipeline
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{
		{Key: "creator_id", Value: userid},
	}}})

	pipeline = append(pipeline, bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "Favourites"},
		{Key: "localField", Value: "creator_id"},
		{Key: "foreignField", Value: "creator_id"},
		{Key: "pipeline", Value: mongo.Pipeline{
			bson.D{{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "c", Value: bson.D{{Key: "$size", Value: "$puzzles"}}},
			}}},
		}},
		{Key: "as", Value: "countFav"},
	}}})

	pipeline = append(pipeline, bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "UserHistory"},
		{Key: "localField", Value: "creator_id"},
		{Key: "foreignField", Value: "creator_id"},
		{Key: "pipeline", Value: mongo.Pipeline{
			bson.D{{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "puzzles", Value: 1},
			}}},
		}},
		{Key: "as", Value: "comp"},
	}}})

	pipeline = append(pipeline, bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "Puzzle"},
		{Key: "localField", Value: "creator_id"},
		{Key: "foreignField", Value: "creator_id"},
		{Key: "pipeline", Value: mongo.Pipeline{
			bson.D{{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "puzzle_id", Value: 1},
			}}},
		}},
		{Key: "as", Value: "countMyP"},
	}}})

	pipeline = append(pipeline, bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "SavedState"},
		{Key: "localField", Value: "creator_id"},
		{Key: "foreignField", Value: "creator_id"},
		{Key: "pipeline", Value: mongo.Pipeline{
			bson.D{{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "States", Value: bson.D{{Key: "$map", Value: bson.D{
					{Key: "input", Value: bson.D{{Key: "$objectToArray", Value: "$States"}}},
					{Key: "as", Value: "e"},
					{Key: "in", Value: "$$e.k"},
				}}}},
			}}},
		}},
		{Key: "as", Value: "sStates"},
	}}})

	pipeline = append(pipeline, bson.D{{Key: "$set", Value: bson.D{
		{Key: "countMyP", Value: bson.D{{Key: "$size", Value: "$countMyP"}}},
		{Key: "countFav", Value: bson.D{{Key: "$cond", Value: bson.D{
			{Key: "if", Value: bson.D{{Key: "$gt", Value: bson.A{
				bson.D{{Key: "$size", Value: "$countFav"}},
				0,
			}}}},
			{Key: "then", Value: bson.D{{Key: "$first", Value: "$countFav"}}},
			{Key: "else", Value: 0},
		}}}},
		{Key: "comp", Value: bson.D{{Key: "$cond", Value: bson.D{
			{Key: "if", Value: bson.D{{Key: "$gt", Value: bson.A{
				bson.D{{Key: "$size", Value: "$comp"}},
				0,
			}}}},
			{Key: "then", Value: bson.D{{Key: "$first", Value: "$comp"}}},
			{Key: "else", Value: bson.D{{Key: "puzzles", Value: bson.A{}}}},
		}}}},
		{Key: "sStates", Value: bson.D{{Key: "$cond", Value: bson.D{
			{Key: "if", Value: bson.D{{Key: "$gt", Value: bson.A{
				bson.D{{Key: "$size", Value: "$sStates"}},
				0,
			}}}},
			{Key: "then", Value: bson.D{{Key: "$first", Value: "$sStates"}}},
			{Key: "else", Value: bson.D{{Key: "States", Value: bson.A{}}}},
		}}}},
	}}})

	pipeline = append(pipeline, bson.D{{Key: "$set", Value: bson.D{
		{Key: "countComp", Value: bson.D{{Key: "$size", Value: "$comp.puzzles"}}},
		{Key: "countFav", Value: "$countFav.c"},
		{Key: "comp", Value: "$comp.puzzles"},
		{Key: "sStates", Value: "$sStates.States"},
	}}})

	pipeline = append(pipeline, bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "countUncomp", Value: bson.D{{Key: "$subtract", Value: bson.A{
			bson.D{{Key: "$size", Value: "$sStates"}},
			bson.D{{Key: "$size", Value: bson.D{{Key: "$setIntersection", Value: bson.A{"$sStates", "$comp"}}}}},
		}}}},
	}}})

	pipeline = append(pipeline, bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "ftkn", Value: 0},
		{Key: "dtkn", Value: 0},
		{Key: "comp", Value: 0},
		{Key: "sStates", Value: 0},
		{Key: "creator_id", Value: 0},
	}}})

	cursor, err := mongoClient.Database("npuzzle").Collection("Akun").
		Aggregate(context.TODO(), pipeline, options.Aggregate())

	if err != nil {
		log.Panic(err)
		return
	}

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	var arr []map[string]interface{}
	for _, e := range results {
		arr = append(arr, e.Map())
	}

	if len(arr) > 0 {
		c.IndentedJSON(http.StatusOK, gin.H{
			"username":         arr[0]["username"],
			"countFavorites":   arr[0]["countFav"],
			"countCompleted":   arr[0]["countComp"],
			"countUncompleted": arr[0]["countUncomp"],
			"countMyPuzzles":   arr[0]["countMyP"],
		})
	} else {
		c.IndentedJSON(http.StatusBadRequest, "USER ID NOT FOUND")
	}
}

//curl -X POST --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmUiOiJTdW4sIDI4IEF1ZyAyMDIyIDAzOjE1OjMyICswOTAwIiwiaWQiOiJhbXNsa2RpdWV3cmE1NDUifQ.PuNqwGfcY0h9NooicnlJRCqXOrbo28LcYdFbz_BBngs" --data "{ \"username\": \"testchange\" }" http://localhost:8080
//curl -X GET --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmUiOiJTdW4sIDI4IEF1ZyAyMDIyIDAzOjE1OjMyICswOTAwIiwiaWQiOiJhbXNsa2RpdWV3cmE1NDUifQ.PuNqwGfcY0h9NooicnlJRCqXOrbo28LcYdFbz_BBngs" http://localhost:8080
