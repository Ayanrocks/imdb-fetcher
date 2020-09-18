package commands

import (
	"errors"
	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
	"log"
	"net/url"
	"strconv"
)

var (
	imdbUrl    string
	itemsCount int
)

var root = &cobra.Command{
	Use:   "imdb",
	Short: "Imdb is a movie fetcher",
	Long:  "Imdb is a command line tool to fetch movies based on the list supplied ",
	RunE:  Fetch,
}

func Execute() error {
	return root.Execute()
}

func init() {
	root.PersistentFlags().StringVar(&imdbUrl, "url", "", "The imdb movie list supplied")
	root.PersistentFlags().IntVar(&itemsCount, "count", 0, "No. of movies to return from top")

}

func Fetch(cmd *cobra.Command, args []string) error {
	if len(args) < 2 || len(args) > 2 {
		return errors.New("only two args must be specified")
	}
	_, err := strconv.Atoi(args[0])
	if err == nil {
		return errors.New("first argument must be an url")
	}
	count, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.New("second argument must be a number")
	}
	u, err := url.ParseRequestURI(args[0])
	if err != nil {
		return errors.New("an url must be supplied as 1st argument")
	}
	log.Println(count, u)
	ch := make(chan []string)
	go movieNameScrapper(ch, u.String())
	movieIDs := <-ch
	movieIDs = movieIDs[:count]
	log.Println(movieIDs)

	for _, v := range movieIDs {
		log.Println(v)
		//getMovieDetails(v)
	}
	getMovieDetails("w")
	return nil
}

func movieNameScrapper(ch chan<- []string, u string) {

	c := colly.NewCollector()
	var movieIDs []string
	//c.OnHTML("tbody.lister-list tr td.titleColumn a", func(element *colly.HTMLElement) {
	//	movieIDs = append(movieIDs, element.Text)
	//	element.Attr("[]")
	//})

	c.OnHTML("tbody.lister-list tr td.ratingColumn div.seen-widget", func(element *colly.HTMLElement) {
		movieIDs = append(movieIDs,element.Attr("data-titleid"))
	})

	c.OnScraped(func(r *colly.Response) {
		log.Println("Ended Scraped")
		ch <- movieIDs
	})

	c.Visit(u)
}

func getMovieDetails(id string) {

	c := colly.NewCollector()

	c.OnHTML("div#title-overview-widget", func(element *colly.HTMLElement) {
		log.Println(element.DOM.Nodes)
	})

	c.OnScraped(func(r *colly.Response) {
		log.Println("Ended Scraped In Details")
	})

	c.Visit(`https://www.imdb.com/title/tt0048473`)
}

// TODO NOW APPEND ALL THE MOVE NAMES TO AN ARRAY AND THEN
// TODO USE we scraping in imdb details page TO REQUEST EACH MOVIE DETAILS
