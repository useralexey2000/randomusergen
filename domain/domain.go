package domain

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Name struct {
	Title string `json:"title"`
	First string `json:"first"`
	Last  string `json:"last"`
}

type Street struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
}

type Coordinates struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type Timezone struct {
	Offset      string `json:"offset"`
	Description string `json:"description"`
}

type Location struct {
	Street      *Street      `json:"street,omitempty"`
	City        string       `json:"city"`
	State       string       `json:"state"`
	Country     string       `json:"country"`
	Postcode    Postcode     `json:"postcode"`
	Coordinates *Coordinates `json:"coordinates,omitempty"`
	Timezone    *Timezone    `json:"timezone,omitempty"`
}

type Postcode string

func (p *Postcode) UnmarshalJSON(data []byte) error {
	var i any
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch i.(type) {
	case int:
		tmp := i.(int)
		*p = Postcode(strconv.Itoa(tmp))
		break
	case string:
		*p = Postcode(i.(string))
	}

	return nil
}

type Login struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Md5      string `json:"md5"`
	Sha1     string `json:"sha1"`
	Sha256   string `json:"sha256"`
}

type Dob struct {
	Date time.Time `json:"date"`
	Age  int       `json:"age"`
}

type Registered struct {
	Date time.Time `json:"date"`
	Age  int       `json:"age"`
}

type ID struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Picture struct {
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Thumbnail string `json:"thumbnail"`
}

type UserData struct {
	Gender     string      `json:"gender"`
	Name       *Name       `json:"name"`
	Location   *Location   `json:"location,omitempty"`
	Email      string      `json:"email"`
	Login      *Login      `json:"login,omitempty"`
	Dob        *Dob        `json:"dob,omitempty"`
	Registered *Registered `json:"registered,omitempty"`
	Phone      string      `json:"phone"`
	Cell       string      `json:"cell"`
	ID         *ID         `json:"id,omitempty"`
	Picture    *Picture    `json:"picture,omitempty"`
	Nat        string      `json:"nat"`
}

type Response struct {
	Results []*UserData `json:"results"`
}

type Request struct {
	QueryParams map[string]string `json:"query"`
}

func MapToQuery(m map[string]string) string {
	counter := len(m) - 1
	var b strings.Builder
	for k, v := range m {
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(v)
		if counter != 0 {
			b.WriteString("&")
		}
		counter--
	}

	return b.String()
}
