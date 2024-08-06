package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	endpoint    string = "http://localhost:8080/"
	contentType string = "application/x-www-form-urlencoded"
)

func inputURL() string {
	fmt.Println("Введите длинный URL")
	reader := bufio.NewReader(os.Stdin)

	long, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	return strings.TrimSuffix(long, "\n")
}

func sendURL(longURL string, verbose bool) string {
	data := url.Values{}
	data.Set("url", longURL)

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	request.Header.Add("Content-Type", contentType)

	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	if verbose {
		fmt.Println("Статус-код ", response.Status)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	return string(body)
}

func main() {

	longURL := inputURL() // вводим "длинную" url ссылку

	shortURL := sendURL(longURL, true) // получаем "короткую" url ссылку

	fmt.Println(shortURL)
}
