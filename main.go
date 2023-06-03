package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var jwtSecret string

var mongoClient mongo.Client
var fbClient firebase.App
var authClient auth.Client

func initClient() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("NOT USING .env FILE")
	} else {
		log.Println("USING .env FILE")
	}

	val, ok := os.LookupEnv("mongoUser")
	if !ok {
		log.Fatal("NEED ENV mongoUser")
	}
	mongoUser := val

	val, ok = os.LookupEnv("mongoPass")
	if !ok {
		log.Fatal("NEED ENV mongoPass")
	}
	mongoPass := val

	val, ok = os.LookupEnv("jwtSecret")
	if !ok {
		log.Fatal("NEED ENV jwtSecret")
	}
	jwtSecret = val

	fbauth := ""
	val, ok = os.LookupEnv("fbauth")
	if !ok {
		val, ok = os.LookupEnv("fbauthenc")
		if !ok {
			log.Fatal("NEED ENV fbauth")
		} else {
			fbauth = val
			val, ok = os.LookupEnv("fbauthkey")
			if !ok {
				log.Fatal("NEED ENV fbauth")
			}
			token, err := jwt.Parse(fbauth, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Fatal("ERROR PARSING FBAUTH")
				}

				return []byte(val), nil
			})

			if err != nil {
				log.Fatal("ERROR PARSING FBAUTH")
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				fbauth = claims["val"].(string)
			} else {
				log.Fatal("ERROR PARSING FBAUTH")
			}
		}
	} else {
		fbauth = val
	}

	// INITIALIZE MONGODB
	var mongoDBConn string = "mongodb+srv://" + mongoUser + ":" + mongoPass + "@cluster0.axykj.mongodb.net/?retryWrites=true&w=majority"

	conClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDBConn))
	if err != nil {
		log.Fatal(err)
	}
	if err = conClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	mongoClient = *conClient

	// INITIALIZE FIREBASE
	opt := option.WithCredentialsJSON([]byte(fbauth))
	fb, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal(err)
	}

	fbClient = *fb

	client, err := fbClient.Auth(context.TODO())
	if err != nil {
		log.Fatal(err)
		return
	}
	authClient = *client
}

func contains(checkArray [5]string, target string) bool {
	for _, e := range checkArray {
		if e == target {
			return true
		}
	}
	return false
}

func main() {
	initClient()

	router := gin.Default()
	var allowedHost = [...]string{"http://127.0.0.1", "http://localhost", "http://178.128.219.41", "https://npuzzle.3mworkshop.com", "http://npuzzle.3mworkshop.com"}

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			hostArray := strings.Split(origin, ":")
			return contains(allowedHost, hostArray[0]+":"+hostArray[1])
		},
	}))

	router.GET("/account/login", getLogin)
	router.GET("/account/login/refresh", refreshTokenLogin)
	router.POST("/account", postAccHome)
	router.GET("/account", getAccHome)

	router.POST("/puzzle", postPuzzleHome)
	router.GET("/puzzle", getPuzzleHome)
	router.PUT("/puzzle", putPuzzleHome)
	router.GET("/puzzle/state", getStateHome)
	router.POST("/puzzle/state", postStateHome)
	router.POST("/puzzle/favourite", postFavHome)
	router.GET("/puzzle/:puzzle_id", getPuzzle)
	router.POST("/puzzle/:puzzle_id", postAnswer)
	router.DELETE("/puzzle/:puzzle_id", deletePuzzle)

	router.Run("0.0.0.0:8080")
}
