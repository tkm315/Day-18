package main

// Literals and closurs
// declare makeHandler func
import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

// os.WriteFile
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

// make filename , read context , return *Page
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, _ := os.ReadFile(filename)
	return &Page{Title: title, Body: body}, nil
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

// must is for panic!
// once parse and exe it every where
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m != nil {
		http.NotFound(w, r)
		return " ", errors.New("invalid page title")
	}
	return m[2], nil
}
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// extract title , loadPage(title) , display
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	//title, err := getTitle(w, r)
	//if err != nil {
	//	return
	//}
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// load or create *page , privide editing form
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	/*title, err := getTitle(w, r)
	if err != nil {
		return
	}*/
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// create *page , find title , find body from form , save page
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	/*title, err := getTitle(w, r)
	if err != nil {
		return
	}*/
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
