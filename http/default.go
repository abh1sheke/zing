package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/abh1sheke/postx/parser"
	"github.com/abh1sheke/postx/result"
)

func DefaultRequest(
	id int,
	c chan *result.Data,
	client *http.Client,
	args *parser.Args,
	wg *sync.WaitGroup,
	logger *log.Logger,
) {
	defer wg.Done()
	var request *http.Request
	var response *http.Response
	var err error
	if *args.Method == "GET" || *args.Method == "HEAD" {
		request, err = http.NewRequest(*args.Method, *args.URL, nil)
	} else {
		data := extractData(args)
		request, err = http.NewRequest(*args.Method, *args.URL, data)
	}

	if err != nil {
		fmt.Println("could not create http request;")
		fmt.Println(`check logs by running "cat $TMPDIR/postx.log".`)
		logger.Printf("could not create http request: %v\n", err)
		os.Exit(1)
	}

	request.Header.Set("User-Agent", "postx/0.1")

	for _, v := range *args.Headers {
		f := strings.Index(v, "=")
		request.Header.Add(v[0:f], v[f+1:])
	}
	response, err = client.Do(request)
	if err != nil {
		fmt.Println("could not perform http request;")
		fmt.Println(`check logs by running "cat $TMPDIR/postx.log".`)
		logger.Printf("could not perform http request: %v\n", err)
		os.Exit(1)
	}
	var body []byte
	body, err = io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("could not read response body.")
		fmt.Println(`check logs by running "cat $TMPDIR/postx.log".`)
		logger.Printf("could not perform http request: %v\n", err)
	} else {
		c <- &result.Data{Body: &body, Response: response}
	}
}
