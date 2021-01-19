package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	postsURL    = "http://jsonplaceholder.typicode.com/posts?userId="
	commentsURL = "http://jsonplaceholder.typicode.com/comments?postId="
	dsn         = "root:root@tcp(127.0.0.1:3306)/golangbeginner?charset=utf8mb4&parseTime=True&loc=Local"
)

type Post struct {
	Id     int    `json:"id"`
	UserId int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comment struct {
	Id     int    `json:"id"`
	PostId int    `json:"postId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

var wg sync.WaitGroup

func insert(model interface{}, db *gorm.DB) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		result := db.Create(model)
		if result.Error != nil {
			log.Println(result)
		}
	}()

}

func unmarshalResponse(url string, i interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &i)
	return err
}

func getComments(postId int, db *gorm.DB) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		var comments []Comment
		err := unmarshalResponse(commentsURL+strconv.Itoa(postId), &comments)
		if err != nil {
			log.Println(err)
			return
		}

		for _, comment := range comments {
			insert(comment, db)
		}
	}()
}

func main() {
	var posts []Post

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
		return
	}

	err = unmarshalResponse(postsURL+strconv.Itoa(7), &posts)
	if err != nil {
		log.Println(err)
		return
	}

	for _, post := range posts {
		insert(post, db)
		getComments(post.Id, db)
	}

	wg.Wait()
}
