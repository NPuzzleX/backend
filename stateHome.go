package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type postStateRequestBody struct {
	Puzzle_id string      `json:"puzzle_id"`
	Data      interface{} `json:"data"`
}

func postStateHome(c *gin.Context) {
	userid, err := validateToken(c.GetHeader("Authorization"), false, true)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "INVALID TOKEN")
		return
	}

	var requestBody postStateRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "NEED PUZZLE ID AND DATA")
		return
	}

	_, err = mongoClient.Database("npuzzle").Collection("SavedState").
		UpdateOne(context.TODO(),
			bson.D{{Key: "creator_id", Value: userid}},
			bson.D{{Key: "$set", Value: bson.D{{Key: "States." + requestBody.Puzzle_id, Value: requestBody.Data}}}},
			options.Update().SetUpsert(true))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, nil)
}

func getStateHome(c *gin.Context) {
	userid, err := validateToken(c.GetHeader("Authorization"), false, true)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	var puzzle_id = ""
	if c.Request.URL.Query().Has("puzzle_id") {
		puzzle_id = c.Request.URL.Query().Get("puzzle_id")
	} else {
		c.IndentedJSON(http.StatusBadRequest, "puzzle_id please")
		return
	}

	var result bson.M
	err = mongoClient.Database("npuzzle").Collection("SavedState").
		FindOne(context.TODO(),
			bson.D{{Key: "creator_id", Value: userid}, {Key: "States." + puzzle_id, Value: bson.D{{Key: "$exists", Value: true}}}},
			options.FindOne().SetProjection(bson.D{{Key: "States." + puzzle_id, Value: 1}, {Key: "_id", Value: 0}}),
		).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.IndentedJSON(http.StatusOK, gin.H{
				"puzzle_id": puzzle_id,
			})
		} else {
			c.IndentedJSON(http.StatusBadRequest, "ERR QUERYING DB")
		}
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"puzzle_id": puzzle_id,
		"data":      result["States"].(bson.M)[puzzle_id],
	})
}

//curl -X POST --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmUiOiJTdW4sIDI4IEF1ZyAyMDIyIDAzOjE1OjMyICswOTAwIiwiaWQiOiJhbXNsa2RpdWV3cmE1NDUifQ.PuNqwGfcY0h9NooicnlJRCqXOrbo28LcYdFbz_BBngs" --data "{ \"username\": \"testchange\" }" http://localhost:8080
//curl -X GET --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmUiOiJTdW4sIDI4IEF1ZyAyMDIyIDAzOjE1OjMyICswOTAwIiwiaWQiOiJhbXNsa2RpdWV3cmE1NDUifQ.PuNqwGfcY0h9NooicnlJRCqXOrbo28LcYdFbz_BBngs" http://localhost:8080
