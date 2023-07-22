package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Comic struct {
	Image    string `json:"img"`
	Title    string `json:"title"`
	Number   int    `json:"num"`
	Previous int
	Random   int
	Next     int
	Last     int
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var comic Comic
	var err error
	var comicCard string

	comicNum := request.QueryStringParameters["comic_num"]

	comic, err = get(comicNum)
	if err != nil {
		fmt.Println(err.Error())
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	comicCard = getComicCard(comic)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       comicCard,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func get(num string) (Comic, error) {
	rand.Seed(time.Now().UnixNano())

	responseLast, err := http.Get("https://xkcd.com/info.0.json")
	if err != nil {
		return Comic{}, err
	}
	defer responseLast.Body.Close()

	bodyLast, err := io.ReadAll(responseLast.Body)
	if err != nil {
		return Comic{}, err
	}

	var lastComic Comic
	err = json.Unmarshal(bodyLast, &lastComic)
	if err != nil {
		return Comic{}, err
	}

	lastComic.Last = lastComic.Number
	lastComic.Next = lastComic.Number
	lastComic.Previous = lastComic.Number - 1
	lastComic.Random = rand.Intn(lastComic.Last-1) + 1

	numInt, _ := strconv.Atoi(num)

	if num == "" || numInt == lastComic.Number {
		return lastComic, nil
	}

	response, err := http.Get(fmt.Sprintf("https://xkcd.com/%s/info.0.json", num))
	if err != nil {
		return Comic{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Comic{}, err
	}

	var comic Comic
	err = json.Unmarshal(body, &comic)
	if err != nil {
		return Comic{}, err
	}

	comic.Last = lastComic.Last
	comic.Next = comic.Number + 1
	comic.Previous = comic.Number - 1
	comic.Random = rand.Intn(comic.Last-1) + 1

	return comic, nil
}

func getComicCard(comic Comic) string {
	return fmt.Sprintf(`
		<div class="max-w-2xl mx-auto my-8 grid grid-cols-1 sm:hidden">
			<button hx-get="/.netlify/functions/comic?comic_num=%d" hx-trigger="click" hx-swap="innerHTML" hx-target="#comic"  class="w-20 mx-auto text-2xl border border-black rounded" style="font-family: 'Shadows Into Light', cursive;">
				Random
			</button>
		</div>
		<div class="max-w-2xl mx-auto my-8 grid grid-cols-4 sm:grid-cols-5">
			<button hx-get="/.netlify/functions/comic?comic_num=1" hx-trigger="click" hx-swap="innerHTML" hx-target="#comic" class="w-16 mx-auto text-2xl border border-black rounded" style="font-family: 'Shadows Into Light', cursive;">
				First
			</button>

			<button hx-get="/.netlify/functions/comic?comic_num=%d" hx-trigger="click" hx-swap="innerHTML" hx-target="#comic"  class="w-16 mx-auto text-2xl border border-black rounded" style="font-family: 'Shadows Into Light', cursive;">
				←
			</button>

			<button hx-get="/.netlify/functions/comic?comic_num=%d" hx-trigger="click" hx-swap="innerHTML" hx-target="#comic" class="w-20 mx-auto text-2xl border hidden sm:block border-black rounded" style="font-family: 'Shadows Into Light', cursive;">
				Random
			</button>

			<button hx-get="/.netlify/functions/comic?comic_num=%d" hx-trigger="click" hx-swap="innerHTML" hx-target="#comic" class="w-16 mx-auto text-2xl border border-black rounded" style="font-family: 'Shadows Into Light', cursive;">
				→
			</button>

			<button hx-get="/.netlify/functions/comic" hx-trigger="click" hx-swap="innerHTML" hx-target="#comic" class="w-16 mx-auto text-2xl border border-black rounded" style="font-family: 'Shadows Into Light', cursive;">
				Last
			</button>
    	</div>
		<div class="mt-5"> 
			<p class="text-3xl text-center" style="font-family: 'Shadows Into Light', cursive;">%d - %s</p>
			<img class="mx-auto object-contain my-5" src="%s" alt="comic">
			<p class="text-sm text-center mt-5" style="font-family: 'Shadows Into Light', cursive;">Source: <a href="https://xkcd.com/%d" target="_blank">https://xkcd.com/%d</a></p>
		</div>`, comic.Random, comic.Previous, comic.Random, comic.Next, comic.Number, comic.Title, comic.Image, comic.Number, comic.Number)
}
