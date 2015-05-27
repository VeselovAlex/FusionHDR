// fusionhdr project fusionhdr.go
package main

import (
	"flag"
	"fmt"
	"fusionhdr/native"
	"github.com/go-martini/martini"
	"html/template"
	"net/http"
	"os"
)

const (
	OUTPUT_NAME_FORMAT     = "fusion-%s.jpg"
	CACHE_FILE_NAME_FORMAT = "input/cache-%s-%d.%s"
)

const (
	IMAGE_SERVE_PATH = "res/images/"
)

var RootTemplate = template.Must(template.ParseFiles("static/index.html"))
var ResultTemplate = template.Must(template.ParseFiles("static/result.html"))

func handleHome(w http.ResponseWriter, r *http.Request) {
	RootTemplate.Execute(w, nil)
}

func handleView(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("i")
	var out string
	if len(id) == 0 {
		out = "noimage.png"
	} else {
		out = IMAGE_SERVE_PATH + fmt.Sprintf(OUTPUT_NAME_FORMAT, id)
	}
	ResultTemplate.Execute(w, out)
}

func handleFusion(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		http.Error(w, "Не могу получить файлы", http.StatusBadRequest)
		return
	}
	form := r.MultipartForm
	files := form.File["images"]
	id, err := cache(files)
	defer deleteCache(id)
	if err != nil {
		http.Error(w, "Неизвестная ошибка", http.StatusInternalServerError)
		return
	}

	inputName := fmt.Sprintf("input/cache-%s/list.txt", id)
	outputName := fmt.Sprintf(OUTPUT_NAME_FORMAT, id)
	native.RunFusion("processed/"+outputName, inputName)
	if _, err := os.Stat("processed/" + outputName); err != nil {
		http.Error(w, "Ошибка при попытке создать HDR", http.StatusInternalServerError)
		return
	}

	redirectURL := "/view?i=" + id
	reqRedirect, err := http.NewRequest("GET", redirectURL, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, reqRedirect, redirectURL, http.StatusFound)
}

var addr = flag.String("-addr", ":8000", "Address for server to start")

func main() {
	flag.Parse()
	createDirIfNotExists("processed")
	createDirIfNotExists("input")

	m := martini.Classic()
	m.Get("/", handleHome)
	m.Use(martini.Static("processed",
		martini.StaticOptions{
			Prefix: IMAGE_SERVE_PATH,
		}))
	m.Use(martini.Static("static"))
	m.Post("/process", handleFusion)
	m.Get("/view/*", handleView)
	m.RunOnAddr(*addr)
}
