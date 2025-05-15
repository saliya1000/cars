package main

import (
	"cars-viewer/handlers"
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	if err := handlers.LoadCarsData(); err != nil {
		fmt.Println("Error loading car data:", err)
		return
	}

	funcMap := template.FuncMap{
		"getManufacturerNameByID": handlers.GetManufacturerNameByID,
	}
	tmpl := template.Must(template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/api/", http.StripPrefix("/api/", http.FileServer(http.Dir("api"))))
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	// http.HandleFunc("/", handlers.HomeHandler(tmpl))
	// http.HandleFunc("/compare", handlers.CompareHandler)
	// http.HandleFunc("/car/", handlers.CarDetailHandler)
	http.HandleFunc("/", handlers.RecoveryMiddleware(handlers.HomeHandler(tmpl)))
	http.HandleFunc("/compare", handlers.RecoveryMiddleware(handlers.CompareHandler))
	http.HandleFunc("/car/", handlers.RecoveryMiddleware(handlers.CarDetailHandler))

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
