package main

import (
	"fmt"
	"io"
	"time"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"log"
	"net/http"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/gorm"
)

// DB 用構造体
type Event struct {
	Released_at time.Time `gorm:"not null"`
	Name string `gorm:"type:text;not null"`
	Idol string `gorm:"type:text;not null"`
	Type int `gorm:"not null"`
}

// テンプレート用構造体
type Row struct {
	Released_at string
	Name string
	Idol string
	Type int
}

type Template struct {
	templates *template.Template
}

var db *gorm.DB

func gormConnect() *gorm.DB {
	db, err := gorm.Open("mysql", "root:pass@/imas?charset=utf8&parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("gorm connect success")
	}
	db.LogMode(true)
	
	// db.SingularTable(true)

	return db
}

func main() {
	log.Println("start main.go")
	db = gormConnect()
	if db.HasTable(&Event{}) {
		println("events table already exists")
	} else {
		println("create events table")
		db.CreateTable(&Event{})
	}
	

	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = t

	e.Use(middleware.Logger()) //全てのリクエストに対してロギング
	e.Use(middleware.Recover())

	e.Static("/css", "public/assets/css")
	e.Static("/js", "public/assets/js")
	e.Static("/img", "public/assets/img")

	e.GET("/", IndexHandler)

	e.Logger.Fatal(e.Start(":14235"))
}

func IndexHandler(c echo.Context) error {
	events := []Event{}
	db.Order("Released_at asc").Find(&events)
	fmt.Println(len(events), "件見つかりました")

	rows := []Row{}
	for _, e := range events {
		var s string = e.Released_at.Format("2006-01-02")
		row := Row{Released_at: s, Name: e.Name, Idol: e.Idol, Type: e.Type}
		rows = append(rows, row)
	}

	// fmt.Println(c.Path())

	return c.Render(http.StatusOK, "index", rows)
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
