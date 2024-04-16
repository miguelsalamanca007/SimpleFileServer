package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"io"

	"github.com/gin-gonic/gin"
)

func listDir(dirname string) ([]string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, file := range files {
		names = append(names,"📁 " + file.Name())
	}
	return names, nil
}

func ServeFiles(c *gin.Context) {
	names, err := listDir("/fileserver")

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Internal server error")
	}

	c.JSON(http.StatusOK, names)
}

func fileHandler(c *gin.Context) {
    // Open the file
    file, err := os.Open("/fileserver/"+c.Param("filename"))
    if err != nil {
        c.String(http.StatusNotFound, "File not found")
        return
    }
    defer file.Close()

    // Set the appropriate headers
    c.Writer.Header().Set("Content-Disposition", "attachment; filename="+c.Param("filename"))
    c.Writer.Header().Set("Content-Type", "application/octet-stream")

    // Copy the file to the response writer
    _, err = io.Copy(c.Writer, file)
    if err != nil {
        c.String(http.StatusInternalServerError, "Internal server error")
        return
    }
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.GET("/file", ServeFiles)
	r.GET("/download/:filename", fileHandler)
	http.ListenAndServe("localhost:8080", r)
}