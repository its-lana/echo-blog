package routes

import (
	"echo-blog/controllers"
	"echo-blog/middlewares"
	"net/http"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	e.GET("/", defaultHandler)
	middlewares.LogMiddlewares(e)

	v1 := e.Group("/api/v1")
	v1Auth := e.Group("/api/v1", middlewares.UserAuthMiddlewares())

	//user login
	v1.POST("/login", controllers.LoginUser)

	//api Blog
	v1.GET("/blogs", controllers.GetAllBlogs)
	v1.GET("/blogs/:id", controllers.GetBlogByID)
	v1Auth.POST("/blogs", controllers.AddNewBlog)
	v1Auth.PUT("/blogs/:id", controllers.UpdateBlog)
	v1Auth.DELETE("/blogs/:id", controllers.DeleteBlog)

	//api User
	v1Auth.GET("/users", controllers.GetAllUser)
	v1Auth.GET("/users/:id", controllers.GetUserByID)
	v1.POST("/users", controllers.AddNewUser)
	v1Auth.PUT("/users/:id", controllers.UpdateUser)
	v1Auth.DELETE("/users/:id", controllers.DeleteUser)

	e.Any("*", catchAllHandler)

	return e
}

func defaultHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to echo-blog!")
}

func catchAllHandler(c echo.Context) error {
	return c.String(http.StatusNotFound, "Sorry, the route path you're looking for doesn't exist!")
}
