package main

import (
    "log"
	"net/http"
	"github.com/gin-gonic/gin"
    "fariv/web-service-gin/rdb"
)

type Album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

type successresponse struct {
    Success     int32  `json:"success"`
    Data  []Album  `json:"data"`
}

type successresponsetwo struct {
    Success     int32  `json:"success"`
    Data  Album  `json:"data"`
}

type errorresponse struct {
    Success     int32  `json:"success"`
    Message  string  `json:"message"`
}

var albums = []Album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
    var albums []Album
    dbConn := rdb.DbConnector()
    if dbConn != nil {
        log.Fatalf("Database initialization failed: %v", dbConn)
    }

    rows, err := rdb.DB.Query("SELECT * FROM album")
    defer rows.Close()
    if err != nil {
        errresp := errorresponse{
            Success: 0,
            Message: "Database error happened",
        }
        c.IndentedJSON(http.StatusInternalServerError, errresp)
    } else {
        for rows.Next() {
            var alb Album
            err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price)
            if err != nil {
                log.Fatalf("getAlbums error: %v", err)
            }
            albums = append(albums, alb)
        }
        successresp := successresponse{
            Success: 1,
            Data: albums,
        }
        c.IndentedJSON(http.StatusOK, successresp)
    }
}

func main () {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbum)
	router.GET("/albums/:id", getAlbumByID)
	router.Run("localhost:8080")
}

func postAlbum(c *gin.Context) {

	var newAlbum Album

	err := c.BindJSON(&newAlbum)
	if err != nil {
        errresp := errorresponse{
            Success: 0,
            Message: "Server error",
        }
		c.IndentedJSON(http.StatusInternalServerError, errresp)
	} else {

        dbConn := rdb.DbConnector()
        if dbConn != nil {
            log.Fatalf("Database initialization failed: %v", dbConn)
            errresp := errorresponse{
                Success: 0,
                Message: "Database initialization failed",
            }
            c.IndentedJSON(http.StatusInternalServerError, errresp)
        } else {

            result, err := rdb.DB.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", newAlbum.Title, newAlbum.Artist, newAlbum.Price)
            if err != nil {
                log.Fatalf("Query failed error: %v", err)
                errresp := errorresponse{
                    Success: 0,
                    Message: "Query failed",
                }
                c.IndentedJSON(http.StatusInternalServerError, errresp)
            } else {

                id,_ := result.LastInsertId()
                row := rdb.DB.QueryRow("SELECT * FROM album WHERE id = ?", id)

                var tmpAlbum Album
                err := row.Scan(&tmpAlbum.ID, &tmpAlbum.Title, &tmpAlbum.Artist, &tmpAlbum.Price)
                if err != nil {
                    log.Fatalf("Error: %v", err)
                    errresp := errorresponse{
                        Success: 0,
                        Message: "Scan failed",
                    }
                    c.IndentedJSON(http.StatusInternalServerError, errresp)
                } else {
                    successresp := successresponsetwo{
                        Success: 1,
                        Data: tmpAlbum,
                    }
                    c.IndentedJSON(http.StatusOK, successresp)
                }
            }
        }
    }

}

func getAlbumByID(c *gin.Context) {
    id := c.Param("id")

    for _, a := range albums {
        if a.ID == id {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}