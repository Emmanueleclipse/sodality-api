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
	user.HandleFunc("/username/{username}", controllers.GetUserByUsername).Methods("GET")

	user.HandleFunc("/profile/", middlewares.IsAuthorized(controllers.GetProfile)).Methods("GET")
	user.HandleFunc("/update", middlewares.IsAuthorized(controllers.UpdateUser)).Methods("PUT")

	// Content API routes
	content := api.PathPrefix("/content").Subrouter()
	content.HandleFunc("/post", middlewares.IsAuthorized(controllers.PostContent)).Methods("POST")
	content.HandleFunc("/all", controllers.GetAllCreatorsContent).Methods("GET")
	content.HandleFunc("/search", controllers.SearchContentByTitle).Methods("GET")
	content.HandleFunc("/delete/{id}", middlewares.IsAuthorized(controllers.DeleteContent)).Methods("DELETE")

	// Creator content
	creator := api.PathPrefix("/creator").Subrouter()
	creator.HandleFunc("/content/{id}", controllers.GetContentByID).Methods("GET")
	creator.HandleFunc("/content/category/{category_name}/", controllers.GetCreatorDirectoryByDirectoryName).Methods("GET")
	creator.HandleFunc("/all/content", middlewares.IsAuthorized(controllers.GetOwnContent)).Methods("GET")
	creator.HandleFunc("/all", controllers.GetAllCreators).Methods("GET")
	creator.HandleFunc("/search", controllers.SearchCreatorByUsername).Methods("GET")
	creator.HandleFunc("/all/content/{creator_id}", controllers.GetCreatorContentById).Methods("GET")
	creator.HandleFunc("/get/category/{category_name}", controllers.GetCreatorsByDirectoryName).Methods("GET")

	// supporter get creator content
	creator.HandleFunc("/content/{id}/supporter", middlewares.IsAuthorized(controllers.GetContentForSpecificSupporterByID)).Methods("GET")
	creator.HandleFunc("/content/all/supporter/", middlewares.IsAuthorized(controllers.GetAllCreatorsContentForSpecificSupporter)).Methods("GET")
	creator.HandleFunc("/{id}/supporter/", middlewares.IsAuthorized(controllers.GetCreatorContentsForSpecificSupporter)).Methods("GET")
	creator.HandleFunc("/content/category/{category_name}/supporter", middlewares.IsAuthorized(controllers.GetCreatorDirectoryByDirectoryNameForSpecificSupporter)).Methods("GET")
	creator.HandleFunc("/{creator_id}/supporters/record", controllers.GetCreatorSupportersRecord).Methods("GET")

	// setting endpoint
	setting := creator.PathPrefix("/setting").Subrouter()
	setting.HandleFunc("/update", middlewares.IsAuthorized(controllers.CreatorSetting)).Methods("PUT")
	setting.HandleFunc("/notification", middlewares.IsAuthorized(controllers.NotificationSetting)).Methods("PUT")
	setting.HandleFunc("/notification/get", middlewares.IsAuthorized(controllers.GetNotificationSetting)).Methods("GET")

	// followers endpoint
	creator.HandleFunc("/follow/{creator_id}", middlewares.IsAuthorized(controllers.FollowCreator)).Methods("POST")
	creator.HandleFunc("/unfollow/{creator_id}", middlewares.IsAuthorized(controllers.UnfollowCreator)).Methods("DELETE")
	creator.HandleFunc("/{user_id}/followers/", controllers.GetCreatorFollowers).Methods("GET")

	// supporter endpoint
	creator.HandleFunc("/{creator_id}/supporters/count", controllers.GetCreatorSupporter).Methods("GET")

	// ipfs endpoint
	// file := api.PathPrefix("/file").Subrouter()
	// file.HandleFunc("/upload", controllers.UploadFile).Methods("POST")

	// donation endpoint
	donate := api.PathPrefix("/donate").Subrouter()
	donate.HandleFunc("/", middlewares.IsAuthorized(controllers.DonateUser)).Methods("POST")
	donate.HandleFunc("/content", middlewares.IsAuthorized(controllers.DonateContent)).Methods("POST")

	// 2FA endpoint
	security := api.PathPrefix("/auth").Subrouter()
	security.HandleFunc("/generate", middlewares.IsAuthorized(controllers.GenerateQR)).Methods("POST")
	security.HandleFunc("/verify", middlewares.IsAuthorized(controllers.VerifyOTP)).Methods("POST")
	security.HandleFunc("/2fa/update", middlewares.IsAuthorized(controllers.Update2FA)).Methods("PUT")

	tiers := creator.PathPrefix("/tiers").Subrouter()
	tiers.HandleFunc("/update", middlewares.IsAuthorized(controllers.AddCreatorTiers)).Methods("PUT")
	tiers.HandleFunc("/{username}", controllers.GetCreatorTierByUserID).Methods("GET")

	return router
}
