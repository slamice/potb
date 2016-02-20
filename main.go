package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
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

// Performers class
type Performer struct {
	ID      int64 `db:"performer_id"`
	Created int64
	English_name    string `json:"english_name" form:"english_name" binding:"required"`
	Japanese_name   string `json:"japanese_name" form:"japanese_name" binding:"required"`
}

// Game class
type Game struct {
	ID      int64 `db:"game_id"`
	Created int64
	English_name    string
	Japanese_name   string
	English_description    string
	Japanese_description   string
}

// News class
type News struct {
	ID      int64 `db:"news_id"`
	Created int64
	Name string
	Description    string
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
	_, err := dbmap.Exec("update team set score = score + 1 where name = ?", name)
	checkErr(err, "Update failed")
}

func AddNews(news []News) {
	for i := range news {	
		newsItem := News{
			Created: time.Now().UnixNano(),
			Name:   news[i].Name,
			Description:  news[i].Description,
		}

		err := dbmap.Insert(&newsItem)
		checkErr(err, "Creating new news failed")
	}
}

// AddGames
func AddGames(games []Game) {
	for i := range games {	
		game := Game{
			Created: time.Now().UnixNano(),
			English_name:   games[i].English_name,
			Japanese_name:  games[i].Japanese_name,
			English_description:    games[i].English_description,
			Japanese_description:   games[i].Japanese_description,
		}

		err := dbmap.Insert(&game)
		checkErr(err, "Creating a new performer failed")
	}
}

// AddPerformers 
func AddPerformers(performers []Performer) {
	for i := range performers {	
		perf := Performer{
			Created: time.Now().UnixNano(),
			English_name:   performers[i].English_name,
			Japanese_name:  performers[i].Japanese_name,
		}

		err := dbmap.Insert(&perf)
		checkErr(err, "Creating a new performer failed")
	}
}

// GetTeams fetches all teams
func GetTeams() []Team {
	var teams []Team
	_, err := dbmap.Select(&teams, "select * from team order by score")
	checkErr(err, "Select failed")
	return teams
}

// GetPerformers fetches all performers
func GetPerformers() []Performer {
	var performers []Performer
	_, err := dbmap.Select(&performers, "select * from performer")
	checkErr(err, "Select failed")
	return performers
}

// GetGames fetches all games
func GetGames() []Game {
	var games []Game
	_, err := dbmap.Select(&games, "select * from game")
	checkErr(err, "Select failed")
	return games
}

// GetNews fetches all news
func GetNews() []News {
	var news []News
	_, err := dbmap.Select(&news, "select * from news")
	checkErr(err, "Select failed")
	return news
}

// TeamsGet gets all teh teams in the db
func TeamsGet(c *gin.Context) {
	content := gin.H{}
	for k, v := range GetTeams() {
		content[strconv.Itoa(k)] = gin.H{
			"Name":  v.Name,
			"Color": v.Color,
			"Score": v.Score,
		}
	}
	c.JSON(200, content)
}

// NewsGet 
func NewsGet(c *gin.Context) {
	content := gin.H{}
	for k, v := range GetNews() {
		content[strconv.Itoa(k)] = gin.H{
			"name": v.Name,
			"description": v.Description,
		}
	}
	c.JSON(200, content)
}

// GamesGet 
func GamesGet(c *gin.Context) {
	content := gin.H{}
	for k, v := range GetGames() {
		content[strconv.Itoa(k)] = gin.H{
			"english_name": v.English_name,
			"japanese_name": v.Japanese_name,
		}
	}
	c.JSON(200, content)
}


// PerformersGet 
func PerformersGet(c *gin.Context) {
	content := gin.H{}
	for k, v := range GetPerformers() {
		content[strconv.Itoa(k)] = gin.H{
			"english_name": v.English_name,
			"japanese_name": v.Japanese_name,
		}
	}
	c.JSON(200, content)
}


// ClearPerformers clears the performers to 0
func ClearPerformers() {
	deleteAllData("performer")
}

// Clear Games
func ClearGames() {
	deleteAllData("game")
}

// Clear News
func ClearNews() {
	deleteAllData("news")
}

// Delete all data
func deleteAllData(table string) {
	_, err := dbmap.Exec("DELETE FROM " + table)
	checkErr(err, "Clear table failed")
}



// ClearScores clears the scores to 0
func ClearScores() {
	_, err := dbmap.Exec("update team set score = 0")
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

// PerformersPost 
func PerformersPost(c *gin.Context) {
	ClearPerformers()
	var performers []Performer
	err := c.Bind(&performers)
	if err != nil {
        panic (err)
    }
	AddPerformers(performers)
	content := gin.H{"result": "Success"}
	c.JSON(200, content)
}

// GamesPost
func GamesPost(c *gin.Context) {
	ClearGames()
	var games []Game
	err := c.Bind(&games)
	if err != nil {
        panic (err)
    }
	AddGames(games)
	content := gin.H{"result": "Success"}
	c.JSON(200, content)
}

// GamesPost
func NewsPost(c *gin.Context) {
	ClearNews()
	var news []News
	err := c.Bind(&news)
	if err != nil {
        panic (err)
    }
	AddNews(news)
	content := gin.H{"result": "Success"}
	c.JSON(200, content)
}

// ScorePost adds to the score
func ScorePost(c *gin.Context) {
	ClearScores()
	content := gin.H{"result": "Success"}
	c.JSON(200, content)
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(Team{}, "team").SetKeys(true, "ID")
	dbmap.AddTableWithName(Performer{}, "performer").SetKeys(true, "ID")
	dbmap.AddTableWithName(Game{}, "game").SetKeys(true, "ID")
	dbmap.AddTableWithName(News{}, "news").SetKeys(true, "ID")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	count, err := dbmap.SelectInt("select count(*) from team")
	checkErr(err, "Create tables failed")
	if count == 0 {
		team := Team{
			Created: time.Now().UnixNano(),
			Name:    "Red",
			Color:   "#DF3535",
			Score:   0,
		}
		err = dbmap.Insert(&team)
		checkErr(err, "insert failed")

		team = Team{
			Created: time.Now().UnixNano(),
			Name:    "White",
			Color:   "#E6E7E9",
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

	// Statics
	app.Static("/assets", "./assets")


	app.PUT("/addscore", ScorePut)
	app.POST("/clearscores", ScorePost)
	app.GET("/getteams", TeamsGet)

	app.POST("/addperformers", PerformersPost)
	app.GET("/getperformers", PerformersGet)

	app.POST("/addgames", GamesPost)
	app.GET("/getgames", GamesGet)

	app.POST("/addnews", NewsPost)
	app.GET("/getnews", NewsGet)

	app.LoadHTMLGlob("templates/*")
	app.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Pirates of Tokyo Bay!!!",
		})
	})
	app.GET("/programs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "programs.tmpl", gin.H{
			"title": "Programs!!!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	app.Run(":" + port)
}
