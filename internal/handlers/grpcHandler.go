package handlers

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"realAPI/internal/database"
	"realAPI/internal/database/models"
	"realAPI/packages/proto/pb"
)

const mongodb_uri = "mongodb+srv://user:user@godatabase.jia0iua.mongodb.net/?retryWrites=true&w=majority&appName=GoDatabase"

type Server struct {
	pb.BatchServiceServer
}

func (s *Server) AddWords(ctx context.Context, req *pb.AddWordsRequest) (*pb.AddWordsReply, error) {

	clientOpts := options.Client().ApplyURI(mongodb_uri).SetMaxPoolSize(20)
	pooledClient, err := mongo.Connect(context.TODO(), clientOpts)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := pooledClient.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	mongoClient := database.NewServer(pooledClient)

	var words []interface{}

	for _, word := range req.Words {
		words = append(words, models.Word{
			Word:         word.Word,
			Translations: word.Translations,
			Description:  word.Description,
		})
	}

	insert, err := mongoClient.BatchWordsInsert(words)
	if err != nil {
		return nil, err
	}

	return &pb.AddWordsReply{
		Message: fmt.Sprintf("Inserted words, %s", insert),
	}, nil
}
