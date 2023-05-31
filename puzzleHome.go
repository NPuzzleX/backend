package main

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	uuid "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type postPuzzleRequestBody struct {
	Data   interface{} `json:"data"`
	Ptype  string      `json:"type"`
	Hidden bool        `json:"hidden"`
}

type putPuzzleRequestBody struct {
	Data      interface{} `json:"data"`
	Ptype     string      `json:"type"`
	Puzzle_id string      `json:"puzzle_id"`
	Hidden    bool        `json:"hidden"`
}

func postPuzzleHome(c *gin.Context) {
	userid, err := validateToken(c.GetHeader("Authorization"), false, true)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "INVALID TOKEN")
		return
	}

	var requestBody postPuzzleRequestBody
	if err = c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "NEED TYPE, DATA, AND HIDDEN")
		return
	}

	gameData := requestBody.Data.(map[string]interface{})["game"]
	if gameData == nil {
		c.IndentedJSON(http.StatusBadRequest, "NEED GAME DATA")
		return
	}

	var result bson.M
	err = mongoClient.Database("npuzzle").Collection("Puzzle").
		FindOne(context.TODO(),
			bson.D{{Key: "data.game", Value: gameData}}, options.FindOne().SetProjection(bson.D{{Key: "creator_id", Value: 1}, {Key: "_id", Value: 0}})).
		Decode(&result)

	if err == mongo.ErrNoDocuments {
		_, err = mongoClient.Database("npuzzle").Collection("Puzzle").
			InsertOne(context.TODO(),
				bson.D{
					{Key: "creator_id", Value: userid},
					{Key: "puzzle_id", Value: uuid.NewString()},
					{Key: "type", Value: requestBody.Ptype},
					{Key: "data", Value: requestBody.Data},
					{Key: "hidden", Value: requestBody.Hidden},
				})

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}
		c.IndentedJSON(http.StatusOK, nil)
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "ERROR DATABASE CONNECTION")
		return
	} else {
		c.IndentedJSON(http.StatusOK, "PUZZLE ALREADY EXIST")
		return
	}
}

func putPuzzleHome(c *gin.Context) {
	userid, err := validateToken(c.GetHeader("Authorization"), false, true)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "INVALID TOKEN")
		return
	}

	var requestBody putPuzzleRequestBody
	if err = c.BindJSON(&requestBody); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "NEED TYPE, DATA, HIDDNE, AND PUZZLE_ID")
		return
	}

	gameData := requestBody.Data.(map[string]interface{})["game"]
	if gameData == nil {
		c.IndentedJSON(http.StatusBadRequest, "NEED GAME DATA")
		return
	}

	var result bson.M
	err = mongoClient.Database("npuzzle").Collection("Puzzle").
		FindOne(context.TODO(),
			bson.D{{Key: "data.game", Value: gameData}}, options.FindOne().SetProjection(bson.D{{Key: "creator_id", Value: 1}, {Key: "_id", Value: 0}})).
		Decode(&result)

	if err == mongo.ErrNoDocuments {
		_, err = mongoClient.Database("npuzzle").Collection("Puzzle").
			UpdateOne(context.TODO(),
				bson.D{{Key: "creator_id", Value: userid}, {Key: "puzzle_id", Value: requestBody.Puzzle_id}},
				bson.D{{Key: "$set", Value: bson.D{{Key: "type", Value: requestBody.Ptype}, {Key: "data", Value: requestBody.Data}, {Key: "hidden", Value: requestBody.Hidden}}}})

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err)
			return
		}

		c.IndentedJSON(http.StatusOK, nil)
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "ERROR DATABASE CONNECTION")
		return
	} else {
		c.IndentedJSON(http.StatusOK, "PUZZLE ALREADY EXIST")
		return
	}
}

func getPuzzleHome(c *gin.Context) {
	userid, err := validateToken(c.GetHeader("Authorization"), false, false)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	var page = 1
	if c.Request.URL.Query().Has("page") {
		page, err = strconv.Atoi(c.Request.URL.Query().Get("page"))
		if err != nil {
			page = 1
		}
	}
	if page < 1 {
		page = 1
	}

	var limit = 5
	if c.Request.URL.Query().Has("limit") {
		limit, err = strconv.Atoi(c.Request.URL.Query().Get("limit"))
		if err != nil {
			limit = 5
		}
	}
	if limit < 1 {
		limit = 1
	}

	var creator_id = ""
	if c.Request.URL.Query().Get("self") == "1" {
		creator_id = userid
	} else if c.Request.URL.Query().Has("creator_id") {
		creator_id = c.Request.URL.Query().Get("creator_id")
	}

	var completed = false
	if c.Request.URL.Query().Has("completed") {
		if c.Request.URL.Query().Get("completed") == "true" {
			completed = true
		} else {
			completed = false
		}
	}

	var favourited = false
	if c.Request.URL.Query().Has("favourited") {
		if c.Request.URL.Query().Get("favourited") == "true" {
			favourited = true
		} else {
			favourited = false
		}
	}

	var pipeline mongo.Pipeline

	var completedPipeline mongo.Pipeline = mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "UserHistory"},
			{Key: "let", Value: bson.D{{Key: "pid", Value: "$puzzle_id"}}},
			{Key: "pipeline", Value: mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$and", Value: bson.A{
						bson.D{{Key: "$expr", Value: bson.D{
							{Key: "$in", Value: bson.A{
								"$$pid", "$puzzles",
							}},
						}}},
						bson.D{{Key: "creator_id", Value: userid}},
					}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "creator_id", Value: 0},
					{Key: "puzzles", Value: 0},
				}}},
			}},
			{Key: "as", Value: "completed"},
		}}},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "completed", Value: bson.D{{Key: "$cond", Value: bson.D{
				{Key: "if", Value: bson.D{{Key: "$gt", Value: bson.A{
					bson.D{{Key: "$size", Value: "$completed"}},
					0,
				}}}},
				{Key: "then", Value: true},
				{Key: "else", Value: false},
			}}}},
		}}},
	}

	var favouritedPipeline mongo.Pipeline = mongo.Pipeline{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "Favourites"},
			{Key: "let", Value: bson.D{{Key: "pid", Value: "$puzzle_id"}}},
			{Key: "pipeline", Value: mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$and", Value: bson.A{
						bson.D{{Key: "$expr", Value: bson.D{
							{Key: "$in", Value: bson.A{
								"$$pid", "$puzzles",
							}},
						}}},
						bson.D{{Key: "creator_id", Value: userid}},
					}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "creator_id", Value: 0},
					{Key: "puzzles", Value: 0},
				}}},
			}},
			{Key: "as", Value: "favourited"},
		}}},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "favourited", Value: bson.D{{Key: "$cond", Value: bson.D{
				{Key: "if", Value: bson.D{{Key: "$gt", Value: bson.A{
					bson.D{{Key: "$size", Value: "$favourited"}},
					0,
				}}}},
				{Key: "then", Value: true},
				{Key: "else", Value: false},
			}}}},
		}}},
	}

	filter := bson.D{
		{Key: "creator_id", Value: 1},
		{Key: "type", Value: 1},
		{Key: "puzzle_id", Value: 1},
		{Key: "thumbnail", Value: bson.D{{Key: "$cond", Value: bson.D{
			{Key: "if", Value: bson.D{{Key: "$eq", Value: bson.A{"$data.thumbnail", nil}}}},
			{Key: "then", Value: nil},
			{Key: "else", Value: "$data.thumbnail"},
		}}}},
	}
	if userid == creator_id {
		filter = append(filter, bson.E{Key: "hidden", Value: 1})
	} else {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "hidden", Value: false}}}})
	}
	pipeline = append(pipeline, bson.D{{Key: "$project", Value: filter}})

	filter = bson.D{}
	if creator_id != "" {
		filter = append(filter, bson.E{Key: "creator_id", Value: creator_id})
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: filter}})
	}

	if completed && (userid != "") {
		pipeline = append(pipeline, completedPipeline...)
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "completed", Value: true}}}})
	}

	if favourited && (userid != "") {
		pipeline = append(pipeline, favouritedPipeline...)
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "favourited", Value: true}}}})
	}

	//------------------------------- TOTAL PAGE QUERY -------------------------------
	cursor, err := mongoClient.Database("npuzzle").Collection("Puzzle").
		Aggregate(context.TODO(), append(pipeline, bson.D{{Key: "$count", Value: "count"}}), options.Aggregate())

	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, "ERR QUERYING DB")
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

	var pageCount int32
	if len(arr) > 0 {
		pageCount = (arr[0]["count"]).(int32)
	}

	arr = nil
	//=============================== TOTAL PAGE QUERY ===============================

	pipeline = append(pipeline, bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: -1}}}})

	if page != 1 {
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: int64((page - 1) * limit)}})
	}

	pipeline = append(pipeline, bson.D{{Key: "$limit", Value: int64(limit)}})

	if !completed {
		if userid != "" {
			pipeline = append(pipeline, completedPipeline...)
		} else {
			pipeline = append(pipeline, bson.D{{Key: "$set", Value: bson.D{{Key: "completed", Value: false}}}})
		}
	}

	if !favourited {
		if userid != "" {
			pipeline = append(pipeline, favouritedPipeline...)
		} else {
			pipeline = append(pipeline, bson.D{{Key: "$set", Value: bson.D{{Key: "favourited", Value: false}}}})
		}
	}

	pipeline = append(pipeline, bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "Akun"},
		{Key: "localField", Value: "creator_id"},
		{Key: "foreignField", Value: "creator_id"},
		{Key: "as", Value: "d"},
	}}})

	pipeline = append(pipeline, bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$d"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}})

	pipeline = append(pipeline, bson.D{{Key: "$addFields", Value: bson.D{{Key: "username", Value: "$d.username"}}}})

	pipeline = append(pipeline, bson.D{{Key: "$project", Value: bson.D{
		{Key: "d", Value: 0},
	}}})

	cursor, err = mongoClient.Database("npuzzle").Collection("Puzzle").
		Aggregate(context.TODO(), pipeline, options.Aggregate())

	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, "ERR QUERYING DB")
		return
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for _, e := range results {
		arr = append(arr, e.Map())
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"data":  arr,
		"pages": math.Ceil(float64(pageCount) / float64(limit)),
	})
}

//curl -X POST --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmUiOiJTdW4sIDI4IEF1ZyAyMDIyIDAzOjE1OjMyICswOTAwIiwiaWQiOiJhbXNsa2RpdWV3cmE1NDUifQ.PuNqwGfcY0h9NooicnlJRCqXOrbo28LcYdFbz_BBngs" --data "{ \"username\": \"testchange\" }" http://localhost:8080
//curl -X GET --header "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmUiOiJNb24sIDIxIE5vdiAyMDIyIDE0OjE0OjQ5ICswMDAwIiwiaWQiOiI3YTlhN2FjMy00MDFlLTQ0YWItOTA0Yi02NDVmNDEwMGYzOGYifQ.VphM9F0tGk4QVf-sU3JKjNZgt5u4PLf8DWvnzoq-smg" "http://localhost:8080/?page=1&limit=1"
