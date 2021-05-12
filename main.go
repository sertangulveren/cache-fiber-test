package main

import (
	"cache-fiber-test/cmd"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"log"
	"net/http"
	"time"
)

type Campaign struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

var campaigns = map[string][]Campaign{
	"west-europe": []Campaign{
		{"WE-C1", "Buy a bike, get 1 liter of gasoline for free!"},
		{"WE-C2", "Buy a car, get driving school training for free!"},
		{"WE-C3", "Buy a jackboot, get 3 bikinis for free!"},
		{"WE-C4", "All products are 1% off!"},
	},
	"east-europe": []Campaign{
		{"EE-C1", "All products are 5% off!"},
	},
}

var queryCount int

func main() {
	app := fiber.New()
	app.Use(cache.New(cache.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			return getCacheKey(c)
		},
		Expiration: time.Second * 30,
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		results := FetchCampaignsFromDB(getCacheKey(c))
		return c.Status(fiber.StatusOK).JSON(results)
	})

	go cmd.Run()

	app.Listen(":3000")
}

func getCacheKey(c *fiber.Ctx) string {
	key := c.Query("region", "west-europe")
	return key
}

func FetchCampaignsFromDB(key string) []Campaign {
	queryCount++
	fmt.Println("Sorgulama Başlatıldı")
	fmt.Printf("Sorgulama Sayısı: %d\n", queryCount)
	results := campaigns[key]
	// 3 saniyelik bir fetch işlemi yapıyoruz. Amaç latency olsun.
	_, err := http.Get("https://deelay.me/3000/https://www.google.com")
	if err != nil {
		log.Fatalln(err)
	}
	return results
}
