package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Card struct {
	FileNo     string
	Codename   string
	Class      string
	Rare       int
	Gender     string
	IsPickedUp bool
}

func main() {
	rand.Seed(time.Now().UnixNano())
	router := gin.Default()
	router.Use(Cors())
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")
	router.GET("/user", handleUser)
	router.GET("/pickOne", handlePickOne)
	router.Run(":8966")
	//mongoURI := "mongodb://localhost:27017"
	//filter := bson.M{"isPickedUp":true}
	//update := bson.D{{"$set", bson.D{{"isPickedUp",false}}}}

	//findManyFromDB(mongoURI,"arknights","遗愿焰火", filter)

	//updateOneToDB(mongoURI,"arknights","遗愿焰火", filter, update)

}

func handleRare() bson.M {
	p := rand.Float64() * 100
	if p >= 0 && p <= 1.4 {
		return bson.M{
			"rare":       6,
			"isPickedUp": false,
		}
	} else if p <= 2 {
		return bson.M{
			"rare":       6,
			"isPickedUp": false,
		}
	} else if p <= 6 {
		return bson.M{
			"rare":       5,
			"isPickedUp": true,
		}
	} else if p <= 10 {
		return bson.M{
			"rare":       5,
			"isPickedUp": false,
		}
	} else if p <= 60 {
		return bson.M{
			"rare": 4,
		}
	} else {
		return bson.M{
			"rare": 3,
		}
	}
}

func handleUser(c *gin.Context) {
	data := bson.M{
		"name":   c.Query("name"),
		"pass":   c.Query("pass"),
		"isUser": c.Query("isUser"),
	}
	c.JSON(200, data) //gin.H{"name": c.Query("name"), "pass": c.Query("pass"), "isUser": c.Query("isUser"),})
}

func handlePickOne(c *gin.Context) {
	mongoURI := "mongodb://localhost:27017"
	p := handleRare()
	res := findManyFromDB(mongoURI, "arknights", "遗愿焰火", p)
	data := res[rand.Intn(len(res))]
	c.JSON(200, data)
}

func findManyFromDB(MongoUri, Database, Collection string, filter bson.M) []Card {
	clientOptions := options.Client().ApplyURI(MongoUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Mongodb!")
	collection := client.Database(Database).Collection(Collection)
	findOptions := options.Find()
	var results []Card
	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem Card
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	fmt.Println("Connection to MongoDB closed.")
	return results
}

func findOneFromDB(MongoUri, Database, Collection string, filter bson.M) Card {
	clientOptions := options.Client().ApplyURI(MongoUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Mongodb!")
	var result Card
	collection := client.Database(Database).Collection(Collection)
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found a single document: %+v\n", result)
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
	return result
}

func updateOneToDB(MongoUri, Database, Collection string, filter bson.M, update bson.D) {
	clientOptions := options.Client().ApplyURI(MongoUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Mongodb!")
	collection := client.Database(Database).Collection(Collection)
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	if res == nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Matched %v documents and updated %v documents.\n", res.MatchedCount, res.ModifiedCount)
	}
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func insertOneToDB(MongoUri, Database, Collection string, data Card) {
	clientOptions := options.Client().ApplyURI(MongoUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Mongodb!")
	collection := client.Database(Database).Collection(Collection)
	res, err := collection.InsertOne(context.TODO(), data)
	if res == nil {
		log.Fatal(err)
	} else {
		id := res.InsertedID
		fmt.Println("insert complete!: ", id)
	}
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func insertManyToDB(MongoUri, Database, Collection string, data []interface{}) {
	clientOptions := options.Client().ApplyURI(MongoUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Mongodb!")
	collection := client.Database(Database).Collection(Collection)
	res, err := collection.InsertMany(context.TODO(), data)
	if res == nil {
		log.Fatal(err)
	} else {
		fmt.Println("insert complete!: ", res)
	}
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization,"+
				" Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection,"+
				" Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, "+
				"Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, "+
				"Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,"+
				"Expires,Last-Modified,Pragma,FooBar")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}
