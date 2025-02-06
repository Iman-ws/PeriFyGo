package routes

import (
	"net/http"

	"PeriFyGo/controllers"
	"PeriFyGo/middlewares"
	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	authCtrl := &controllers.AuthController{}
	verifyCtrl := &controllers.VerifyController{}
	productCtrl := &controllers.ProductController{}
	profileCtrl := &controllers.ProfileController{}
	adminCtrl := &controllers.AdminController{}

	// Auth endpoints.
	r.HandleFunc("/register", authCtrl.Register).Methods("POST")
	r.HandleFunc("/login", authCtrl.Login).Methods("POST")
	r.HandleFunc("/verify", verifyCtrl.VerifyEmail).Methods("GET")

	// Profile endpoint.
	r.Handle("/profile", middlewares.AuthMiddleware(http.HandlerFunc(profileCtrl.GetProfile), "")).Methods("GET")

	// Products endpoints.
	r.HandleFunc("/products", productCtrl.GetProducts).Methods("GET")
	r.Handle("/products", middlewares.AuthMiddleware(http.HandlerFunc(productCtrl.CreateProduct), "admin")).Methods("POST")
	r.Handle("/products/{id}", middlewares.AuthMiddleware(http.HandlerFunc(productCtrl.UpdateProduct), "admin")).Methods("PUT")
	r.Handle("/products/{id}", middlewares.AuthMiddleware(http.HandlerFunc(productCtrl.DeleteProduct), "admin")).Methods("DELETE")

	// Admin endpoints.
	r.Handle("/admin/users", middlewares.AuthMiddleware(http.HandlerFunc(adminCtrl.GetAllUsers), "admin")).Methods("GET")
	r.Handle("/admin/users/{id}", middlewares.AuthMiddleware(http.HandlerFunc(adminCtrl.GetUser), "admin")).Methods("GET")
	r.Handle("/admin/users/{id}", middlewares.AuthMiddleware(http.HandlerFunc(adminCtrl.UpdateUser), "admin")).Methods("PUT")
	r.Handle("/admin/users/{id}", middlewares.AuthMiddleware(http.HandlerFunc(adminCtrl.DeleteUser), "admin")).Methods("DELETE")
	// Endpoint for sending message.
	r.Handle("/admin/message", middlewares.AuthMiddleware(http.HandlerFunc(adminCtrl.SendMessage), "admin")).Methods("POST")

	// Serve static files.
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	return r
}
