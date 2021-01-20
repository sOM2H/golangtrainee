package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dsn = "root:root@tcp(127.0.0.1:3306)/golangbeginner?charset=utf8mb4&parseTime=True&loc=Local"
)

var db *gorm.DB

// create table posts(id INT NOT NULL, user_id INT NOT NULL, title VARCHAR(100) NOT NULL, body VARCHAR(8000) NOT NULL);
type Post struct {
	Id     int    `json:"id"`
	UserId int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// create table comments(id INT NOT NULL, post_id INT NOT NULL, name VARCHAR(100) NOT NULL, email VARCHAR(100) NOT NULL, body VARCHAR(8000) NOT NULL);
type Comment struct {
	Id     int    `json:"id"`
	PostId int    `json:"postId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "not found")
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var post Post

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error")

	}
	json.Unmarshal(reqBody, &post)

	result := db.Save(&post)
	if result.Error != nil {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(post)

}

func getPosts(w http.ResponseWriter, r *http.Request) {
	var posts []Post
	db.Find(&posts)
	json.NewEncoder(w).Encode(posts)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]

	var post Post

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error")

	}
	json.Unmarshal(reqBody, &post)

	post.Id, _ = strconv.Atoi(postID)

	result := db.Save(&post)
	if result.Error != nil {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(post)

}

func getPost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]
	var post Post
	err := db.First(&post, postID)
	if err.Error != nil {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(post)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["id"]
	err := db.Delete(&Post{}, postID)
	if err.Error != nil {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	fmt.Fprint(w, "Deleted")
}

func createComment(w http.ResponseWriter, r *http.Request) {
	var comment Comment

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error")

	}
	json.Unmarshal(reqBody, &comment)

	result := db.Save(&comment)
	if result.Error != nil {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(comment)

}

func getComments(w http.ResponseWriter, r *http.Request) {
	var comments []Comment
	db.Find(&comments)
	json.NewEncoder(w).Encode(comments)
}

func updateComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]

	var comment Comment

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error")

	}
	json.Unmarshal(reqBody, &comment)

	comment.Id, _ = strconv.Atoi(commentID)

	result := db.Save(&comment)
	if result.Error != nil {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(comment)

}

func getComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]
	var comment Comment
	err := db.First(&comment, commentID)
	if err.Error != nil {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(comment)
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := mux.Vars(r)["id"]
	err := db.Delete(&Comment{}, commentID)
	if err.Error != nil {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	fmt.Fprint(w, "Deleted")
}

func main() {
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
		return
	}
	db = database

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/posts", getPosts).Methods("GET")
	router.HandleFunc("/posts/{id}", getPost).Methods("GET")
	router.HandleFunc("/posts/{id}", updatePost).Methods("PATCH")
	router.HandleFunc("/posts/{id}", deletePost).Methods("DELETE")

	router.HandleFunc("/comments", createComment).Methods("POST")
	router.HandleFunc("/comments", getComments).Methods("GET")
	router.HandleFunc("/comments/{id}", getComment).Methods("GET")
	router.HandleFunc("/comments/{id}", updateComment).Methods("PATCH")
	router.HandleFunc("/comments/{id}", deleteComment).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
