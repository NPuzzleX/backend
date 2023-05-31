package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type postFavRequestBody struct {
	Puzzle_id string `json:"puzzle_id"`
}

func postFavHome(c *gin.Context) {
	userid, err := validateToken(c.GetHeader("Authorization"), false, true)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "INVALID TOKEN")
		return
	}

	var requestBody postFavRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "NEED PUZZLE ID AND DATA")
		return
	}

	res, err := mongoClient.Database("npuzzle").Collection("Favourites").
		UpdateOne(context.TODO(),
			bson.D{{Key: "creator_id", Value: userid}, {Key: "puzzles", Value: bson.D{{Key: "$in", Value: bson.A{requestBody.Puzzle_id}}}}},
			bson.D{{Key: "$pull", Value: bson.D{{Key: "puzzles", Value: requestBody.Puzzle_id}}}},
		)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	if res.MatchedCount == 0 {
		_, err = mongoClient.Database("npuzzle").Collection("Favourites").
			UpdateOne(context.TODO(),
				bson.D{{Key: "creator_id", Value: userid}},
				bson.D{{Key: "$addToSet", Value: bson.D{{Key: "puzzles", Value: requestBody.Puzzle_id}}}},
				options.Update().SetUpsert(true),
			)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}
	}

	c.IndentedJSON(http.StatusOK, nil)
}

//curl -X POST --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmUiOiJTdW4sIDI4IEF1ZyAyMDIyIDAzOjE1OjMyICswOTAwIiwiaWQiOiJhbXNsa2RpdWV3cmE1NDUifQ.PuNqwGfcY0h9NooicnlJRCqXOrbo28LcYdFbz_BBngs" --data "{ \"username\": \"testchange\" }" http://localhost:8080
//curl -X GET --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmUiOiJTdW4sIDI4IEF1ZyAyMDIyIDAzOjE1OjMyICswOTAwIiwiaWQiOiJhbXNsa2RpdWV3cmE1NDUifQ.PuNqwGfcY0h9NooicnlJRCqXOrbo28LcYdFbz_BBngs" http://localhost:8080
