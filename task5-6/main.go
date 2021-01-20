package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/echo-swagger/example/docs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dsn = "root:root@tcp(127.0.0.1:3306)/golangbeginner?charset=utf8mb4&parseTime=True&loc=Local"
)

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

var db *gorm.DB

func respond(c echo.Context, code int, i interface{}) error {
	if strings.HasPrefix(c.Request().Header.Get("Content-Type"), "application/xml") {
		return c.XML(code, i)
	}
	return c.JSON(code, i)
}

func getPosts(c echo.Context) error {
	var posts []Post
	db.Find(&posts)

	return respond(c, http.StatusOK, posts)
}

func getPost(c echo.Context) error {
	postID := c.Param("id")
	var post Post
	err := db.First(&post, postID)
	if err.Error != nil {
		return respond(c, http.StatusNotFound, "")
	}

	return respond(c, http.StatusOK, post)
}

func createPost(c echo.Context) error {
	id, _ := strconv.Atoi(c.FormValue("id"))
	user_id, _ := strconv.Atoi(c.FormValue("userId"))
	title := c.FormValue("title")
	body := c.FormValue("body")

	post := Post{Id: id, UserId: user_id, Title: title, Body: body}

	err := db.Create(&post)
	if err.Error != nil {
		return respond(c, http.StatusBadRequest, "")
	}

	return respond(c, http.StatusOK, post)
}

func updatePost(c echo.Context) error {
	id, _ := strconv.Atoi(c.FormValue("id"))
	user_id, _ := strconv.Atoi(c.FormValue("userId"))
	title := c.FormValue("title")
	body := c.FormValue("body")

	post := Post{Id: id, UserId: user_id, Title: title, Body: body}

	err := db.Save(&post)
	if err.Error != nil {
		return respond(c, http.StatusBadRequest, "")
	}

	return respond(c, http.StatusOK, post)

}

func deletePost(c echo.Context) error {
	postID := c.Param("id")
	err := db.Delete(&Post{}, postID)
	if err != nil {
		return respond(c, http.StatusNotFound, "")
	}
	return respond(c, http.StatusOK, "Deleted")
}

func getComments(c echo.Context) error {
	var comments []Comment
	db.Find(&comments)

	return respond(c, http.StatusOK, comments)
}

func getComment(c echo.Context) error {
	commentID := c.Param("id")
	var comment Comment
	err := db.First(&comment, commentID)
	if err.Error != nil {
		return respond(c, http.StatusNotFound, "")
	}

	return respond(c, http.StatusOK, comment)
}

func createComment(c echo.Context) error {
	id, _ := strconv.Atoi(c.FormValue("id"))
	post_id, _ := strconv.Atoi(c.FormValue("postId"))
	name := c.FormValue("name")
	email := c.FormValue("email")
	body := c.FormValue("body")

	comment := Comment{Id: id, PostId: post_id, Name: name, Email: email, Body: body}

	err := db.Create(&comment)
	if err.Error != nil {
		return respond(c, http.StatusBadRequest, "")
	}

	return respond(c, http.StatusOK, comment)
}

func updateComment(c echo.Context) error {
	id, _ := strconv.Atoi(c.FormValue("id"))
	post_id, _ := strconv.Atoi(c.FormValue("postId"))
	name := c.FormValue("name")
	email := c.FormValue("email")
	body := c.FormValue("body")

	comment := Comment{Id: id, PostId: post_id, Name: name, Email: email, Body: body}

	err := db.Save(&comment)
	if err.Error != nil {
		return respond(c, http.StatusBadRequest, "")
	}

	return respond(c, http.StatusOK, comment)
}

func deleteComment(c echo.Context) error {
	commentID := c.Param("id")
	err := db.Delete(&Comment{}, commentID)
	if err != nil {
		return respond(c, http.StatusNotFound, "")
	}
	return respond(c, http.StatusOK, "Deleted")
}

func main() {
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
		return
	}
	db = database

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/posts/", getPosts)
	e.POST("/posts/", createPost)
	e.GET("/posts/:id/", getPost)
	e.PUT("/posts/:id/", updatePost)
	e.DELETE("/posts/:id/", deletePost)

	e.GET("/comments", getComments)
	e.POST("/comments", createComment)
	e.GET("/comments/:id", getComment)
	e.PUT("/comments/:id", updateComment)
	e.DELETE("/commetns/:id", deleteComment)

	e.Logger.Fatal(e.Start(":1323"))
}
