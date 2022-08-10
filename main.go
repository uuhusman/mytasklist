package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"database/sql"
	// "encoding/json"
	"fmt"
	// "html/template"
	// "path"

	_ "github.com/go-sql-driver/mysql"
)

type task struct {
	Id       int
	Task     string
	Assignee string
	Deadline string
	Status   string
}

func connect() (*sql.DB, error) {
	// db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_task")
	db, err := sql.Open("mysql", "freedb_usman:?Zmp67zTwEp7$%J@tcp(sql.freedb.tech:3306)/freedb_db_task")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	port := os.Getenv("PORT")
	// port := "8080";

	if port == "" {
		// log.Fatal("$PORT must be set")
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	// router := gin.Default()
    router.POST("/get_data", getData)
    router.POST("/save", saveData)
    router.POST("/update", updateData)

	// http.HandleFunc("/get_data", ActionData)
	// http.HandleFunc("/save", handleSave)
	// http.HandleFunc("/update", handleUpdate)

	router.Run(":" + port)

	// fmt.Println(port)
}

func getData(c *gin.Context) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	// var id = 1
	// rows, err := db.Query("select id, task, assignee, deadline, status from tbl_tasking where id = ?", id)
	rows, err := db.Query("select Id, Task, Assignee, Deadline, Status from tbl_tasking")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	var result []task

	for rows.Next() {
		var each = task{}
		var err = rows.Scan(&each.Id, &each.Task, &each.Assignee, &each.Deadline, &each.Status)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		result = append(result, each)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}

    c.IndentedJSON(http.StatusOK, result)
}

func saveData(c *gin.Context){

	// var tasklist task

		payload := struct {
			Id       int    `json:"id"`
			Task     string `json:"task"`
			Assignee string `json:"assignee"`
			Deadline string `json:"deadline"`
		}{}

		err := c.BindJSON(&payload)
		if err != nil {
			log.Fatal(err)
		}

		db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		if payload.Id > 0 {
			_, err = db.Exec("update tbl_tasking set Task = ?, Assignee = ? , Deadline = ?  where id = ?", payload.Task, payload.Assignee, payload.Deadline, payload.Id)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("update success!")
		} else {
			_, err = db.Exec("INSERT INTO `tbl_tasking` (`Task`, `Assignee`, `Deadline`) VALUES (?, ?, ?)", payload.Task, payload.Assignee, payload.Deadline)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("insert success!")
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success"})
}

func updateData(c *gin.Context){
	payload := struct {
		Id     int    `json:"id"`
		Status string `json:"status"`
	}{}

	err := c.BindJSON(&payload)
	if err != nil {
		log.Fatal(err)
	}

	db, err := connect()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer db.Close()

		if payload.Status != "delete" {
			_, err = db.Exec("update tbl_tasking set Status = ? where id = ?", payload.Status, payload.Id)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Println("update success!")
		} else {
			_, err = db.Exec("delete from tbl_tasking where id = ?", payload.Id)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Println("delete success!")
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success"})
}