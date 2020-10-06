package zi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"zi/api"
)

var url string

// Get gets a key from database selected in Zi function
type get func(string) api.Pair
type set func(api.Pair) api.Pair
type del func(string) string
type rename func(old string, new string) string
type getAll func() []api.Pair

// ZI is main struct interface
type ZI struct {
	Get    get
	Set    set
	Del    del
	Rename rename
	GetAll getAll
}

var ziGoodReturn = ZI{Get: func(key string) api.Pair {
	u := url + "/get?key=" + key
	data, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Body.Close()

	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		log.Fatal(err)
	}
	r := api.Pair{}
	json.Unmarshal([]byte(body), &r)
	return r
}, Set: func(d api.Pair) api.Pair {
	in, _ := json.Marshal(d)
	u := url + "/set?data=" + strings.ReplaceAll(string(in), " ", "%20")
	data, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Body.Close()

	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		log.Fatal(err)
	}
	r := api.Pair{}
	json.Unmarshal([]byte(body), &r)
	return r
}, Del: func(key string) string {
	u := url + "/del?key=" + key
	data, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Body.Close()
	return "OK"
}, Rename: func(old string, new string) string {
	u := url + "/rename?origin=" + old + "&new=" + new
	data, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Body.Close()
	return "OK"
}, GetAll: func() []api.Pair {
	u := url + "/getall"
	data, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Body.Close()

	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		log.Fatal(err)
	}
	r := []api.Pair{}
	json.Unmarshal([]byte(body), &r)
	return r
}}

// Zi is the main function for the Zi go library.
func Zi(u string) (ZI, error) {
	url = u
	data, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Body.Close()

	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		log.Fatal(err)
	}
	if string(body) != "OK" {
		return ziGoodReturn, errors.New("Not valid zi database")
	} else {
		return ziGoodReturn, nil
	}
}