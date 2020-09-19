# IMDB Movie Detail Fetcher
Fetch any movie detail from any imdb list. Returns a JSON list

## Example 

## Basic Usage
```shell
./imdb https://www.imdb.com/india/top-rated-indian-movies/ 2

[{"title":"Pather Panchali ","movie_releasing_year":"(1955)","imdb_rating":"8.6","summary":"Impoverished priest Harihar Ray, dreaming of a better life for himself and his family, leaves his rural Bengal village in search of work.","duration":"2h 5min","genre":"Drama"},{"title":"Ratsasan ","movie_releasing_year":"(2018)","imdb_rating":"8.7","summary":"A Sub-Inspector sets out in pursuit of a mysterious serial killer who targets teen school girls and murders them brutally","duration":"2h 50min","genre":"Action"}]

```

## Cache Usage (fast Output)
```shell
./imdb https://www.imdb.com/india/top-rated-indian-movies/ 2 --cache . ## Specify Cache Dir 

```

