package main

import (
	"awesomeProject/scrapper"
	"github.com/labstack/echo"
	"os"
	"strings"
)

const fileName string = "jobs.csv"

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	defer os.Remove(fileName)
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	exportFilename := term + " jobs.csv"
	scrapper.Scrape(term)
	return c.Attachment(fileName, exportFilename)
}

func main(){
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}