package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"github.com/nleof/goyesql"
	"time"
)

var (
	err error
	db *sql.DB
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	queries    goyesql.Queries
	router *chi.Mux
)

type Post struct {
	ID      int    `json: "id"`
	Title   string `json: "title"`
	Content string `json: "content"`
	CreatedAt time.Time `json: "created_at"`
}

func init() {
	PostgresHost = os.Getenv("POSTGRE_HOST")
	PostgresPort = os.Getenv("POSTGRE_PORT")
	PostgresUser = os.Getenv("POSTGRE_USER")
	PostgresPassword = os.Getenv("POSTGRE_PASSWORD")

	queries = goyesql.MustParseFile("queries/queries.sql")

	router = chi.NewRouter()
	router.Use(middleware.Recoverer)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		PostgresHost, PostgresPort, PostgresUser, PostgresPassword, "go-crud")
	db, err = sql.Open("postgres", psqlInfo)

	catch(err)
}

// UpdatePost update a  specific post
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	id := chi.URLParam(r, "id")
	json.NewDecoder(r.Body).Decode(&post)

	query, err := db.Prepare(fmt.Sprintf(queries["update-post"], "$1", "$2", "$3"))
	fmt.Println(fmt.Sprintf(queries["update-post"], "$1", "$2", "$3"))
	catch(err)

	_, err = query.Exec(post.Title, post.Content, id)
	catch(err)

	defer query.Close()

	respondwithJSON(w, http.StatusOK, map[string]string{"message": "update successfully"})

}

// DeletePost remove a specific post
func DeletePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query, err := db.Prepare(fmt.Sprintf(queries["delete-post"], "$1"))
	catch(err)

	_, err = query.Exec(id)
	catch(err)
	query.Close()

	respondwithJSON(w, http.StatusOK, map[string]string{"message": "successfully deleted"})
}

// CreatePost create a new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	json.NewDecoder(r.Body).Decode(&post)
	fmt.Println(post)

	query, err := db.Prepare(fmt.Sprintf(queries["create-post"], "$1", "$2"))
	catch(err)

	_, err = query.Exec(post.Title, post.Content)
	catch(err)
	defer query.Close()

	respondwithJSON(w, http.StatusCreated, map[string]string{"message": "successfully created"})
}

func AllPosts(w http.ResponseWriter, r *http.Request) {
	posts := []Post{}

	rows, err := db.Query(queries["select-all"])
	if err != nil {
		catch(err)
	}
	defer rows.Close()
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.ID, &post.Content, &post.Title, &post.CreatedAt)
		if err != nil {
			catch(err)
		}
		posts = append(posts, post)
	}

	respondwithJSON(w, http.StatusOK, posts)
}

func DetailPost(w http.ResponseWriter, r *http.Request) {
	post := Post{}
	id := chi.URLParam(r, "id")

	db.QueryRow(fmt.Sprintf(queries["select-post"], "$1"), id).
		Scan(&post.ID, &post.Content, &post.Title, &post.CreatedAt)
	if err != nil {
		catch(err)
	}

	respondwithJSON(w, http.StatusOK, post)
}

func routers() *chi.Mux {
	router.Get("/posts", AllPosts)
	router.Get("/posts/{id}", DetailPost)
	router.Post("/posts", CreatePost)
	router.Put("/posts/{id}", UpdatePost)
	router.Delete("/posts/{id}", DeletePost)

	return router
}

func main () {
	routers()
	http.ListenAndServe(":8005", Logger())
}