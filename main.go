package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
)

//http://phalt.co/a-simple-api-in-go/

var dbmap = initDb()

// Team class
type Team struct {
	ID      int64 `db:"team_id"`
	Created int64
	Name    string
	Color   string
	Score   int64
}

// AddTeam creates a new team
func AddTeam(name string, color string) Team {
	team := Team{
		Created: time.Now().UnixNano(),
		Name:    name,
		Color:   color,
		Score:   0,
	}

	err := dbmap.Insert(&team)
	checkErr(err, "Insert failed")
	return team
}

// AddScore adds the score by 1
func AddScore(name string) {
	_, err := dbmap.Exec("update teams set score = score + 1 where name = ?", name)
	checkErr(err, "Update failed")
}

// GetTeams fetches all teams
func GetTeams() []Team {
	var teams []Team
	_, err := dbmap.Select(&teams, "select * from teams order by score")
	checkErr(err, "Select failed")
	return teams
}

// ClearScore clears the scores to 0
func ClearScore() {
	_, err := dbmap.Exec("update teams set score = 0")
	checkErr(err, "Update failed")
}

// TODO score clear

// ScorePut adding a score
func ScorePut(c *gin.Context) {
	var json Team
	c.Bind(&json)
	AddScore(json.Name)
	content := gin.H{"result": "Success"}
	c.JSON(200, content)
}

// ScorePost clears the score
func ScorePost(c *gin.Context) {
	ClearScore()
	content := gin.H{"result": "Success"}
	c.JSON(200, content)
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(Team{}, "teams").SetKeys(true, "ID")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	count, err := dbmap.SelectInt("select count(*) from teams")
	checkErr(err, "Create tables failed")
	if count == 0 {
		team := Team{
			Created: time.Now().UnixNano(),
			Name:    "red",
			Color:   "rgba(223,53,53,0)",
			Score:   0,
		}
		err = dbmap.Insert(&team)
		checkErr(err, "insert failed")

		team = Team{
			Created: time.Now().UnixNano(),
			Name:    "white",
			Color:   "rgba(230,231,233,0)",
			Score:   0,
		}
		err = dbmap.Insert(&team)
		checkErr(err, "insert failed")
	}

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func main() {
	app := gin.Default()
	app.PUT("/addscore", ScorePut)
	app.POST("/clearscore", ScorePost)

	content := gin.H{}
	for k, v := range GetTeams() {
		content[strconv.Itoa(k)] = gin.H{
			"Name":  v.Name,
			"Color": v.Color,
			"Score": v.Score,
		}
	}

	app.LoadHTMLGlob("templates/*")
	app.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
			"teams": content,
		})
	})

	app.Run(":8000")
}
