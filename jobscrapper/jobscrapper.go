package jobscrapper

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type extractedJob struct {
	id string
	title string
	location string
	salary string
	summary string
}

const baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main(){
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages()

	for i := 0; i < totalPages; i++{
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++{
		extractedJob := <- c
		jobs = append(jobs, extractedJob...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

func getPage(page int, mainC chan <- []extractedJob){
	var jobs []extractedJob
	c := make(chan extractedJob)

	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")

	searchCards.Each(func(i int, card *goquery.Selection){
		go extractJob(card, c)
	})

	for i:=0; i < searchCards.Length(); i++{
		job := <-c
		jobs = append(jobs, job)
	}
	mainC <- jobs
}
func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".title > a").Text())
	location := cleanString(card.Find(".sjcl").Text())
	salary := cleanString(card.Find(".salaryText").Text())
	summary := cleanString(card.Find(".summary").Text())
	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary,
	}
}

func cleanString(str string) string{
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection){
		pages = s.Find("a").Length()
	})
	return pages
}

func writeJobs(jobs []extractedJob){
	const filename = "jobs.csv"
	file, err := os.Create(filename)
	checkErr(err)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	c := make(chan []string)

	for _, job := range jobs {
		go makeJobData(job, c)
	}

	for i:=0; i < len(jobs); i++{
		jobData := <-c
		wErr := w.Write(jobData)
		checkErr(wErr)
	}
}

func makeJobData(job extractedJob, c chan<- []string) {
	var jobDetail = "https://kr.indeed.com/viewjob?jk=" + job.id
	c <- []string{jobDetail, job.title, job.location, job.salary, job.summary}
}


func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response){
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status : ", res.StatusCode)
	}
}

