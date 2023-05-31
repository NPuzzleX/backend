package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getPuzzle(c *gin.Context) {
	userid, err := validateToken(c.GetHeader("Authorization"), false, false)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "INVALID TOKEN")
		return
	}

	puzzle_id := c.Param("puzzle_id")

	var result bson.M
	err = mongoClient.Database("npuzzle").Collection("Puzzle").
		FindOne(context.TODO(),
			bson.D{{Key: "$and", Value: bson.A{
				bson.D{{Key: "$or", Value: bson.A{
					bson.D{{Key: "hidden", Value: false}},
					bson.D{{Key: "creator_id", Value: userid}},
				}}},
				bson.D{{Key: "puzzle_id", Value: puzzle_id}},
			}}},
			options.FindOne().SetProjection(bson.D{{Key: "data", Value: bson.D{{Key: "$cond", Value: bson.D{
				{Key: "if", Value: bson.D{{Key: "$eq", Value: bson.A{"$data.thumbnail", nil}}}},
				{Key: "then", Value: "$data"},
				{Key: "else", Value: "$data.game"},
			}}}}, {Key: "type", Value: 1}}),
		).Decode(&result)
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, "USER NOT FOUND")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"type": result["type"],
		"data": result["data"],
	})
}

func postAnswer(c *gin.Context) {

	userid, err := validateToken(c.GetHeader("Authorization"), false, true)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "INVALID TOKEN")
		return
	}

	/*
		userid := ""
		var err error
	*/

	puzzle_id := c.Param("puzzle_id")

	var requestBody postPuzzleRequestBody
	if err = c.BindJSON(&requestBody); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest, "NEED TYPE AND DATA")
		return
	}

	if checkAnswer(requestBody.Ptype, requestBody.Data) {
		_, err = mongoClient.Database("npuzzle").Collection("UserHistory").
			UpdateOne(context.TODO(),
				bson.D{{Key: "creator_id", Value: userid}},
				bson.D{{Key: "$addToSet", Value: bson.D{{Key: "puzzles", Value: puzzle_id}}}},
				options.Update().SetUpsert(true),
			)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "ERR PUSHING TO DATABASE")
			return
		}

		c.IndentedJSON(http.StatusOK, true)
	} else {
		c.IndentedJSON(http.StatusOK, false)
	}
}

func deletePuzzle(c *gin.Context) {
	userid, err := validateToken(c.GetHeader("Authorization"), false, true)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "INVALID TOKEN")
		return
	}
	puzzle_id := c.Param("puzzle_id")

	_, err = mongoClient.Database("npuzzle").Collection("Puzzle").
		DeleteOne(context.TODO(),
			bson.D{{Key: "creator_id", Value: userid}, {Key: "puzzle_id", Value: puzzle_id}},
		)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "ERR DELETING IN DATABASE")
		return
	}

	c.IndentedJSON(http.StatusOK, true)

	_, err = mongoClient.Database("npuzzle").Collection("Favourites").
		UpdateMany(context.TODO(),
			bson.D{},
			bson.D{{Key: "$pull", Value: bson.D{{Key: "puzzles", Value: puzzle_id}}}},
		)

	if err != nil {
		log.Println(err)
	}

	_, err = mongoClient.Database("npuzzle").Collection("UserHistory").
		UpdateMany(context.TODO(),
			bson.D{},
			bson.D{{Key: "$pull", Value: bson.D{{Key: "puzzles", Value: puzzle_id}}}},
		)

	if err != nil {
		log.Println(err)
	}

	_, err = mongoClient.Database("npuzzle").Collection("SavedState").
		UpdateMany(context.TODO(),
			bson.D{},
			bson.D{{Key: "$unset", Value: bson.D{{Key: "States." + puzzle_id, Value: ""}}}},
		)

	if err != nil {
		log.Println(err)
	}
}
