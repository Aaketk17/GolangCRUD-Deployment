package routes

import (
	controller "github.com/Aaketk17/GolangCRUD-Deployment/controllers"

	middleware "github.com/Aaketk17/GolangCRUD-Deployment/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(userRoutes *fiber.App) {

	userRoutes.Post("/signup", controller.UserSignUp)
	userRoutes.Post("/login", controller.UserLogin)

	userAuth := userRoutes.Group("/users")
	userAuth.Use(middleware.UserAuthMiddleware)
	userAuth.Post("/addadmin", controller.AddAdminUser)
	userAuth.Get("/getuser/:id", controller.GetUser)
	userAuth.Get("/getusers", controller.GetUsers)
	userAuth.Delete("/delete/:id", controller.DeleteUser)
	userAuth.Put("/update/:id", controller.UpdateUsers)
	userAuth.Post("/logout", controller.Logout)
}
