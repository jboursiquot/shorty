package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/jboursiquot/shorty"
)

type Router struct {
	*gin.Engine
	baseURL   string
	shortener *shorty.Shortener
}

func newRouter(cfg *config) (*Router, error) {
	u, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, err
	}

	r := Router{
		Engine:    gin.Default(),
		baseURL:   u.String(),
		shortener: shorty.NewShortener(),
	}

	v1 := r.Group(cfg.Stage + "/v1")
	v1.POST("/shorten", r.handleShorten)
	r.GET("/:key", r.handleResolve)

	return &r, nil
}

func (r *Router) handleShorten(c *gin.Context) {
	type urlToShorten struct {
		URL string `json:"url"`
	}
	var u urlToShorten
	if err := c.ShouldBindJSON(&u); err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key := r.shortener.Shorten(u.URL)
	url := fmt.Sprintf("%s/%s", r.baseURL, key)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (r *Router) handleResolve(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameter: key"})
		return
	}

	url := r.shortener.Resolve(key)
	if url == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}
