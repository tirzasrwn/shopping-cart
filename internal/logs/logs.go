package logs

import (
	"context"
	"errors"
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
	ID         string      `bson:"_id,omitempty" json:"id,omitempty"`
	Method     string      `json:"method" bson:"method"`
	Param      string      `json:"param" bson:"param"`
	Response   interface{} `json:"response" bson:"response"`
	StatusCode int         `json:"status_code" bson:"status_code"`
	User       string      `json:"user" bson:"user"`
	CreatedAt  time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time   `bson:"updated_at" json:"updated_at"`
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
		Method:     l.Method,
		Param:      l.Param,
		Response:   l.Response,
		StatusCode: l.StatusCode,
		User:       l.User,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
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
				{Key: "method", Value: l.Method},
				{Key: "param", Value: l.Param},
				{Key: "response", Value: l.Response},
				{Key: "status_code", Value: l.StatusCode},
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
	le := LogEntry{
		Method:     c.Request.Method,
		Param:      c.Request.RequestURI,
		Response:   data,
		StatusCode: status,
		User:       userEmail,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	return le.Insert()
}
