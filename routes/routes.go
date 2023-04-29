package routes

import (
	"sodality/controllers"
	middlewares "sodality/handlers"

	"github.com/gorilla/mux"
)

// Routes -> define endpoints
func Routes() *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()

	// User API routes
	user := api.PathPrefix("/user").Subrouter()
	user.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	user.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	user.HandleFunc("/{id}", controllers.GetUserByID).Methods("GET")
	user.HandleFunc("/profile/", middlewares.IsAuthorized(controllers.GetProfile)).Methods("GET")
	user.HandleFunc("/update", middlewares.IsAuthorized(controllers.UpdateUser)).Methods("PUT")

	// Content API routes
	content := api.PathPrefix("/content").Subrouter()
	content.HandleFunc("/post", middlewares.IsAuthorized(controllers.PostContent)).Methods("POST")
	content.HandleFunc("/all", controllers.GetAllCreatorsContent).Methods("GET")
	content.HandleFunc("/{search}", controllers.SearchContentByTitle).Methods("GET")
	content.HandleFunc("/delete/{id}", middlewares.IsAuthorized(controllers.DeleteContent)).Methods("DELETE")

	// Creator content
	creator := api.PathPrefix("/creator").Subrouter()
	creator.HandleFunc("/content/{id}", controllers.GetContentByID).Methods("GET")
	creator.HandleFunc("/content/category/{category_name}/", controllers.GetCreatorDirectoryByDirectoryName).Methods("GET")
	creator.HandleFunc("/all/content", middlewares.IsAuthorized(controllers.GetOwnContent)).Methods("GET")

	// setting endpoint
	setting := creator.PathPrefix("/setting").Subrouter()
	setting.HandleFunc("/update", middlewares.IsAuthorized(controllers.CreatorSetting)).Methods("PUT")

	// followers endpoint
	creator.HandleFunc("/follow/{creator_id}", middlewares.IsAuthorized(controllers.FollowCreator)).Methods("POST")
	creator.HandleFunc("/unfollow/{creator_id}", middlewares.IsAuthorized(controllers.UnfollowCreator)).Methods("DELETE")
	creator.HandleFunc("/{user_id}/followers/", controllers.GetCreatorFollowers).Methods("GET")

	// supporter endpoint
	creator.HandleFunc("/{creator_id}/supporter/", controllers.GetCreatorSupporter).Methods("GET")

	// ipfs endpoint
	file := api.PathPrefix("/file").Subrouter()
	file.HandleFunc("/upload", controllers.UploadFile).Methods("POST")

	// donation endpoint
	donate := api.PathPrefix("/donate").Subrouter()
	donate.HandleFunc("/", middlewares.IsAuthorized(controllers.DonateUser)).Methods("POST")
	donate.HandleFunc("/content", middlewares.IsAuthorized(controllers.DonateContent)).Methods("POST")

	return router
}
