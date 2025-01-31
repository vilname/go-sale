package helper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func GetWebClient(url string) []byte {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	return body
}

func AddQueryParamSlice(query *url.Values, paramName string, params []uint64) {
	if params != nil {
		for _, param := range params {
			query.Add(paramName, strconv.FormatUint(param, 10))
		}
	}
}
