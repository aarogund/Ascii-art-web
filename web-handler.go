package main

import (
	"html/template"
	"net/http"
	"strings"
)

type PageData struct {
	Result    template.HTML
	Text      string
	Banner    string
	Color     string
	Align     string
	substring string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	tmpl.Execute(w, PageData{Result: ""})

}
func asciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "400 Bad Request!", http.StatusBadRequest)
		return
	}
	if r.Method == "POST" {
		tmpl, err := template.ParseFiles("template/index.html")
		if err != nil {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		r.ParseForm()
		input := r.FormValue("input")
		input = template.HTMLEscapeString(input)
		colorFlag := r.FormValue("color")
		alignFlag := r.FormValue("align")
		substring := r.FormValue("substring")
		if input == "" {
			http.Error(w, "400 Bad Request!", http.StatusBadRequest)
			return
		}
		
		filename := r.FormValue("filename")
		bannerMap := loadBanner(filename)
		if bannerMap == nil {
			http.Error(w, "404 Not FOund!", http.StatusNotFound)
			return
		}

		words := strings.Fields(input)

		// step 1 - always build plain art first
		// build wordArts - plain for measuring
		plainArts := []string{}
		coloredArts := []string{}

		for _, word := range words {
			// always build plain art for measuring
			plain, err := printArt(word, bannerMap)
			if err != nil {
				http.Error(w, "400 Bad Request!", http.StatusBadRequest)
				return
			}
			plainArts = append(plainArts, plain)

			// build colored art for printing
			if colorFlag != "" {
				coloredArts = append(coloredArts, printArtColorWeb(colorFlag, word, substring, bannerMap))
			} else {
				coloredArts = append(coloredArts, plain)
			}
		}

		gaps := len(plainArts) - 1
		result := ""
		if alignFlag != "" {
			result += (alignment(coloredArts, plainArts, gaps, 800, alignFlag))
		} else {
			result += (printPlain(coloredArts))
		}
		tmpl.Execute(w, PageData{
			Result: template.HTML(result),
			Text:   input,
			Banner: filename,
			Color:  colorFlag,
			Align:  alignFlag,
		})
	}

}

func shareHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Error(w, "400 Bad Request!", http.StatusBadRequest)
		return
	}
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("template/index.html")
		if err != nil {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}

		text := r.URL.Query().Get("text")
		if text == "" {
			http.Error(w, "400 Bad Request!", http.StatusBadRequest)
			return
		}
		filename := r.URL.Query().Get("banner")
		bannerMap := loadBanner(filename)
		if bannerMap == nil {
			http.Error(w, "404 Not FOund!", http.StatusNotFound)
			return
		}

		result, err := printArt(text, bannerMap)
		if err != nil {
			http.Error(w, "500 Internal Server Error!", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, PageData{Result: template.HTML(result)})
	}
}
