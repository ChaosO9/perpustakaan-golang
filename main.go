package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"perpustakaan-golang/controllers"
	"perpustakaan-golang/db"
	"perpustakaan-golang/forms"

	"github.com/gin-contrib/gzip"
	uuid "github.com/google/uuid"
	"github.com/joho/godotenv"

	docs "perpustakaan-golang/docs"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//CORSMiddleware ...
//CORS (Cross-Origin Resource Sharing)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

//RequestIDMiddleware ...
//Generate a unique ID and attach it to each request for future reference or use
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.New()
		c.Writer.Header().Set("X-Request-Id", uuid.String())
		c.Next()
	}
}

var auth = new(controllers.AuthController)

//TokenAuthMiddleware ...
//JWT Authentication middleware attached to each request that needs to be authenitcated to validate the access_token in the header
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth.TokenValid(c)
		c.Next()
	}
}
// @title Perpustakaan Golang API
// @version 1.0
// @description API for library management system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1
func main() {
	//Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	//Start the default gin server
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	//Custom form validator
	binding.Validator = new(forms.DefaultValidator)

	r.Use(CORSMiddleware())
	r.Use(RequestIDMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	//Start PostgreSQL database
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	db.SetDB(dbConn)

	//Start Redis on database 1 - it's used to store the JWT but you can use it for anythig else
	err = db.ConnectRedis()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }

	v1 := r.Group("api/v1")
	{
		/*** START USER ***/
		user := new(controllers.UserController)
		v1.POST("/user/login", user.Login)
		v1.POST("/user/register", user.Register)
		v1.GET("/user/logout", user.Logout)

		/*** START AUTH ***/
		auth := new(controllers.AuthController)

		v1.POST("/token/refresh", auth.Refresh)

		/*** START Article ***/
		// article := new(controllers.ArticleController)

		// v1.POST("/article", TokenAuthMiddleware(), article.Create)
		// v1.GET("/articles", TokenAuthMiddleware(), article.All)
		// v1.GET("/article/:id", TokenAuthMiddleware(), article.One)
		// v1.PUT("/article/:id", TokenAuthMiddleware(), article.Update)
		// v1.DELETE("/article/:id", TokenAuthMiddleware(), article.Delete)

		// ... other routes ...
		transaction := new(controllers.TransactionController)
		v1.POST("/transactions/borrow", transaction.Borrow)
		v1.PUT("/transactions/return", transaction.Return)
		v1.GET("/user/:id/transactions", transaction.GetTransactionsByMemberID)

		report := new(controllers.ReportController)
		v1.GET("/reports/borrowed_books", report.BorrowedBookReports)
		v1.GET("/reports/returned_books", report.ReturnedBookReports)

		book := new(controllers.BookController)
        v1.POST("/books", book.Create)
        v1.GET("/books", book.All)
        v1.GET("/books/:id", book.One)
        v1.PUT("/books/:id", book.Update)
        v1.DELETE("/books/:id", book.Delete)

		staticFile := r.Group("/uploads")
		{
			// @Summary Get Book Image
			// @Description Retrieves a book image by filename
			// @Tags Static Files
			// @Produce image/*  // Or specify more specific MIME types if needed
			// @Param filename path string true "Image filename"
			// @Success 200 {file} binary "Image file"
			// @Router /uploads/book/{filename} [get]
			staticFile.GET("/book/:filename", func(c *gin.Context) {
				filename := c.Param("filename")
				filePath := "./uploads/book/" + filename 
				
				extension := filepath.Ext(filename) 
		
				contentType := "image/jpeg"
				switch extension {
				case ".png":
					contentType = "image/png"
				case ".gif":
					contentType = "image/gif"
				}
				c.Header("Content-Type", contentType)
				c.File(filePath)
			})
			
			// @Summary Get User Image
			// @Description Retrieves a user image by filename
			// @Tags Static Files
			// @Produce image/*
			// @Param filename path string true "Image filename"
			// @Success 200 {file} binary "Image file"
			// @Router /uploads/user/{filename} [get]
			staticFile.GET("/user/:filename", func(c *gin.Context) {
				filename := c.Param("filename")
				filePath := "./uploads/user/" + filename 
				
				extension := filepath.Ext(filename) 
		
				contentType := "image/jpeg"
				switch extension {
				case ".png":
					contentType = "image/png"
				case ".gif":
					contentType = "image/gif"
				}
				c.Header("Content-Type", contentType)
				c.File(filePath)
			})
		}	
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// r.LoadHTMLGlob("./public/html/*")

	r.Static("/public", "./public")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"ginBoilerplateVersion": "v0.03",
			"goVersion":             runtime.Version(),
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})

    port := os.Getenv("PORT")

    if os.Getenv("ENV") == "LOCAL" {  
        r.Run(":" + port) 
    }
    r.Run(":" + port)

	log.Printf("\n\n PORT: %s \n ENV: %s \n SSL: %s \n Version: %s \n\n", port, os.Getenv("ENV"), os.Getenv("SSL"), os.Getenv("API_VERSION"))

	if os.Getenv("SSL") == "TRUE" {

		//Generated using sh generate-certificate.sh
		SSLKeys := &struct {
			CERT string
			KEY  string
		}{
			CERT: "./cert/myCA.cer",
			KEY:  "./cert/myCA.key",
		}

		r.RunTLS(":"+port, SSLKeys.CERT, SSLKeys.KEY)
	} else {
		r.Run(":" + port)
	}

}
