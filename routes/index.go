package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	fb "github.com/huandu/facebook/v2"
	"github.com/joho/godotenv"
	"github.com/jybp/ebay"
	"github.com/sosedoff/go-craigslist"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"os"
)

type Router struct {
}

func NewRouter() *Router {
	return &Router{}
}

type SearchParams struct {
	Value    string `form:"value"`
	Category string `form:"category"`
	Location string `form:"location"`
	PriceMin int    `form:"priceMin"`
	PriceMax int    `form:"priceMax"`
}

func (r Router) SearchAll(c *gin.Context) {
	var data SearchParams
	err := c.BindQuery(&data)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusBadRequest, gin.H{"message": "improper request"})
		return
	}

	res, err := SearchEbay(data.Value)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusBadRequest, gin.H{"message": "improper request"})
		return
	}

	craig, err := SearchCraigslist(data)

	if err != nil {
		c.Abort()
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}
	face, err := SearchFacebook(data)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Ebay":       res.ItemSummaries,
		"Craigslist": craig.Listings,
		"Facebook":   face,
	})
}

func SearchEbay(value string) (ebay.Search, error) {
	if err := godotenv.Load(".env"); err != nil {
		return ebay.Search{}, nil
	}

	cfg := clientcredentials.Config{
		ClientID:     os.Getenv("EBAY_CLIENT_KEY"),
		ClientSecret: os.Getenv("EBAY_CLIENT_SECRET"),
		TokenURL:     ebay.OAuth20SandboxEndpoint.TokenURL,
		Scopes:       []string{ebay.ScopeRoot /* your scopes */},
	}
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ebay.TokenSource(cfg.TokenSource(ctx)))
	client := ebay.NewSandboxClient(tc)

	// Get an item detail.
	return client.Buy.Browse.Search(ctx, ebay.OptBrowseSearch(value))
}

func SearchCraigslist(params SearchParams) (*craigslist.SearchResults, error) {
	opts := craigslist.SearchOptions{
		Category: params.Category, // cars+trucks
		Query:    params.Value,
		HasImage: true,
		MinPrice: params.PriceMin,
		MaxPrice: params.PriceMax,
	}

	return craigslist.Search(params.Location, opts)

}

func SearchFacebook(params SearchParams) (fb.Result, error) {
	res, err := fb.Get("/538744468", fb.Params{
		"fields":       "first_name",
		"access_token": "a-valid-access-token",
	})
	return res, err

}
