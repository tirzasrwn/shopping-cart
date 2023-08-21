package logs

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbTimeout = 5 * time.Second

var client *mongo.Client
var LogEntries LogEntryFunc

func New(mongo *mongo.Client) error {
	client = mongo
	LogEntries = &LogEntry{}
	return nil
}

type LogEntry struct {
	ID             string      `bson:"_id,omitempty" json:"id,omitempty"`
	RequestMethod  string      `json:"request_method" bson:"request_method"`
	RequestURL     string      `json:"request_url" bson:"request_url"`
	RequestBody    interface{} `json:"request_body" bson:"request_body"`
	RequestHeader  interface{} `json:"request_header" bson:"request_header"`
	ResponseBody   interface{} `json:"response_body" bson:"response_body"`
	ResponseCode   int         `json:"response_code" bson:"response_code"`
	ResponseHeader interface{} `json:"response_header" bson:"response_header"`
	User           string      `json:"user" bson:"user"`
	CreatedAt      time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time   `bson:"updated_at" json:"updated_at"`
}

type LogEntryFunc interface {
	Insert() error
	All() ([]*LogEntry, error)
	GetOne(id string) (*LogEntry, error)
	DropCollection() error
	Update() (*mongo.UpdateResult, error)
	LogThisRequest(c *gin.Context, status int, data interface{}) error
}

func (l *LogEntry) Insert() error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		RequestMethod:  l.RequestMethod,
		RequestURL:     l.RequestURL,
		RequestBody:    l.RequestBody,
		RequestHeader:  l.RequestHeader,
		ResponseBody:   l.ResponseBody,
		ResponseCode:   l.ResponseCode,
		ResponseHeader: l.ResponseHeader,
		User:           l.User,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into log:", err)
		return err
	}
	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}
	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, err
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docId},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "request_method", Value: l.RequestMethod},
				{Key: "request_url", Value: l.RequestURL},
				{Key: "request_body", Value: l.RequestBody},
				{Key: "request_header", Value: l.RequestHeader},
				{Key: "response_body", Value: l.ResponseBody},
				{Key: "response_code", Value: l.ResponseCode},
				{Key: "response_header", Value: l.ResponseHeader},
				{Key: "user", Value: l.User},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getUserEmailFromContex(c *gin.Context) (string, error) {
	// get auth header
	authHeader := c.Request.Header.Get("Authorization")
	// sanity check
	if authHeader == "" {
		return "", errors.New("no auth header")
	}
	// split the header on spaces
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", errors.New("invalid auth header")
	}
	// check to see if we have the word Bearer
	if headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}
	tokenString := headerParts[1]
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}
	var userEmail string
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userEmail = claims["email"].(string)
	} else {
		return "", nil
	}
	return userEmail, nil
}

func (l *LogEntry) LogThisRequest(c *gin.Context, status int, data interface{}) error {
	userEmail, _ := getUserEmailFromContex(c)
	if userEmail == "" {
		userEmail = "public"
	}
	// Get request body.
	requestBodyRaw, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	var result interface{}
	json.Unmarshal(requestBodyRaw, &result)

	le := LogEntry{
		RequestMethod: c.Request.Method,
		RequestURL:    c.Request.RequestURI,
		RequestBody:   result,
		RequestHeader: c.Request.Header,
		ResponseBody:  data,
		ResponseCode:  status,
		User:          userEmail,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return le.Insert()
}
