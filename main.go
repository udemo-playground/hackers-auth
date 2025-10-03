package main

import (
	"net/http"
	"time"

	_ "github.com/cloudbees-days/hackers-auth/docs" // Import the generated docs

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Hackers Auth API
// @version         1.0
// @description     A simple authentication service for demo purposes
// @host            localhost:8080
// @BasePath        /

// User represents the user model
type User struct {
	Username   string `json:"username" example:"betauser"`
	Password   string `json:"-" example:"betauser"` // "-" tag means this field won't be included in JSON
	Company    string `json:"company" example:"acme global"`
	BetaAccess bool   `json:"beta_access" example:"true"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"betauser"`
	Password string `json:"password" binding:"required" example:"betauser"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  User   `json:"user"`
}

// UserCredentials represents the simplified user response for the /users endpoint
type UserCredentials struct {
	Username string `json:"username" example:"betauser"`
	Password string `json:"password" example:"betauser"`
}

// hardcoded user list
var users = []User{
	{
		Username:   "betauser",
		Password:   "betauser",
		Company:    "acme global",
		BetaAccess: true,
	},
	{
		Username:   "normaluser",
		Password:   "normaluser",
		Company:    "generic co",
		BetaAccess: false,
	},
}

var jwtSecret = []byte("super-secure-password")

func findUser(username, password string) *User {
	for _, user := range users {
		if user.Username == username && user.Password == password {
			return &user
		}
	}
	return nil
}

// @Summary     Login user
// @Description Authenticate user and return JWT token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body LoginRequest true "Login credentials"
// @Success     200 {object} LoginResponse
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Router      /login [post]
func login(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user := findUser(loginReq.Username, loginReq.Password)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":    user.Username,
		"company":     user.Company,
		"beta_access": user.BetaAccess,
		"exp":         time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: tokenString,
		User:  *user,
	})
}

// @Summary     List all users
// @Description Returns a list of all available users with their credentials (for demo purposes only)
// @Tags        users
// @Accept      json
// @Produce     json
// @Success     200 {array} UserCredentials
// @Router      /users [get]
func listUsers(c *gin.Context) {
	credentials := make([]UserCredentials, len(users))
	for i, user := range users {
		credentials[i] = UserCredentials{
			Username: user.Username,
			Password: user.Password,
		}
	}
	c.JSON(http.StatusOK, credentials)
}

func main() {
	r := gin.Default()

	// CORS middleware configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // For demo purposes, allow all origins
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Login endpoint
	r.POST("/login", login)

	// List users endpoint
	r.GET("/users", listUsers)

	r.Run(":8080")
}
