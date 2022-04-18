package main

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	LoadSettings()
	DB = Connect()
	DB.AutoMigrate(&Property{})

	business_type_values := map[string]bool{"RENTAL": true, "SALE": true}
	listing_type_values := map[string]bool{"DEVELOPMENT": true, "USED": true}

	mapSites := viper.GetStringMap("sites")
	todosSites := GetKeys(mapSites)

	r := gin.Default()

	r.GET("/locations/:location", func(c *gin.Context) {
		location := c.Param("location")
		locations := ListLocations(location, "vivareal")

		c.JSON(200, locations)
	})

	r.GET("listings/:business_type/:listing_type/:city/:locationId/:neighborhood/:state/:stateAcronym/:zone", func(c *gin.Context) {
		location := map[string]string{
			"city":         c.Param("city"),
			"zone":         c.Param("zone"),
			"state":        c.Param("state"),
			"locationId":   c.Param("locationId"),
			"neighborhood": c.Param("neighborhood"),
			"stateAcronym": c.Param("stateAcronym"),
		}

		business_type := c.Param("business_type")
		listing_type := c.Param("listing_type")

		if _, ok := business_type_values[business_type]; !ok {
			c.String(400, "Business types allowed ['RENTAL', 'SALE']")
			return
		}

		if _, ok := listing_type_values[listing_type]; !ok {
			c.String(400, "Listing types allowed ['DEVELOPMENT', 'USED']")
			return
		}

		var errs []error
		var oks []string
		var wg sync.WaitGroup
		channel_err := make(chan []error)
		channel_ok := make(chan string)

		for _, site := range todosSites {
			wg.Add(1)

			go func(s string, w sync.WaitGroup) {
				defer w.Done()
				msg, err := FetchListings(DB, s, location, business_type, listing_type)
				if err != nil {
					channel_err <- err
				} else {
					channel_ok <- msg
				}
			}(site, wg)
		}
		wg.Wait()

		for err := range channel_err {
			if err != nil {
				errs = append(errs, err...)
			}
		}

		for ok := range channel_ok {
			if ok != "" {
				oks = append(oks, ok)
			}
		}

		if errs != nil {
			c.JSON(400, errs)
			return
		} else {
			c.JSON(200, fmt.Sprintf("Salvo com sucesso: %v", oks))
		}

	})

	port := viper.GetString("PORT")
	r.Run(":" + port)
}
