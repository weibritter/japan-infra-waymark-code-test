package main

import (
	"os"
	"fmt"
	"time"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)


func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"cache-control", "x-requested-with"},
		AllowAllOrigins: true,
		MaxAge: 12 * time.Hour,
	}))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20  // 8 MiB
	r.POST("/upload", func(c *gin.Context) {
		// single file
		file, err := c.FormFile("qqfile")
		
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		filename := c.PostForm("qqfilename")
		
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		AddFileToS3(filename)
		
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func AddFileToS3(filename string) {
	bucket := "forcodetest"
	
	file, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Unable to open file %q, %v", err)
    }

    defer file.Close()
	
	sess, err := session.NewSession(&aws.Config{
        Region: aws.String("ap-northeast-1")},
    )
	
	uploader := s3manager.NewUploader(sess)
	
	_, err = uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucket),
        Key: aws.String(filename),
        Body: file,
    })
    if err != nil {
        fmt.Printf("Unable to upload %q to %q, %v", filename, bucket, err)
		return
    }

    fmt.Printf("Successfully uploaded %q to %q\n", filename, bucket)	
}
