package handler

import (
	"encoding/json"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/labstack/echo/v4"
	"go-echo/model"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"net/http"
	"time"
)

func (h *Handler) TestHome(c echo.Context) (err error) {
	return c.JSON(http.StatusOK,map[string]string{"message": "Hello World"})
}
func (h *Handler) TestStream(c echo.Context) (err error) {
	type (
		Geolocation struct {
			Altitude  float64
			Latitude  float64
			Longitude float64
		}
	)
	var (
		locations = []Geolocation{
			{-97, 37.819929, -122.478255},
			{1899, 39.096849, -120.032351},
			{2619, 37.865101, -119.538329},
			{42, 33.812092, -117.918974},
			{15, 37.77493, -122.419416},
		}
	)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)

	enc := json.NewEncoder(c.Response())
	for _, l := range locations {
		if err := enc.Encode(l); err != nil {
			return err
		}
		c.Response().Flush()
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}
func (h *Handler) InitUsers(c echo.Context) (err error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("foobar"), bcrypt.DefaultCost)
	u := &model.User{ID: bson.NewObjectId()}
	u.Email="foo@bar.com"
	u.Password= string(hashedPassword)
	db := h.DB.Clone()
	defer db.Close()
	if err = db.DB("twitter").C("users").Insert(u); err != nil {

	}
	for n := 0; n < 10; n++ {
		u := &model.User{ID: bson.NewObjectId()}
		u.Email=gofakeit.Email()
		u.Password= string(hashedPassword)
		db := h.DB.Clone()
		defer db.Close()
		if err = db.DB("twitter").C("users").Insert(u); err != nil {

		}
	}
	return c.JSON(http.StatusOK,map[string]string{"message": "Users Seeded"})
}
func (h *Handler) InitFeed(c echo.Context) (err error) {
	db := h.DB.Clone()
	users := []*model.User{}

	if err = db.DB("twitter").C("users").
		Find(nil).
		//Limit(5).
		Select(bson.M{"ID":1}).
		All(&users); err != nil {
		return err
	}
	defer db.Close()
	for _,u:=range users{
		for n := 0; n < 10; n++ {
			p := &model.Post{
				ID:   bson.NewObjectId(),
				From: u.ID.Hex(),
				To: users[rand.Intn(len(users))].ID.Hex(),
			}
			p.Message=gofakeit.HipsterSentence(10)
			if err = db.DB("twitter").C("posts").Insert(p); err != nil {
				return
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Posts Seeded"})
}
