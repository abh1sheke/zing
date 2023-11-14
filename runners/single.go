package runners

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	lhttp "github.com/abh1sheke/postx/http"
	"github.com/abh1sheke/postx/logging"
	"github.com/abh1sheke/postx/parser"
	"github.com/abh1sheke/postx/result"
)

func Single(args *parser.Args) {
	r := result.InitResultList(uint(*args.Parallel))
	var method RequestFunc

	switch *args.Method {
	case "FORM":
		method = lhttp.FormRequest
	default:
		if args.Files != nil {
			method = lhttp.MultipartRequest
		} else {
			method = lhttp.DefaultRequest
		}
	}

	c := make(chan *result.Data)
	startTime := time.Now()
	defer func() {
		c <- nil
		logging.SaveOutput(args, r)
		if *args.Time {
			fmt.Printf(
				"took %vms for %v requests.\n",
				time.Since(startTime).Milliseconds(),
				*args.Parallel,
			)
		}
	}()
	go r.Consumer(c)
	wg := new(sync.WaitGroup)
	client := new(http.Client)
	for i := 1; i <= *args.Parallel; i++ {
		wg.Add(1)
		go method(i, c, client, args, wg)
	}
	wg.Wait()
}
