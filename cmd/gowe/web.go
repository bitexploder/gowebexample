package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/bitexploder/gowebexample/context"
	"github.com/bitexploder/gowebexample/model"

	"github.com/asdine/storm"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

var Db *storm.DB

var store = sessions.NewCookieStore(
	[]byte("You probably want to change this"),
	[]byte("Seriously, I mean it. Change it."))

///// HELPERS
func loadTmpl(path string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		fmt.Errorf("Error parsing template: %s", path)
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		fmt.Println(err)
	}

	return buf.String(), err
}

///// HTTP Handlers
func HomeHandler(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, _ := loadTmpl(config.StaticDir+"/index.html", nil)
		fmt.Fprint(w, s)
	}
}

func ListUsersHandler(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := model.GetUsers(Db)
		if err != nil {
			fmt.Fprintf(w, "err: %+v\n")
			return
		}

		s, err := loadTmpl(config.TemplateDir+"/list.html", users)
		if err != nil {
			fmt.Printf("error loading template: %s\n", err)
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Fprint(w, s)
	}
}

func intVar(vars map[string]string, k string) int64 {
	var vv int64
	if v, ok := vars[k]; ok {
		vv, _ = strconv.ParseInt(v, 0, 32)
	}
	return vv
}

func EditUserHandler(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		vars := mux.Vars(r)
		id := int(intVar(vars, "id"))
		user := model.User{}

		if id != 0 {
			user, err = model.GetUser(Db, id)
			if err != nil {
				fmt.Fprintf(w, "err: %s\n", err)
				return
			}
		}

		if r.Method == "POST" {
			r.ParseForm()
			user.Name = r.Form["name"][0]
			user.Email = r.Form["email"][0]
			user.Username = r.Form["username"][0]

			err = model.UpdateUser(Db, user)
			if err != nil {
				fmt.Fprintf(w, "err: %s\n", err)
			}
			http.Redirect(w, r, "/users", 301)

		}

		s, err := loadTmpl(config.TemplateDir+"/edit.html", user)
		if err != nil {
			fmt.Printf("error loading template: %s\n", err)
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprint(w, s)
	}
}

func DeleteUserHandler(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		vars := mux.Vars(r)
		id := int(intVar(vars, "id"))

		err = model.DeleteUser(Db, id)
		if err != nil {
			fmt.Fprintf(w, "err: %s", err)
			return
		}

		fmt.Fprintf(w, "Deleting user: %d", id)
	}
}

///// Authentication
func LoginHandler(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := context.Get(r, "session").(*sessions.Session)

		loginTmpl := config.TemplateDir + "/login.html"

		params := struct {
			Flashes []interface{}
		}{}

		if r.Method == "GET" {
			params.Flashes = session.Flashes()
			s, err := loadTmpl(loginTmpl, params)
			if err != nil {
				fmt.Printf("error loading template: %s\n", err)
				http.Error(w, err.Error(), 500)
				return
			}
			session.Save(r, w)
			fmt.Fprint(w, s)

		}

		if r.Method == "POST" {
			r.ParseForm()
			username := r.Form["username"][0]
			password := r.Form["password"][0]
			u, err := model.GetUserByUsername(Db, username)
			if err != nil {
				session.AddFlash("err: " + err.Error())
				err = session.Save(r, w)
				if err != nil {
					log.Printf("error saving session: %s\n", err)
				}
				http.Redirect(w, r, "/login", 301)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(u.PassHash), []byte(password))
			if err != nil {
				session.AddFlash("err: " + err.Error())
				err = session.Save(r, w)
				if err != nil {
					log.Printf("error saving session: %s\n", err)
				}

				http.Redirect(w, r, "/login", 301)
				return
			}

			session.Values["id"] = u.ID
			err = session.Save(r, w)
			if err != nil {
				log.Printf("error saving session: %s\n", err)
			}
			http.Redirect(w, r, "/", 301)
		}
	}
}

func LogoutHandler(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := context.Get(r, "session").(*sessions.Session)
		delete(session.Values, "id")
		session.Save(r, w)
		http.Redirect(w, r, "/login", 301)
	}
}

///// MIDDLEWARE
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}

func ContextManager(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "gowe")
		if err != nil {
			log.Printf("ContextManager: err: %s\n", err)
			return
		}

		r = context.Set(r, "session", session)

		if id, ok := session.Values["id"]; ok {
			u, err := model.GetUser(Db, id.(int))
			if err != nil {
				r = context.Set(r, "user", nil)
			} else {
				r = context.Set(r, "user", u)
			}
		} else {
			r = context.Set(r, "user", nil)
		}

		h.ServeHTTP(w, r)

		context.Clear(r)
	})
}

func RequireLogin(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u := context.Get(r, "user"); u != nil {
			h.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", 302)
		}
	})
}

func Logger(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s", r.URL)
		h.ServeHTTP(w, r)
	})
}

////////////////
type Config struct {
	Listen      string
	DbPath      string
	StaticDir   string
	TemplateDir string
}

func Web(config Config) {
	router := mux.NewRouter()
	router.HandleFunc("/", Use(HomeHandler(config), RequireLogin))
	router.HandleFunc("/login", LoginHandler(config))
	router.HandleFunc("/logout", Use(LogoutHandler(config), RequireLogin))
	router.HandleFunc("/users", Use(ListUsersHandler(config), RequireLogin))
	router.HandleFunc("/edit/{id:[0-9]+}", Use(EditUserHandler(config), RequireLogin))
	router.HandleFunc("/edit/{id:[0-9]+}/delete", Use(DeleteUserHandler(config), RequireLogin))

	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(
			http.Dir(config.StaticDir)))) // Is this Lisp?

	h := Use(router.ServeHTTP, Logger, ContextManager)

	http.ListenAndServe(config.Listen, h)
}

func main() {
	listen := flag.String("listen", "127.0.0.1:8080", "Listen address and port")
	dbPath := flag.String("dbpath", "gowe.db", "Database path")
	staticDir := flag.String("static", "static", "Static files to serve")
	templateDir := flag.String("template", "template", "Template file directory")
	flag.Parse()

	c := Config{
		Listen:      *listen,
		DbPath:      *dbPath,
		StaticDir:   *staticDir,
		TemplateDir: *templateDir,
	}

	var err error
	Db, err = storm.Open(c.DbPath)
	if err != nil {
		log.Printf("err: %s\n", err)
		return
	}

	defer Db.Close()

	Web(c)
}
