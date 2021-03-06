package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/madjlzz/madlens/controllers"
	"github.com/madjlzz/madlens/middleware"
	"github.com/madjlzz/madlens/models"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "madlens_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()

	//us.DestructiveReset()
	services.AutoMigrate()

	r := mux.NewRouter()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMiddleware := middleware.RequireUser{
		User: userMw,
	}

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	// Image routes
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	// Gallery routes
	r.Handle("/galleries", requireUserMiddleware.ApplyFn(galleriesC.Index)).Methods("GET")
	r.Handle("/galleries/new", requireUserMiddleware.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMiddleware.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGalleryName)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMiddleware.ApplyFn(galleriesC.Edit)).Methods("GET").Name(controllers.EditGalleryName)
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMiddleware.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMiddleware.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMiddleware.ApplyFn(galleriesC.ImageUpload)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMiddleware.ApplyFn(galleriesC.ImageDelete)).Methods("POST")

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", userMw.Apply(r))
}
