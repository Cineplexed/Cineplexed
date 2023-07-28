package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
	"github.com/gin-contrib/cors"
	_ "cineplexed.com/docs"
	"strconv"
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// moviesByName godoc
// @Summary moviesByName
// @Description Get a list of possible movies by it's title
// @Tags movie
// @Param Item body docs_Title true "Title"
// @Accept */*
// @Produce json
// @Router /getMovieOptions [GET]
func moviesByName(c *gin.Context) {
	checkTime()
	log("INFO", "Getting movie options...")
	title := c.Query("title")
	if title != "" {
		c.JSON(http.StatusOK, getMovieByName(title))
	} else {
		c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Please enter a title"})
	}
}

// movieWithDetails godoc
// @Summary movieWithDetails
// @Description Get a movie with extensive details using it's ID
// @Tags movie
// @Param Item body docs_ID true "ID"
// @Accept */*
// @Produce json
// @Router /getMovieDetails [GET]
func movieWithDetails(c *gin.Context) {
	log("INFO", "Getting movie details...")
	id := c.Query("id")
	numId, err := strconv.Atoi(id)
	if err != nil {	
		log("ERROR", "Cannot get movie details")
		c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Please enter an ID"})
	} else {
		if numId != 0 {
			c.JSON(http.StatusOK, getMovieWithDetail(numId))
		}
	}
}

// Hint godoc
// @Summary Hint
// @Description Get a hint towards the daily movie
// @Tags movie
// @Accept */*
// @Produce json
// @Router /getHint [GET]
func getHint(c *gin.Context) {
	log("INFO", "Getting hint...")
	var entry selections
	result := db.Last(&entry)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		c.JSON(http.StatusBadRequest, Response{Success: false, Context: result.Error.Error()})
	} else {
		c.JSON(http.StatusOK, Hint{Tagline: entry.Tagline, Overview: entry.Overview})
		log("INFO", "Hint given")
	}
}

// makeUser godoc
// @Summary makeUser
// @Description create a new user
// @Tags users
// @Param UserData body Users true "User Data"
// @Accept */*
// @Produce json
// @Router /makeUser [POST]
func makeUser(c *gin.Context) {
	body, err := c.GetRawData()
	var creds Users
	json.Unmarshal(body, &creds)
	if err != nil {
		log("ERROR", "Failed to create user")
		c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Failed to read request body"})
	} else {
		hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
		if err != nil {
			log("ERROR", "Failed to encrypt user data")
			c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Failed to encrypt"})
		} else {
			var user = User{
				ID: uuid.New().String(), 
				Username: creds.Username, 
				Password: string(hash), 
				CreatedAt: time.Now().String(),
				DeletedAt: "",
				UpdatedAt: "",
				SolvedPuzzles: 0,
				FailedPuzzles: 0,
				LastSolvedPuzzle: "",
				Active: true}
			result := db.Create(&user)
			if result.Error != nil {
				log("ERROR", "Failed to post user data to postgres")
				c.JSON(http.StatusBadRequest, Response{Success: false, Context: result.Error.Error()})
			} else {
				log("INFO", "Created user")
				c.JSON(http.StatusOK, Response{Success: true, Context: "Created User"})
			}
		}
	}
}

// validateUser godoc
// @Summary validateUser
// @Description validate a user with a username and password
// @Tags users
// @Param UserData body Users true "User Data"
// @Accept */*
// @Produce json
// @Router /validateUser [POST]
func validateUser(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		log("ERROR", "Error reading response body")
		c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Failed to read request body"})
	} else {
		var creds Users
		json.Unmarshal(body, &creds)

		var validateCreds User
		db.Where("username = ?", creds.Username).First(&validateCreds)
		if validateCreds.Active {
			err := bcrypt.CompareHashAndPassword([]byte(validateCreds.Password), []byte(creds.Password))
			if err != nil {
				fmt.Println(err)
				log("WARNING", "Incorrect credentials on log in")
				c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Please enter a valid username and password"})
			} else {
				var key UserInfo
				db.Table("users").Where("username = ? AND password = ?", creds.Username, validateCreds.Password).First(&key)
				log("INFO", "Validated user")
				c.JSON(http.StatusOK, key)
			}
		} else {
			log("WARNING", "Attempted to sign into inactive account")
			c.JSON(http.StatusOK, Response{Success: false, Context: "Please enter a valid username and password"})
		}
	}
}

// deleteUser godoc
// @Summary deleteUser
// @Description delete a user with a User-Id
// @Tags users
// @Param User-Id header string true "UserID" 
// @Accept */*
// @Produce json
// @Router /deleteUser [DELETE]
func deleteUser(c *gin.Context) {
	header := c.Request.Header["User-Id"]
	if len(header) == 0 {
		log("ERROR", "Invalid User-Id given")
		c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Please enter a User ID"})
	} else {
		var entry User
		result := db.Where("id = ?", string(header[0])).First(&entry)
		if entry.Active {
			if result.Error != nil {
				log("ERROR", "Invalid User-Id given")
				c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Please enter a valid User ID"})
			} else {
				entry.Active = false
				entry.DeletedAt = time.Now().String()
				result := db.Save(&entry)
				if result.Error != nil {
					log("ERROR", "Failed to post deletion to postgres")
					c.JSON(http.StatusBadRequest, Response{Success: false, Context: result.Error.Error()})
				} else {
					log("INFO", "Deleted user")
					c.JSON(http.StatusOK, Response{Success: true, Context: "Deleted User"})
				}
			}
		} else {
			log("ERROR", "Attempted to delete an inactive account")
			c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Please enter a valid user ID"})
		}
	}
}

// updateUser godoc
// @Summary updateUser
// @Description update a user with a User-Id and new username and password
// @Tags users
// @Param User-Id header string true "UserID" 
// @Param UserData body Users true "User Data"
// @Accept */*
// @Produce json
// @Router /updateUser [PATCH]
func updateUser(c *gin.Context) {
	header := c.Request.Header["User-Id"]
	if len(header) == 0 {
		log("ERROR", "Invalid User-Id given")
		c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Please enter a valid user ID"})
	} else {
		body, err := c.GetRawData()
		if err != nil {
			log("ERROR", "Failed to read request body")
			c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Failed to read request body"})
		} else {
			var entry User
			result := db.Where("id = ?", string(header[0])).First(&entry)
			if entry.Active {
				if result.Error != nil {
					log("ERROR", "Invalid User-Id given")
					c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Please enter a valid user ID"})
				} else {
					var creds Users
					json.Unmarshal(body, &creds)
					fmt.Println(creds)
					hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
					if err != nil {
						log("ERROR", "Failed to encrypt")
						c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Failed to encrypt"})
					} else {
						entry.Username = creds.Username
						entry.Password = string(hash)
						entry.UpdatedAt = time.Now().String()
						result := db.Save(&entry)
						if result.Error != nil {
							log("ERROR", "Could not post user info to postgres")
							c.JSON(http.StatusBadRequest, Response{Success: false, Context: result.Error.Error()})
						} else {
							log("INFO", "Updated user info")
							c.JSON(http.StatusOK, Response{Success: true, Context: "Updated info"})
						}
					}
				}
			} else {
				log("WARNING", "Attempted to update inactive account")
				c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Please enter a valid user ID"})
			}

		}
	}
}

// finishGame godoc
// @Summary finishGame
// @Description update success rates of user and daily
// @Tags users
// @Param User-Id header string false "UserID" 
// @Param UserData body Users true "User Data"
// @Accept */*
// @Produce json
// @Router /finishGame [POST]
func finishGame(c *gin.Context) {
	userFailed := false
	//Get data from body
	body, err := c.GetRawData()
	if err != nil {
		log("ERROR", "Failed to read request body")
		c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Failed to read request body"})
	} else {
		//Covert body data to boolean 
		var status GameStatus
		json.Unmarshal(body, &status)

		//Check if User-Id is present
		header := c.Request.Header["User-Id"]
		if len(header) != 0 {
			if len(header[0]) == 0 {
				log("ERROR", "Invalid user-id entered")
				c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Please enter a valid user ID"})
				userFailed = true
			} else {
				//Update user field to reflect win/loss
				var entry User
				result := db.Where("id = ?", string(header[0])).First(&entry)
				if entry.Active {
					if result.Error != nil {
						log("ERROR", "Invalid user ID entered")
						c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Please enter a valid user ID"})
						userFailed = true
					} else {
						if status.Won {
							entry.SolvedPuzzles = entry.SolvedPuzzles + 1
							entry.LastSolvedPuzzle = time.Now().String()
						} else {
							entry.FailedPuzzles = entry.FailedPuzzles + 1
						}
						result := db.Save(&entry)
						if result.Error != nil {
							log("ERROR", "Failed to post info to postgres")
							c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Failed to post info to postgres"})
							userFailed = true
						} else {
							log("INFO", "Submitted daily win info")
						}
					}
				} else {
					log("ERROR", "Attempted to complete a puzzle on inactive account")
					c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Please enter a valid user ID"})
					userFailed = true
				}
			}
		}
		if !userFailed {
			//User-Id not present
			var daily selections
			result := db.Last(&daily)
			if result.Error != nil {
				log("ERROR", "Failed to pull info from postgres")
				c.JSON(http.StatusBadRequest, Response{Success: false, Context: result.Error.Error()})
			} else {
				//Update daily movie field to reflect win/loss
				if status.Won {
					daily.NumCorrect = daily.NumCorrect + 1
				} else {
					daily.NumIncorrect = daily.NumIncorrect + 1
				}
				result := db.Model(&daily).Where("date = ?", daily.Date).Updates(map[string]interface{} {
					"NumCorrect": daily.NumCorrect,
					"NumIncorrect": daily.NumIncorrect, 
				})
				if result.Error != nil {
					log("ERROR", "Failed to post info to postgres")
					c.JSON(http.StatusBadRequest, Response{Success: false, Context: "Failed to post info to postgres"})
				} else {
					log("INFO", "Submitted daily info")
					c.JSON(http.StatusOK, Response{Success: true, Context: "Submitted info"})
				}
			}
		}
	}
}

func getHost() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("ERROR")
		return ""
	} else {
		return os.Getenv("host")
	}
}

// @title Gin Swagger Cineplexed
// @version 1.0
// @description This is Cineplexed.
// @host localhost:5050
// @BasePath /
// @schemes http
func main() {
	getEnv()
	getTargetTime()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "User-Id"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}

	router.Use(cors.New(config))

	router.GET("/getMovieOptions", moviesByName)
	router.GET("/getMovieDetails", movieWithDetails)
	router.GET("/getHint", getHint)

	router.POST("/makeUser", makeUser)
	router.POST("/validateUser", validateUser)
	router.DELETE("/deleteUser", deleteUser)
	router.PATCH("/updateUser", updateUser)
	router.POST("/finishGame", finishGame)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(getHost())
}