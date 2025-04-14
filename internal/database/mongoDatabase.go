package database

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"realAPI/api/apiUtils"
	"realAPI/api/types"
	"realAPI/internal/database/models"
	"realAPI/packages/proto/pb"
)

//mongodb+srv://user:user@godatabase.jia0iua.mongodb.net/?retryWrites=true&w=majority&appName=GoDatabase

type ServerGRPC struct {
	Database *mongo.Client
	pb.BatchServiceServer
}

type Server struct {
	Database *mongo.Client
}

func NewServerGRPC(ctx context.Context, database *mongo.Client) pb.BatchServiceServer {

	var server = &ServerGRPC{Database: database}

	return server
}

func NewServer(database *mongo.Client) *Server {
	return &Server{Database: database}
}

func (mongo *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	var result bson.M

	//coll := mongo.Client.Database("Words").Collection("WordsCollection")

	err := mongo.Database.Ping(context.TODO(), nil)

	if err != nil {
		log.Error(err)
		apiUtils.WriteError(w, "mongo db ping error", http.StatusInternalServerError)

	} else {
		fmt.Println(result)
		apiUtils.WriteSuccess(w, "MongoDB is up and running")
	}
}

func (mongo *Server) AddWord(w http.ResponseWriter, r *http.Request) {
	var requestModel types.WordRequestModel

	userError := json.NewDecoder(r.Body).Decode(&requestModel)

	if userError != nil {
		apiUtils.RequestErrorHandler(w, userError)
	}

	word := models.Word{
		Word:         requestModel.Word,
		Translations: requestModel.Translations,
		Description:  requestModel.Description,
	}

	coll := mongo.Database.Database("realAPI").Collection("words")

	result, err := coll.InsertOne(context.TODO(), word)

	if err != nil {
		log.Error(err)
		apiUtils.InternalErrorHandler(w)
		return
	}

	apiUtils.WriteSuccess(w, result.InsertedID)
}

func (mongo *Server) UpdateWord(w http.ResponseWriter, r *http.Request) {
	var requestModel types.WordRequestModel

	userError := json.NewDecoder(r.Body).Decode(&requestModel)

	if userError != nil {
		apiUtils.RequestErrorHandler(w, userError)
	}

	filter := bson.D{{"_id", convertHexToObjectId(requestModel.Id)}}

	updateModel := bson.D{{"Word", requestModel.Word}, {"Description", requestModel.Description}, {"Translations", requestModel.Translations}}

	coll := mongo.Database.Database("realAPI").Collection("words")

	result, err := coll.ReplaceOne(context.TODO(), filter, updateModel)

	if err != nil {
		log.Error(err)
		apiUtils.InternalErrorHandler(w)
		return
	}

	apiUtils.WriteSuccess(w, result.ModifiedCount)
}

func (mongo *Server) BatchWordsInsert(words []interface{}) (insertedRowsCount int, err error) {

	coll := mongo.Database.Database("realAPI").Collection("words")

	result, err := coll.InsertMany(context.TODO(), words)

	if err != nil {
		log.Error(err)
		return 0, err
	}

	return len(result.InsertedIDs), nil
}

func convertHexToObjectId(hex string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		log.Println("Invalid id")
	}

	return objectId
}
