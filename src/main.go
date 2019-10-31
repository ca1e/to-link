package main

import (
    "encoding/base64"
    "net/http"
    "runtime"
    "regexp"
    "tolnk/config"
    "tolnk/murshort"
    "tolnk/kvstore"
    "github.com/gin-gonic/gin"
)

var LOCAL_URL = config.Conf.LocalUrl

const (
    RESULT_OK = "OK"
    RESULT_ERROR = "error"
)

type MSG struct {
    TinyUrl  string `json:"tinyurl"`
    LongUrl  string `json:"longurl"`
}

func long2short(lng string) string {
    tinyUrl := murshort.Mur3h62(lng)
    return tinyUrl
}

func checkUrl(url string) error {
    _, err := http.Get(url)
    return err
}

func create(c *gin.Context){
    longUrl_base64 := c.Query("url")
    alias := c.Query("alias")

    longUrl_bytes, err := base64.StdEncoding.DecodeString(longUrl_base64)
    if(err != nil){
        c.JSON(http.StatusOK, gin.H{
            "result": RESULT_ERROR,
            "data": "Decode error",
        })
        return
    }

    longUrl := string(longUrl_bytes)

    err = checkUrl(longUrl)
    if(err != nil){
        c.JSON(http.StatusOK, gin.H{
            "result": RESULT_ERROR,
            "data": "Invalid url",
        })
        return
    }

    if(alias != "") {
        _, ok := kvstore.Exist(alias)
        if(ok) {
            c.JSON(http.StatusOK, gin.H{
                "result": RESULT_ERROR,
                "data": "Alias existed",
            })
            return
        }
    }else{
        alias = long2short(longUrl)
    }

    match, _ := regexp.MatchString("([A-Za-z0-9_]+)", alias)
    if( len(alias) < 3 || len(alias) > 16 && !match ){
        c.JSON(http.StatusOK, gin.H{
            "result": RESULT_ERROR,
            "data": "Invalid alias",
        })
        return
    }

    tinyUrl := alias

    err = kvstore.Store(tinyUrl, longUrl)
    if(err != nil){
        c.JSON(http.StatusOK, gin.H{
            "result": RESULT_ERROR,
            "data": "Limited",
        })
        return
    }

    data := &MSG{
        TinyUrl: LOCAL_URL + tinyUrl,
        LongUrl: longUrl,
    }
    c.JSON(http.StatusOK, gin.H{
        "result": RESULT_OK,
        "data": data,
    })
}

func redirect(c *gin.Context){
    guid := c.Param("guid")
    longUrl, ok := kvstore.Exist(guid)
    if(ok){
        println("go to:", longUrl)
        c.Redirect(http.StatusFound, longUrl) //301 http.StatusMovedPermanently
    }else {
        c.String(http.StatusNotFound, "")
    }
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU());
    // gin.SetMode(gin.ReleaseMode)

    r := gin.Default()

    // r.Static("/static", "../html/static")
    r.StaticFile("/", "../html/index.html")
    r.GET("/create.php", create)
    r.GET("/t/:guid", redirect)

    r.Run()
}
