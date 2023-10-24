package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Developer struct {
	Name  string `json:"name"`
	Year  string `json:"year"`
	Major string `json:"major"`
	About string `json:"about"`
}

func getDevelopers(c *gin.Context) {
	// Open a connection to PlanetScale
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
		return
	} else {
		log.Println("Connected to PlanetScale...")
	}

	// Query the database
	rows, err := db.Query("SELECT * FROM developers")
	if err != nil {
		log.Fatal("(getDevelopers) db.Query:", err)
		return
	}
	defer rows.Close()

	developers := []Developer{}
	for rows.Next() {
		var developer Developer
		err := rows.Scan(&developer.Name, &developer.Year, &developer.Major, &developer.About)
		if err != nil {
			log.Fatal("(getDevelopers) res.Scan:", err)
			return
		}
		developers = append(developers, developer)
	}

	c.JSON(http.StatusOK, developers)
}

func postDevelopers(c *gin.Context) {
	// validate key before proceeding
	if c.GetHeader("key") != os.Getenv("ADMIN_KEY") {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Open a connection to PlanetScale
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
		return
	} else {
		log.Println("Connected to PlanetScale...")
	}

	// bind input json
	var newDeveloper Developer
	err = c.BindJSON(&newDeveloper)
	if err != nil {
		log.Fatal("(postDevlopers) c.BindJSON:", err)
		return
	}

	// Insert into the database
	query := `INSERT INTO developers (name, year, major, about) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, newDeveloper.Name, newDeveloper.Year, newDeveloper.Major, newDeveloper.About)
	if err != nil {
		log.Fatal("(postDevelopers) db.Exec:", err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func deleteDevelopers(c *gin.Context) {
	// validate key before proceeding
	if c.GetHeader("key") != os.Getenv("ADMIN_KEY") {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Open a connection to PlanetScale
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
		return
	} else {
		log.Println("Connected to PlanetScale...")
	}

	// bind input json
	var developerName string
	err = c.BindJSON(&developerName)
	if err != nil {
		log.Fatal("(deleteDevelopers) c.BindJSON:", err)
		return
	}

	// Insert into the database
	query := `DELETE FROM developers WHERE name = ?`
	res, err := db.Exec(query, developerName)
	if err != nil {
		log.Fatal("(deleteDevelopers) db.Exec:", err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func validate(c *gin.Context) {
	pin := c.PostForm("pin")
	new_url := os.Getenv("FRONTEND_URL") + "/game/" + pin
	fmt.Println(new_url)
	c.Redirect(http.StatusMovedPermanently, new_url)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func gameSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade to websocket:", err)
		return
	}
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("failed to read message: ", err)
			return
		}
		fmt.Println("message type: ", t)
		fmt.Println("message: ", string(msg))
	}
}

func gameHandler(c *gin.Context) {
	gameId := c.Param("gameId")
	fmt.Println(gameId)
	fmt.Println("Game ID: ", gameId)
	fmt.Println("------------------")
	fmt.Println("Websocket Test:")
	gameSocket(c.Writer, c.Request)
}

func main() {

	fmt.Println("Starting server...")

	// Load connection string from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env", err)
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(cors.Default())

	// default route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": "hungrymonkey",
		})
	})

	// route for developers
	r.GET("/developers", getDevelopers)
	// ! POST and DELETE routes not tested yet
	r.POST("/developers", postDevelopers)
	r.DELETE("/developers", deleteDevelopers)

	// route pin
	r.POST("/validate", validate)

	// game web socket
	r.GET("/game/gameId", gameHandler)

	r.Run(":8000")

}
