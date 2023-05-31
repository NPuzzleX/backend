package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getLogin(c *gin.Context) {
	var tkn string = ""
	var loginType string = ""
	var uid string = ""

	var cid string = ""
	if c.Request.URL.Query().Has("ftkn") {
		tkn = c.Request.URL.Query().Get("ftkn")
		loginType = "ftkn"

		token, err := authClient.VerifyIDToken(context.TODO(), tkn)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "ERROR VERIVYING TOKEN")
			return
		}
		uid = token.UID

		h := sha256.New()
		h.Write([]byte(uid))
		var hashedUid = base64.URLEncoding.EncodeToString(h.Sum(nil))

		var result bson.M
		err = mongoClient.Database("npuzzle").Collection("Akun").
			FindOne(context.TODO(),
				bson.D{{Key: loginType, Value: hashedUid}}, options.FindOne().SetProjection(bson.D{{Key: "creator_id", Value: 1}, {Key: "_id", Value: 0}})).
			Decode(&result)

		if err == mongo.ErrNoDocuments {
			var userName = ""
			if loginType == "ftkn" {
				userdt, err := authClient.GetUser(context.TODO(), uid)
				if err != nil {
					c.IndentedJSON(http.StatusBadRequest, "ERROR RETRIEVING FIREBASE USER DATA")
				}
				userName = userdt.DisplayName
			}

			cid = uuid.New().String()
			res, err := mongoClient.Database("npuzzle").Collection("Akun").InsertOne(context.TODO(), bson.D{{
				Key: "creator_id", Value: cid}, {Key: loginType, Value: hashedUid}, {Key: "username", Value: userName}})

			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, "ERROR DATABASE CONNECTION")
				return
			}
			log.Println(res.InsertedID)
		} else if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "ERROR DATABASE CONNECTION")
			return
		} else {
			cid = result["creator_id"].(string)
		}
	}

	tokenString, err := createJWT(jwt.MapClaims{
		"id":     cid,
		"expire": time.Now().Add(1 * time.Hour).Format(time.RFC1123Z),
	})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "ERROR CREATING TOKEN")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func refreshTokenLogin(c *gin.Context) {
	userid, err := validateToken(c.GetHeader("Authorization"), true, false)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	if userid != "" {
		var result bson.M
		err = mongoClient.Database("npuzzle").Collection("Akun").
			FindOne(context.TODO(),
				bson.D{{Key: "creator_id", Value: userid}}, options.FindOne().SetProjection(bson.D{{Key: "ftkn", Value: 0}, {Key: "dtkn", Value: 0}})).
			Decode(&result)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, "USER ID NOT FOUND")
			return
		}
	}

	tokenString, err := createJWT(jwt.MapClaims{
		"id":     userid,
		"expire": time.Now().Add(1 * time.Hour).Format(time.RFC1123Z),
	})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "ERROR CREATING TOKEN")
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

//curl -X GET http://localhost:8080/login?ftkn=eyJhbGciOiJSUzI1NiIsImtpZCI6IjJkMjNmMzc0MDI1ZWQzNTNmOTg0YjUxMWE3Y2NlNDlhMzFkMzFiZDIiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiTW9ubW9uIE1uZW1vbmljIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hLS9BRmRadWNvdlVDd3E5QjVKWWJqUC1MVHNwVGc4bGRDakZPd29YNk1BcFVWUj1zOTYtYyIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS9ucHV6emxlLWUzNjVjIiwiYXVkIjoibnB1enpsZS1lMzY1YyIsImF1dGhfdGltZSI6MTY2MzEyMjAzOSwidXNlcl9pZCI6Im94ckVXeUhTUmJURUpVZEZSQURhcVByTjFYTzIiLCJzdWIiOiJveHJFV3lIU1JiVEVKVWRGUkFEYXFQck4xWE8yIiwiaWF0IjoxNjYzMTIyMDM5LCJleHAiOjE2NjMxMjU2MzksImVtYWlsIjoibWlrYXR6dWtpYXRoYW5hc0BnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJnb29nbGUuY29tIjpbIjEwNDAzNTcxMzg5NTg2MDYwOTUzMSJdLCJlbWFpbCI6WyJtaWthdHp1a2lhdGhhbmFzQGdtYWlsLmNvbSJdfSwic2lnbl9pbl9wcm92aWRlciI6Imdvb2dsZS5jb20ifX0.F-7T2k7d-cYEx-l-j-ppRvH6yrchXMDDHy7LmFqdjefVnEYmIaaDmtpqhvbU5xRaAQpUHHRhLaont7UezHtv9mjm71wfnNlNTCW6azLws4R_D1nWEB5uCXwW8EgPe7-AQn2cApHOJJhiRzEYIYy5GXdK_p_UB55AEaxotdKOl0RErUf23LguWuB9MqPbEQoSJ7dC9bDUENaDCzHnqCfFh6VFaLcFQfFxgw7RncsSeYQ_OgxM5mLvzP0a9LpPZEI26uYmpj3wLWnNQITXu7kLUHtevIoZTfImvSfHT0t0uez7XnaYtTKMCVNlv1gCUoxGNa50rTNjIIy9y-9Hdxqa5Q
//curl -X GET --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmUiOiJXZWQsIDE0IFNlcCAyMDIyIDE0OjIxOjQ4ICswOTAwIiwiaWQiOiI3MGQzYTQzNi1jMGFlLTRmOGQtOWU0NS0zYmIyMThjOTFmZTUifQ.69AQB5HBAlh31bcyPJnPdM1aI6d75pNybng8Q9uuhkw" http://localhost:8080/login/refresh
