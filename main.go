package main

import (
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

		listings, err := FetchListings(DB, "vivareal", location, business_type, listing_type)

		if err != nil {
			c.JSON(400, err)
			return
		} else {
			c.JSON(200, listings)
		}

	})

	port := viper.GetString("PORT")
	r.Run(":" + port)
}
