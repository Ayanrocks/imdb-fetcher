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
	//
	//resp, err := http.Get(u.String())
	//if err != nil {
	//	return errors.New("error sending request")
	//}
	//defer resp.Body.Close()
	//_, err = ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return errors.New("error parsing request body")
	//}

	c := colly.NewCollector()
	c.OnHTML("tbody.lister-list tr td.titleColumn a", func(element *colly.HTMLElement) {
		log.Println(element.Text)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})


	c.Visit(u.String())

	return nil
}
