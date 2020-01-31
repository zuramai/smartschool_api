package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"gonum.org/v1/gonum/floats"

	"github.com/gorilla/mux"

	"github.com/zuramai/smartschool_api/models"
	"go.mongodb.org/mongo-driver/bson"
)

var users []models.User

func UserIndex(w http.ResponseWriter, r *http.Request) {
	// timeStart := time.Now()
	users := []models.User{}
	user := models.User{}

	ctx := context.TODO() // Options to the database.
	coll, err := models.GetDB("main").Collection("users").Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
	}

	for coll.Next(ctx) {
		coll.Decode(&user)
		users = append(users, user)

		user = models.User{}
	}
	// timeEnd := time.Since(timeStart)
	// fmt.Println("Time Elapsed: ", timeEnd)
	respondJSON(w, 200, "Success get all data users", users)

}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	// username := r.FormValue("username")
	// password := r.FormValue("password")
	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	// name := r.FormValue("name")
	// phone := r.FormValue("phone")

	// var photoName string
	// photo, _, err := r.FormFile("photo")
	// if err != nil {
	// 	photoName = "default.png"
	// } else {
	// 	dir, _ := os.Getwd()
	// 	timeNowMs := strconv.FormatInt(makeTimestampMilli(), 10)
	// 	photoName := timeNowMs + ".png"
	// 	photoDir := dir + "/assets/images/attendances/" + photoName
	// 	f, err := os.Create(photoDir)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	io.Copy(f, photo)
	// }

	// roleString := r.FormValue("role")
	// var role uint
	// if roleString == "Superadmin" {
	// 	role = 4
	// } else if roleString == "Admin" {
	// 	role = 3
	// } else if roleString == "Guru" {
	// 	role = 2
	// } else {
	// 	role = 1
	// }

	// user := models.User{
	// 	Username:  username,
	// 	Password:  string(hashedPassword),
	// 	Name:      name,
	// 	Phone:     phone,
	// 	PhotoName: photoName,
	// 	RoleID:    role,
	// 	Status:    true,
	// }

	// models.GetDB("main").NewRecord(user)
	// models.GetDB("main").Create(&user)

	// respondJSON(w, 200, "Successfully register user!", user)

}

func UserStore(w http.ResponseWriter, r *http.Request) {

}

func UserDetail(w http.ResponseWriter, r *http.Request) {

}

func UserUpdate(w http.ResponseWriter, r *http.Request) {

}

func UserDelete(w http.ResponseWriter, r *http.Request) {

}

func UserV2Index(w http.ResponseWriter, r *http.Request) {
	timeStart := time.Now()
	users := []models.User{}
	user := models.User{}

	ctx := context.TODO() // Options to the database.
	coll, err := models.GetDB("main").Collection("users").Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
	}

	for coll.Next(ctx) {
		coll.Decode(&user)
		users = append(users, user)

		user = models.User{}
	}
	timeEnd := time.Since(timeStart)
	fmt.Println("Time Elapsed: ", timeEnd)
	respondJSON(w, 200, "Success get all data users", users)
}

func UserV2Detail(w http.ResponseWriter, r *http.Request) {

	userID := mux.Vars(r)["id"]

	fmt.Println(userID)
	var user models.User
	err := models.GetDB("main").Collection("users").FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(&user)

	if err != nil {
		fmt.Println(err)
		respondJSON(w, 200, "User not found", map[string]interface{}{})
		return
	}

	respondJSON(w, 200, "Get User Data Embeddings", user)
	return
}

func UserV2Verify(w http.ResponseWriter, r *http.Request) {
	var user models.UserVerify
	var checkUser models.User
	var checkUserValidation models.User
	// photo, _, errFile := r.FormFile("photo")
	// if errFile != nil {
	// 	respondErrorValidationJSON(w, 422, "Error", map[string]interface{}{
	// 		"photo": "Photo is required",
	// 	})
	// }
	// fmt.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&user)
	stringUserID := strconv.Itoa(user.UserID)

	errCheckUser := models.GetDB("main").Collection("users").FindOne(context.TODO(), bson.M{"user_id": stringUserID}).Decode(&checkUserValidation)
	if errCheckUser != nil {
		fmt.Println(errCheckUser)
	}
	if checkUserValidation.Status == 1 {
		respondJSON(w, 422, "You are already verified", map[string]interface{}{})
		return
	}
	fmt.Println(checkUser)

	if err != nil && err != io.EOF {
		// fmt.Println(err)
		fmt.Println("Error", err.Error())
		respondJSON(w, 422, "Error", map[string]interface{}{})
		return
	}
	// fmt.Println(user)
	errCheck := models.GetDB("main").Collection("users").FindOneAndUpdate(context.TODO(), bson.M{"user_id": stringUserID}, bson.M{"$set": bson.M{"embeddings": user.Embeddings, "status": 1}}).Decode(&checkUser)

	if errCheck != nil {
		fmt.Println("not found", errCheck.Error())
		// fmt.Println("Input string : ", user)
		respondJSON(w, 422, "Verify failed", map[string]interface{}{})
		fmt.Println("not found", errCheck)
		fmt.Println("Input string : ", user)
		respondJSON(w, 422, "Verify failed "+errCheck.Error(), map[string]interface{}{})
		return
	}

	respondJSON(w, 200, "Verify success!", map[string]interface{}{
		"user": checkUser,
	})
	return
}

func UserV2Embeddings(w http.ResponseWriter, r *http.Request) {
	var userEmbedding models.UserEmbeddings
	var userEmbeddings []models.UserEmbeddings
	cursor, err := models.GetDB("main").Collection("users").Find(context.TODO(), bson.M{})

	if err != nil {
		fmt.Println(err)
	}

	for cursor.Next(context.TODO()) {
		cursor.Decode(&userEmbedding)
		userEmbeddings = append(userEmbeddings, userEmbedding)
		userEmbedding = models.UserEmbeddings{}
	}

	respondJSON(w, 200, "Get All Data Embeddings", userEmbeddings)
	return
}

func UserRecognize(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var userEmbedding models.UserEmbeddings
	var userEmbeddings []models.UserEmbeddings
	var recognition models.UserRecognition
	var recognitionList []models.UserRecognition
	var log models.Log

	json.NewDecoder(r.Body).Decode(&recognition)

	if len(recognition.Embedding) == 0 {
		respondErrorValidationJSON(w, 422, "Input Embedding Null", map[string]interface{}{})
		return
	}

	cursor, err := models.GetDB("main").Collection("users").Find(context.TODO(), bson.M{})

	if err != nil {
		fmt.Println(err)
	}

	for cursor.Next(context.TODO()) {
		cursor.Decode(&userEmbedding)
		userEmbeddings = append(userEmbeddings, userEmbedding)
		userEmbedding = models.UserEmbeddings{}
	}

	for _, UserEmbeddingList := range userEmbeddings {
		// Index := index
		if len(UserEmbeddingList.Embeddings) == 0 {
			continue
		} else {
			var val []float64
			for _, embeddingList := range UserEmbeddingList.Embeddings {
				val = append(val, euclideanDistance(embeddingList, recognition.Embedding))
			}
			recognition.UserID = UserEmbeddingList.UserID
			recognition.Name = UserEmbeddingList.Name
			recognition.Accuracy = floats.Min(val)
			recognition.Elapsed = time.Since(start).String()
			recognitionList = append(recognitionList, recognition)
		}
		// log.Println(maximum)
	}
	if len(recognitionList) == 0 {
		respondErrorValidationJSON(w, 422, "Cannot Recognize Face!", map[string]interface{}{})
		return
	} else {
		var acculist []float64
		for _, value := range recognitionList {
			acculist = append(acculist, value.Accuracy)
		}
		res := recognitionList[floats.MinIdx(acculist)]
		fmt.Println(res.UserID)

		if res.Accuracy <= 0.2 {
			attendanceBody := models.AttendanceBody{
				UserID:        res.UserID,
				CameraID:      res.CameraID,
				PhotoEncoding: res.PhotoEncoding,
			}

			log = models.Log{
				UserID:   res.UserID,
				CameraID: res.CameraID,
			}
			newAttendance(w, attendanceBody)
			insertLog := logStore(log)
			fmt.Println("insert log :", insertLog)
		}

		respondJSON(w, 200, "Returned Matching Identities", map[string]interface{}{
			"user_id":  res.UserID,
			"name":     res.Name,
			"accuracy": res.Accuracy,
			"elapsed":  res.Elapsed,
		})
		return
	}
}

func euclideanDistance(emb1, emb2 []float64) float64 {
	val := 0.0
	for i := range emb1 {
		val += math.Pow(emb1[i]-emb2[i], 2)
	}
	return val
}

func UserV2EmbeddingsClear(w http.ResponseWriter, r *http.Request) {
	res, err := models.GetDB("main").Collection("users").UpdateMany(context.TODO(), bson.M{}, bson.M{"$set": bson.M{"embeddings": ""}})
	if err != nil {
		fmt.Println(err)
	}
	respondJSON(w, 200, "Success Clear All Embeddings", res)
}

func UserV2EmbeddingsClearOnUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	res, err := models.GetDB("main").Collection("users").UpdateOne(context.TODO(), bson.M{"user_id": userID}, bson.M{"$set": bson.M{"embeddings": ""}})
	if err != nil {
		fmt.Println(err)
	}
	respondJSON(w, 200, "Success Clear User Embeddings", res)
}
