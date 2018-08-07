package main

import (
	"html"
	"html/template"
	"io"
	"net/http"

	"github.com/ipfans/echo-session"
	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

type LoginForm struct {
	UserId       string
	Password     string
	ErrorMessage string
}

type Complete struct {
	Success bool
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	// set template
	t := &Template{
		templates: template.Must(template.ParseGlob("../views/*.html")),
	}
	e.Renderer = t

	// set session
	store := session.NewCookieStore([]byte("secret-key"))
	store.MaxAge(90000)
	e.Use(session.Sessions("Error Session", store))
	e.GET("/login", ShowLoginHtml)
	e.POST("/login", Login)

	// up server
	e.Logger.Fatal(e.Start(":8080"))
	e.Logger.Fatal(e.StartAutoTLS(":443"))
}

func ShowLoginHtml(c echo.Context) error {
	session := session.Default(c)

	loginId := session.Get("loginCompleted")
	if loginId != nil && loginId == "completed" {

	}
	return c.Render(http.StatusOK, "login", LoginForm{})
}

func Login(c echo.Context) error {
	loginForm := LoginForm{
		UserId:   c.FormValue("userId"),
		Password: c.FormValue("password"),
	}

	userId := html.EscapeString(loginForm.UserId)
	password := html.EscapeString(loginForm.Password)

	if userId != "userId" && password != "password" {
		loginForm.ErrorMessage = "ユーザーID または パスワードが間違っています"
		return c.Render(http.StatusOK, "login", loginForm)
	}

	//セッションにデータを保存する
	session := session.Default(c)
	session.Set("loginCompleted", "completed")
	session.Save()

	return c.Render(http.StatusOK, "login", LoginForm{})

}
