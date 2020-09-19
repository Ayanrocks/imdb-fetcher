package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
	"log"
	"net/url"
	"strconv"
	"strings"
)

var (
	cacheDir string
)

type MovieDetails struct {
	Title              string `json:"title"`
	MovieReleasingYear string `json:"movie_releasing_year"`
	IMDBRating         string `json:"imdb_rating"`
	Summary            string `json:"summary"`
	Duration           string `json:"duration"`
	Genre              string `json:"genre"`
}

var root = &cobra.Command{
	Use:     "imdb [link] [limit count]",
	Example: "  imdb https://www.imdb.com/india/top-rated-indian-movies/ 2",
	Short:   "Imdb is a movie fetcher",
	Long:    "Imdb is a command line too16Ml to fetch movie details based on the imdb list supplied with count limit",
	RunE:    Fetch,
}

func Execute() error {
	return root.Execute()
}

func init() {
	root.PersistentFlags().StringVarP(&cacheDir, "cache", "c", "", "Specify cache dir to reuse the link")
}

func Fetch(_ *cobra.Command, args []string) error {
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
	ch := make(chan []string)
	go movieNameScrapper(ch, u.String())
	movieIDs := <-ch
	movieIDs = movieIDs[:count]
	var movieDetailsArr []MovieDetails
	movieCh := make(chan MovieDetails)
	for _, v := range movieIDs {
		go getMovieDetails(v, movieCh)
	}
	for range movieIDs {
		detail := <-movieCh
		movieDetailsArr = append(movieDetailsArr, detail)
	}
	jsonDetail, err := json.Marshal(&movieDetailsArr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonDetail))
	return nil
}

func movieNameScrapper(ch chan<- []string, u string) {

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(1),
		colly.CacheDir(cacheDir),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1})
	var movieIDs []string
	c.OnHTML("tbody.lister-list tr td.ratingColumn div.seen-widget", func(element *colly.HTMLElement) {
		movieIDs = append(movieIDs, element.Attr("data-titleid"))
	})
	c.OnScraped(func(r *colly.Response) {
		ch <- movieIDs
	})
	c.Visit(u)
}

func getMovieDetails(id string, ch chan<- MovieDetails) {

	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(1),
		colly.CacheDir(cacheDir),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 1})

	movieDetail := MovieDetails{}

	c.OnHTML("div#title-overview-widget", func(element *colly.HTMLElement) {

		element.ForEach("div.vital > div.title_block > div.title_bar_wrapper", func(_ int, chEl *colly.HTMLElement) {
			movieDetail.IMDBRating = chEl.ChildText("div.ratings_wrapper > div.imdbRating > div.ratingValue > strong > span")
			movieTitle := strings.Split(chEl.ChildText("div.titleBar > div.title_wrapper > h1"), " ")
			m := strings.Split(movieTitle[len(movieTitle)-1], "(")
			movieTitle = movieTitle[:len(movieTitle)-1]
			movieTitle = append(movieTitle, m[0])
			movieDetail.Title = strings.Join(movieTitle, " ")
			movieDetail.MovieReleasingYear = chEl.ChildText("div.titleBar > div.title_wrapper > h1 > span")
			movieDetail.Duration = chEl.ChildText("div.titleBar > div.title_wrapper > div.subtext > time")
			movieDetail.Genre = chEl.ChildText("div.titleBar > div.title_wrapper > div.subtext > a:first-of-type")
		})

		element.ForEach("div.plot_summary_wrapper", func(_ int, chEl *colly.HTMLElement) {
			movieDetail.Summary = chEl.ChildText("div.plot_summary  > div.summary_text")
		})
	})

	c.OnScraped(func(r *colly.Response) {
		ch <- movieDetail
	})

	c.Visit(`https://www.imdb.com/title/` + id)
}

// TODO NOW APPEND ALL THE MOVE NAMES TO AN ARRAY AND THEN
// TODO USE we scraping in imdb details page TO REQUEST EACH MOVIE DETAILS
