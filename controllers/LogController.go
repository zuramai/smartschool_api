package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zuramai/smartschool_api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func LogIndex(w http.ResponseWriter, r *http.Request) {
	var logs []models.Log
	var log models.Log
	cur, err := models.GetDB("main").Collection("logs").Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println(err)
	}
	for cur.Next(context.TODO()) {
		cur.Decode(&log)
		logs = append(logs, log)
		log = models.Log{}
	}

	respondJSON(w, 200, "Success Retrieve All Logs Data", logs)
}

func LogStore(w http.ResponseWriter, r *http.Request) {

}

func logStore(log models.Log) bool {
	var checkLog models.Log
	err := models.GetDB("main").Collection("users").FindOne(context.TODO(), bson.M{"user_id": log.UserID}).Decode(&checkLog)

	if err != nil {
		// no data
		models.GetDB("main").Collection("users").InsertOne(context.TODO(), &log)
		return true
	}
	return false
}

func LogDetail(w http.ResponseWriter, r *http.Request) {

}
