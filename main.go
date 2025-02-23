package main

import (
	"fmt"
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
)

type Article struct {
Id uint16
Title, Anons, FullText string 
}

var posts = []Article{}
var showPost = Article{}




func index(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")


if err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}

db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Выборка данных
res, err := db.Query("SELECT * FROM `articles`")
	fmt.Println("Подключено к  MySQL!")
	if err != nil {
		panic(err)
	}

	posts = []Article{}
	for res.Next(){
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		if err != nil {
			http.Error (w,"Ошибка при парсинге шаблона: " + err.Error(), http.StatusInternalServerError)
			return
		}

		posts = append(posts, post)
		
	}

	//t.ExecuteTemplate(w, "index", posts)
	err = t.ExecuteTemplate(w, "index", posts)
    if err != nil {
        http.Error(w, "Ошибка при выполнении шаблона: " + err.Error(), http.StatusInternalServerError)
        return
    }
	
}


func create(w http.ResponseWriter, r *http.Request){
	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	
	t.ExecuteTemplate(w, "create", nil)
	}



	func save_article(w http.ResponseWriter, r *http.Request){
		title := r.FormValue("title")
		anons := r.FormValue("anons")
		full_text := r.FormValue("full_text")


if title == "" || anons == "" || full_text == "" {
fmt.Fprintf(w, "Не все поля заполнены")
} else{

db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
if err != nil {
	panic(err)
}
defer db.Close()

//Установка данных
insert, err :=db.Query(fmt.Sprintf ("INSERT INTO `articles` (`title`, `anons`, `full_text`) Values ('%s', '%s', '%s')", title, anons, full_text))

if err != nil {
panic(err)
}
defer insert.Close()

http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}


func show_post(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")


if err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golang")
	if err != nil {
		panic(err)
	}
	defer db.Close()
		// Выборка данных
err = db.QueryRow("SELECT id, title, anons, full_text FROM articles WHERE id = ?", vars["id"]).Scan(&showPost.Id, &showPost.Title, &showPost.Anons, &showPost.FullText)
fmt.Println("Подключено к  MySQL!")
if err != nil {
	panic(err)
}

// No need to loop through results, as we are fetching a single row

//t.ExecuteTemplate(w, "index", posts)
err = t.ExecuteTemplate(w, "show", showPost)
if err != nil {
	http.Error(w, "Ошибка при выполнении шаблона: " + err.Error(), http.StatusInternalServerError)
	return
}

}

func handleFunc(){
rtr := mux.NewRouter()
rtr.HandleFunc("/", index).Methods("GET")
rtr.HandleFunc("/create", create).Methods("GET")
rtr.HandleFunc("/save_article", save_article).Methods("POST")
rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

http.Handle("/", rtr)
http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
http.ListenAndServe(":8080", nil)
}




func main() {
handleFunc()
}
//