package usergen

import (
	"encoding/json"
	"fmt"
	"net/http"
	"randomusergen/apiclient"
	"randomusergen/domain"
	"randomusergen/repo"

	"github.com/gin-gonic/gin"
)

type Saver struct {
	client apiclient.Client
	db     repo.UserRepo
}

func New(c apiclient.Client, db repo.UserRepo) *Saver {
	return &Saver{
		client: c,
		db:     db,
	}
}

func (s *Saver) CreateUsers() func(c *gin.Context) {
	return func(c *gin.Context) {
		var req domain.Request
		err := json.NewDecoder(c.Request.Body).Decode(&req)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, "can't decode request")
			return
		}

		query := domain.MapToQuery(req.QueryParams)
		fmt.Println(query)

		res, err := s.client.Get(c.Request.Context(), query)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, "can't get users")
			return
		}

		// fmt.Println(res.Results)

		n, err := s.db.SaveAll(c.Request.Context(), res.Results)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, "can't save users")
			return
		}

		c.JSON(http.StatusOK, fmt.Sprint(n, " users saved"))
	}
}
