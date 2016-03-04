package main

import (
	"encoding/json"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"bytes"

	"github.com/jinzhu/gorm"
	"github.com/nfnt/resize"
	"github.com/serbe/ncp"
)

var (
	urls = []string{
		"http://nnm-club.me/forum/viewforum.php?f=218",
		"http://nnm-club.me/forum/viewforum.php?f=270",
	}
	commands = []string{
		"get",
		"update",
		"name",
		"rating",
        "poster",
	}
)

// App struct variables
type App struct {
	db  gorm.DB
	net *ncp.NCp
	hd  string
}

type config struct {
	Nnm struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	} `json:"nnmclub"`
	Pq struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
		Sslmode  string `json:"sslmode"`
	} `json:"postgresql"`
	Hd string `json:"httpdir"`
}

func getConfig() (config, error) {
	c := config{}
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(file, &c)
	return c, err
}

func contain(args []string, str string) bool {
	result := false
	for _, item := range args {
		if item == str {
			result = true
			return result
		}
	}
	return result
}

func containCommand(args []string) bool {
	result := false
	for _, item := range commands {
		if contain(args, item) {
			result = true
			return result
		}
	}
	return result
}

func exit(err error) {
	if err == nil {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func (a *App) checkName(ncf ncp.Film) ncp.Film {
	if ncf.Name != strings.ToUpper(ncf.Name) {
		return ncf
	}
	name, err := a.getMovieName(ncf)
	if err == nil {
		ncf.Name = name
		return ncf
	}
	return ncf
}

func (a *App) getPoster(url string) (string, error) {
	var poster string
	resp, err := http.Get(url)
	if err != nil {
		return poster, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return poster, err
	}
	img, err := jpeg.Decode(bytes.NewReader(body))
	if err != nil {
		return poster, err
	}
	m := resize.Resize(15, 0, img, resize.Lanczos3)
	outName := strings.Replace(url, "/", "", -1)
	outName = strings.Replace(url, ":", "", -1)
	if len(outName) < 20 {
		outName = outName[:len(outName)-4]
	} else {
		outName = outName[len(outName)-20 : len(outName)-4]
	}
	out, err := os.Create(a.hd + outName + ".jpg")
	if err != nil {
		return poster, err
	}
	defer out.Close()
	jpeg.Encode(out, m, nil)
	return poster, nil
}
