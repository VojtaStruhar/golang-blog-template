package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"strconv"
)

var router *chi.Mux
var db *sql.DB

type Article struct {
	ID      int           `json:"id"`
	Title   string        `json:"title"`
	Content template.HTML `json:"content"`
}

func main() {
	router = chi.NewRouter()
	/*
		Note the use of the Recoverer middleware.
		When the catch() function panics, this middleware will recover the server, log the error with a stack trace,
		and send a 500 Internal Server Error response to the client.
	*/
	router.Use(middleware.Recoverer)

	var err error
	db, err = connect()
	catch(err)

	/*
		ChangeMethod: This middleware will change the request method to PUT or DELETE if the request method is POST and the
		form field _method is set to PUT or DELETE, respectively. This is required because HTML forms only support GET and
		POST methods.
	*/
	router.Use(ChangeMethod)
	router.Get("/", GetAllArticles)
	router.Route("/articles", func(r chi.Router) {
		r.Get("/", NewArticle)
		r.Post("/", CreateArticle)
		r.Route("/{articleID}", func(r chi.Router) {
			r.Use(ArticleCtx)
			r.Get("/", GetArticle)       // GET /articles/1234
			r.Put("/", UpdateArticle)    // PUT /articles/1234
			r.Delete("/", DeleteArticle) // DELETE /articles/1234
			r.Get("/edit", EditArticle)  // GET /articles/1234/edit
		})
	})

	err = http.ListenAndServe(":8005", router)
	catch(err)
}

func catch(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func ChangeMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch method := r.PostFormValue("_method"); method {
			case http.MethodPut:
				fallthrough
			case http.MethodPatch:
				fallthrough
			case http.MethodDelete:
				r.Method = method
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}

func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID := chi.URLParam(r, "articleID")
		article, err := dbGetArticle(articleID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), "article", article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := dbGetAllArticles()
	catch(err)

	fmt.Println("GetAllArticles:", articles)

	t, _ := template.ParseFiles("templates/index.html")
	err = t.Execute(w, articles)
	catch(err)
}

func NewArticle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/new.html")
}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")
	article := &Article{
		Title:   title,
		Content: template.HTML(content),
	}

	err := dbCreateArticle(article)
	catch(err)
	http.Redirect(w, r, "/", http.StatusFound)
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*Article)
	fmt.Println(article)
	//TODO: Render template
}

func EditArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*Article)
	fmt.Println(article)
	// TODO: Render template
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*Article)

	title := r.FormValue("title")
	content := r.FormValue("content")
	newArticle := &Article{
		Title:   title,
		Content: template.HTML(content),
	}

	err := dbUpdateArticle(strconv.Itoa(article.ID), newArticle)
	catch(err)
	http.Redirect(w, r, fmt.Sprintf("/articles/%d", article.ID), http.StatusFound)
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value("article").(*Article)
	err := dbDeleteArticle(strconv.Itoa(article.ID))
	catch(err)

	http.Redirect(w, r, "/", http.StatusFound)
}
