package main

import "github.com/julienschmidt/httprouter"

func initRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/favicon.ico", faviconHandler)
	router.GET("/ns/favicon.ico", faviconHandler)
	router.GET("/", getSession(indexHandler))
	router.GET("/ns/", getSession(indexHandler))
	router.GET("/ping", getSession(indexPing))
	router.GET("/ns/ping", getSession(indexPing))

	// Clean the session cache.
	router.GET("/ns/clean-sessions", checkPermission(cleanSessionsHandler, "admin"))
	// Changelog page.
	router.GET("/ns/changelog", checkPermission(changelogHandler, "admin"))
	// Test.
	router.GET("/ns/test", checkPermission(testPageHandler, "admin"))
	router.POST("/ns/test/send-email", checkPermission(testSendMailPost, "admin"))

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// ALDO
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Products list page.
	router.GET("/ns/aldo/products", checkPermission(aldoProductsHandler, "read"))
	// Product page.
	router.GET("/ns/aldo/product/:code", checkPermission(aldoProductHandler, "read"))
	// Create product on zunka server.
	router.POST("/ns/aldo/product/:code", checkPermission(aldoProductHandlerPost, "write"))
	// Check product change.
	router.POST("/ns/aldo/product/:code/checked", checkPermission(aldoProductCheckedHandlerPost, "write"))
	// Product removed from site, so remove his reference from the site system.
	router.DELETE("/ns/aldo/product/mongodb_id/:code", checkApiAuthorization(aldoProductMongodbIdHandlerDelete))
	// Categories page.
	router.GET("/ns/aldo/categories", checkPermission(aldoCategoriesHandler, "read"))
	// Save categories.
	router.POST("/ns/aldo/categories", checkPermission(aldoCategoriesHandlerPost, "write"))

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// ALLNATIONS
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Products list page.
	router.GET("/ns/allnations/products", checkPermission(allnationsProductsHandler, "read"))
	// Product page.
	router.GET("/ns/allnations/product/:code", checkPermission(allnationsProductHandler, "read"))
	// Create product on zunka server.
	router.POST("/ns/allnations/product/:code", checkPermission(allnationsProductHandlerPost, "write"))
	// Check product change.
	router.POST("/ns/allnations/product/:code/checked", checkPermission(allnationsProductCheckedHandlerPost, "write"))
	// Product removed from site, so remove his reference from zunkasrv.
	router.DELETE("/ns/allnations/product/zunka_product_id/:code", checkApiAuthorization(allnationsProductZunkaProductIdHandlerDelete))
	// Filter page.
	router.GET("/ns/allnations/filters", checkPermission(allnationsFiltersHandler, "read"))
	// Save filter.
	router.POST("/ns/allnations/filters", checkPermission(allnationsFiltersHandlerPost, "write"))
	// Categories page.
	router.GET("/ns/allnations/categories", checkPermission(allnationsCategoriesHandler, "read"))
	// Save categories.
	router.POST("/ns/allnations/categories", checkPermission(allnationsCategoriesHandlerPost, "write"))
	// Makers page.
	router.GET("/ns/allnations/makers", checkPermission(allnationsMakersHandler, "read"))
	// Save categories.
	router.POST("/ns/allnations/makers", checkPermission(allnationsMakersHandlerPost, "write"))

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// MERCADO LIVRE
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Login
	router.GET("/ns/ml/auth/login", checkPermission(mercadoLivreAuthLoginHandler, "write"))
	// User code
	router.GET("/ns/ml/auth/user", mercadoLivreAuthUserHandler)
	// Notification
	router.GET("/ns/ml/notification", mercadoLivreNotificationHandler)

	// Load user code from zunka server. Used by zunka server in development.
	router.GET("/ns/ml/user/load-code", checkPermission(mercadoLivreLoadUserCode, "read"))

	// Show user code
	router.GET("/ns/ml/user/code", checkPermission(mercadoLivreUserCodeHandler, "read"))
	router.GET("/ns/ml/api/user-code", checkApiAuthorization(mercadoLivreAPIUserCodeHandler))

	// User info
	router.GET("/ns/ml/users/info", checkPermission(mercadoLivreUsersInfoHandler, "read"))
	// User products
	router.GET("/ns/ml/users/products", checkPermission(mercadoLivreUsersProductsHandler, "read"))
	// Active products
	router.GET("/ns/ml/active-products", checkPermission(mercadoLivreActiveProductsHandler, "read"))
	// Product
	router.GET("/ns/ml/product/:id", checkPermission(mercadoLivreRawProductHandler, "read"))

	// Autheticate user.
	// router.GET("/ns/ml/auth/login", checkPermission(mercadoLivreAuthUserHandler, "read"))
	// router.POST("/ns/ml/auth/login", checkPermission(mercadoLivreAuthUserHandlerPost, "write"))

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// AUTHENTICATION
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Signup
	router.GET("/ns/auth/signup", confirmNoLogged(authSignupHandler))
	router.POST("/ns/auth/signup", confirmNoLogged(authSignupHandlerPost))
	router.GET("/ns/auth/signup/confirmation/:uuid", confirmNoLogged(authSignupConfirmationHandler))

	// Signin / signout
	router.GET("/ns/auth/signin", confirmNoLogged(authSigninHandler))
	router.POST("/ns/auth/signin", confirmNoLogged(authSigninHandlerPost))
	router.GET("/ns/auth/signout", authSignoutHandler)

	// Password
	router.GET("/ns/auth/password/recovery", confirmNoLogged(passwordRecoveryHandler))
	router.POST("/ns/auth/password/recovery", confirmNoLogged(passwordRecoveryHandlerPost))
	router.GET("/ns/auth/password/reset", confirmNoLogged(passwordResetHandler))

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// User
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	router.GET("/ns/user/account", checkPermission(userAccountHandler, ""))
	router.GET("/ns/user/change/name", checkPermission(userChangeNameHandler, ""))
	router.POST("/ns/user/change/name", checkPermission(userChangeNameHandlerPost, ""))
	router.GET("/ns/user/change/email", checkPermission(userChangeEmailHandler, ""))
	router.POST("/ns/user/change/email", checkPermission(userChangeEmailHandlerPost, ""))
	router.GET("/ns/user/change/email-confirmation/:uuid", checkPermission(userChangeEmailConfirmationHandler, ""))
	router.GET("/ns/user/change/mobile", checkPermission(userChangeMobileHandler, ""))
	router.POST("/ns/user/change/mobile", checkPermission(userChangeMobileHandlerPost, ""))

	// Entrance.
	router.GET("/ns/user_add", userAddHandler)

	return router
}
