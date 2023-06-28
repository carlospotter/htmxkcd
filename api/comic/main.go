package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Comic struct {
	Image  string `json:"img"`
	Title  string `json:"title"`
	Number int    `json:"num"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// lc, ok := lambdacontext.FromContext(ctx)
	// if !ok {
	// 	return &events.APIGatewayProxyResponse{
	// 		StatusCode: 503,
	// 		Body:       "Something went wrong :(",
	// 	}, nil
	// }

	// make http request
	comic, err := getComic("")
	if err != nil {
		fmt.Println(err.Error())
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	comicCard := getComicCard(comic)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       comicCard,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func getComic(queryParam string) (Comic, error) {
	response, err := http.Get("https://xkcd.com/info.0.json")
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

	return comic, nil
}

func getComicCard(comic Comic) string {
	return fmt.Sprintf(`
		<div class="mt-5"> 
			<p class="text-3xl text-center" style="font-family: 'Shadows Into Light', cursive;">%d - %s</p>
			<img class="mx-auto object-contain my-5" src="%s" alt="comic">
			<p class="text-sm text-center mt-5" style="font-family: 'Shadows Into Light', cursive;">Source: <a href="https://xkcd.com/%d" target="_blank">https://xkcd.com/%d</a></p>
		</div>`, comic.Number, comic.Title, comic.Image, comic.Number, comic.Number)
}
