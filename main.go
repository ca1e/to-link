package main

import (
	"embed"
	"encoding/base64"
	"net/http"
	"regexp"
	"runtime"

	"github.com/ca1e/to-link/internal/config"
	"github.com/ca1e/to-link/internal/kvstore"
	"github.com/ca1e/to-link/internal/murshort"
	"github.com/gin-gonic/gin"
)

//go:embed static/index.html
var IndexHTML embed.FS

var LOCAL_URL = config.Conf.LocalUrl

const (
	RESULT_OK    = "OK"
	RESULT_ERROR = "error"
)

type MSG struct {
	TinyUrl string `json:"tinyurl"`
	LongUrl string `json:"longurl"`
}

func long2short(lng string) string {
	tinyUrl := murshort.Mur3h62(lng)
	return tinyUrl
}

func checkUrl(url string) error {
	_, err := http.Get(url)
	return err
}

func create(c *gin.Context) {
	longUrl_base64 := c.Query("url")
	alias := c.Query("alias")

	longUrl_bytes, err := base64.StdEncoding.DecodeString(longUrl_base64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": RESULT_ERROR,
			"data":   "Decode error",
		})
		return
	}

	longUrl := string(longUrl_bytes)

	err = checkUrl(longUrl)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": RESULT_ERROR,
			"data":   "Invalid url",
		})
		return
	}

	if alias != "" {
		_, ok := kvstore.Exist(alias)
		if ok {
			c.JSON(http.StatusOK, gin.H{
				"result": RESULT_ERROR,
				"data":   "Alias existed",
			})
			return
		}
	} else {
		alias = long2short(longUrl)
	}

	match, _ := regexp.MatchString("([A-Za-z0-9_]+)", alias)
	if len(alias) < 3 || len(alias) > 16 && !match {
		c.JSON(http.StatusOK, gin.H{
			"result": RESULT_ERROR,
			"data":   "Invalid alias",
		})
		return
	}

	tinyUrl := alias

	err = kvstore.Store(tinyUrl, longUrl)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": RESULT_ERROR,
			"data":   "Limited",
		})
		return
	}

	data := &MSG{
		TinyUrl: LOCAL_URL + tinyUrl,
		LongUrl: longUrl,
	}
	c.JSON(http.StatusOK, gin.H{
		"result": RESULT_OK,
		"data":   data,
	})
}

func redirect(c *gin.Context) {
	guid := c.Param("guid")
	longUrl, ok := kvstore.Exist(guid)
	if ok {
		println("go to:", longUrl)
		c.Redirect(http.StatusFound, longUrl) //301 http.StatusMovedPermanently
	} else {
		c.String(http.StatusNotFound, "")
	}
}

func CORSMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(CORSMiddleware)

	// r.StaticFS("/", http.FS(IndexHTML))
	r.StaticFile("/", "static/index.html")
	r.GET("/create.php", create)
	r.GET("/t/:guid", redirect)

	r.Run()
}
