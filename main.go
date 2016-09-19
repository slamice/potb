package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"potb-server/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"potb-server/Godeps/_workspace/src/github.com/go-gorp/gorp"
	_ "potb-server/Godeps/_workspace/src/github.com/lib/pq"
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
	ID            int64 `db:"performer_id"`
	Created       int64
	English_name  string `json:"english_name" form:"english_name" binding:"required"`
	Japanese_name string `json:"japanese_name" form:"japanese_name" binding:"required"`
}

// Game class
type Game struct {
	ID                   int64 `db:"game_id"`
	Created              int64
	English_name         string
	Japanese_name        string
	English_description  string
	Japanese_description string
}

// News class
type News struct {
	ID          int64 `db:"news_id"`
	Created     int64
	Name        string
	Description string
}

// News class
type Program struct {
	ID      int64 `db:"program_id"`
	Created     int64
	ProgramDate string
}

// UpdateProgramDate replaces date for the program
func UpdateProgramDate(programDate string) Program {
	_, err := dbmap.Exec(fmt.Sprintf("update program set programddate = '%s'", programDate))
	checkErr(err, "Update failed for adding a new program date")

	return GetProgram()
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
	checkErr(err, "Insert failed for Team")
	return team
}

// AddScore adds the score by 1
func AddScore(name string) {
	_, err := dbmap.Exec("update team set score = score + 1 where name = $1", name)
	checkErr(err, "Update failed for updating a team score")
}

func AddNews(news []News) {
	for i := range news {
		newsItem := News{
			Created:     time.Now().UnixNano(),
			Name:        news[i].Name,
			Description: news[i].Description,
		}

		err := dbmap.Insert(&newsItem)
		checkErr(err, "Creating new news failed")
	}
}

// AddGames
func AddGames(games []Game) {
	for i := range games {
		game := Game{
			Created:              time.Now().UnixNano(),
			English_name:         games[i].English_name,
			Japanese_name:        games[i].Japanese_name,
			English_description:  games[i].English_description,
			Japanese_description: games[i].Japanese_description,
		}

		err := dbmap.Insert(&game)
		checkErr(err, "Creating a new performer failed")
	}
}

// AddPerformers
func AddPerformers(performers []Performer) {
	for i := range performers {
		perf := Performer{
			Created:       time.Now().UnixNano(),
			English_name:  performers[i].English_name,
			Japanese_name: performers[i].Japanese_name,
		}

		err := dbmap.Insert(&perf)
		checkErr(err, "Creating a new performer failed")
	}
}

// GetTeams fetches all teams
func GetTeams() []Team {
	var teams []Team
	_, err := dbmap.Select(&teams, "select * from team order by score")
	checkErr(err, "Team select failed")
	return teams
}

// GetPerformers fetches all performers
func GetPerformers() []Performer {
	var performers []Performer
	_, err := dbmap.Select(&performers, "select * from performer")
	checkErr(err, "Performer select failed")
	return performers
}

// GetProgram fetchs the current program date
func GetProgram() Program {
	var program Program
	err := dbmap.SelectOne(&program, "select * from program limit 1")
	checkErr(err, "Program Date select failed")
	return program
}

// GetGames fetches all games
func GetGames() []Game {
	var games []Game
	_, err := dbmap.Select(&games, "select * from game")
	checkErr(err, "Game select failed")
	return games
}

// GetNews fetches all news
func GetNews() []News {
	var news []News
	_, err := dbmap.Select(&news, "select * from news")
	checkErr(err, "News select failed")
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

// ProgramGet
func ProgramGet(c *gin.Context) {
	content := gin.H{}
	content["data"] = GetProgram()
	c.JSON(200, content)
}

// NewsGet
func NewsGet(c *gin.Context) {
	content := gin.H{}
	content["data"] = GetNews()
	c.JSON(200, content)
}

// GamesGet
func GamesGet(c *gin.Context) {
	content := gin.H{}
	content["data"] = GetGames()
	c.JSON(200, content)
}

// PerformersGet
func PerformersGet(c *gin.Context) {
	content := gin.H{}
	content["data"] = GetPerformers()
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

// Clear Programs
func ClearProgram() {
	deleteAllData("program")
}

// Delete all data
func deleteAllData(table string) {
	_, err := dbmap.Exec("DELETE FROM " + table)
	checkErr(err, fmt.Sprintf("Clear table failed for table %s", table))
}

// ClearScores clears the scores to 0
func ClearScores() {
	_, err := dbmap.Exec("update team set score = 0")
	checkErr(err, "Update failed for scores")
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

// ProgramPost
func ProgramPost(c *gin.Context) {
	var program Program
	err := c.Bind(&program)
	if err != nil {
		panic(err)
	}
	UpdateProgramDate(program.ProgramDate)
	content := gin.H{"result": "Success"}
	c.JSON(200, content)
}

// PerformersPost
func PerformersPost(c *gin.Context) {
	ClearPerformers()
	var performers []Performer
	err := c.Bind(&performers)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	AddGames(games)
	content := gin.H{"result": "Success"}
	c.JSON(200, content)
}

// NewsPost
func NewsPost(c *gin.Context) {
	ClearNews()
	var news []News
	err := c.Bind(&news)
	if err != nil {
		panic(err)
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

	username := os.Getenv("PIRATES_DB_USERNAME")
	password := os.Getenv("PIRATES_DB_PASSWORD")
	host := os.Getenv("PIRATES_DB_HOST")
	dbName := os.Getenv("PIRATES_DB")

	connection_str := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", username, password, host, dbName)
	db, err := sql.Open("postgres", connection_str)
	checkErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(Team{}, "team").SetKeys(true, "ID")
	dbmap.AddTableWithName(Performer{}, "performer").SetKeys(true, "ID")
	dbmap.AddTableWithName(Game{}, "game").SetKeys(true, "ID")
	dbmap.AddTableWithName(News{}, "news").SetKeys(true, "ID")
	dbmap.AddTableWithName(Program{}, "program").SetKeys(true, "ID")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	// Create default count for teams
	count, err := dbmap.SelectInt("select count(*) from team")
	checkErr(err, "Create tables failed")
	if count == 0 {
		team := Team{
			Created: time.Now().UnixNano(),
			Name:    "red",
			Color:   "#DF3535",
			Score:   0,
		}
		err = dbmap.Insert(&team)
		checkErr(err, "insert failed")

		team = Team{
			Created: time.Now().UnixNano(),
			Name:    "white",
			Color:   "#E6E7E9",
			Score:   0,
		}
		err = dbmap.Insert(&team)
		checkErr(err, "insert failed")
	}

	// Create default program date
	count, err = dbmap.SelectInt("select count(*) from program")
	checkErr(err, "Create tables failed")
	if count == 0 {
		program := Program{
			Created: 	 time.Now().UnixNano(),
			ProgramDate: "2016-01-01T01:01:00+09:00",
		}
		err = dbmap.Insert(&program)
		checkErr(err, "Init insert failed for Program Date")
	}

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func main() {
	router := gin.Default()

	// Statics
	router.Static("/assets", "./assets")

	router.PUT("/addscore", ScorePut)
	router.POST("/clearscores", ScorePost)
	router.GET("/getteams", TeamsGet)

	router.POST("/addperformers", PerformersPost)
	router.GET("/getperformers", PerformersGet)

	router.POST("/addgames", GamesPost)
	router.GET("/getgames", GamesGet)

	router.POST("/addnews", NewsPost)
	router.GET("/getnews", NewsGet)

	router.POST("/addprogramdate", ProgramPost)
	router.GET("/getprogramdate", ProgramGet)

	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Pirates of Tokyo Bay!!!",
		})
	})
	router.GET("/programs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "programs.tmpl", gin.H{
			"title": "Programs!!!",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

    s := &http.Server{
        Addr:           ":" + port,
        Handler:        router,
        ReadTimeout:    30 * time.Second,
        WriteTimeout:   30 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    
	if err := s.ListenAndServe(); err != nil {
	    log.Fatal("ListenAndServe: ", err)
	}

}
