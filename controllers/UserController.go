package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zuramai/smartschool_api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var users []models.User

func UserIndex(w http.ResponseWriter, r *http.Request) {
	// search := r.URL.Query().Get("search")
	// // _, perPage, offset := u.Paginate(r)
	// models.GetDB("master").Find(&users)
	// // models.GetDB("master").Where("name LIKE ?", "%"+search+"%").Offset(offset).Limit(perPage).Order("id asc").Find(&users)
	// respondJSON(w, 200, "Success get all data users, search: "+search, users)
}

type UserMongo struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}

func UserIndexMongo(w http.ResponseWriter, r *http.Request) {
	const (
		// Name of the database.
		DBName = "smart_school"
		URI    = "mongodb://localhost/smart_school"
	)

	ctx := context.Background() // Options to the database.
	clientOpts := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		fmt.Println(err)
		return
	}
	db := client.Database(DBName)

	notesResult := []UserMongo{}
	n := UserMongo{}
	coll := db.Collection("users")
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return
	} // Iterate through the returned cursor.
	for cursor.Next(ctx) {
		cursor.Decode(&n)
		notesResult = append(notesResult, n)
	}

	respondJSON(w, 200, "Success get all data users", notesResult)

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

}

func UserV2Detail(w http.ResponseWriter, r *http.Request) {

}

func UserV2Verify(w http.ResponseWriter, r *http.Request) {

}
