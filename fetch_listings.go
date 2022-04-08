package main

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func FetchListings(
	origin string,
	location map[string]string,
	business_type string,
	listing_type string,
) []Listing {

	size := 24
	site_info := viper.Get("sites").(map[string]interface{})[origin].(map[string]interface{})
	base_url := fmt.Sprintf("https://%s/v2/listings", site_info["api"])

	max_page := viper.Get("max_page")
	max_page_int := max_page.(int64)

	headers := CreateHeaders(origin)
	query := createQuery(origin, location, business_type, listing_type, size)

	qtd_listings := qtdListings(base_url, query, headers)
	total_pages := int64(qtd_listings / size)
	if max_page_int <= 0 {
		max_page_int = total_pages
	} else {
		max_page_int = Min(max_page_int, total_pages)
	}

	log.Info(fmt.Sprintf("Getting %d/%d pages with %d listings from '%s'", max_page_int, total_pages, qtd_listings, origin))

	var all_listings []Listing

	for page := 1; page <= int(max_page_int); page++ {
		log.Info(fmt.Sprintf("Getting page %d from '%s'", page, origin))
		query["from"] = page * query["size"].(int)

		bytes_data := MakeRequest(base_url, query, headers)

		page_listings := getListings(bytes_data)

		for _, listing := range page_listings {
			all_listings = append(all_listings, listing)
		}

		if page < int(max_page_int) {
			time.Sleep(300 * time.Millisecond)
		}
	}

	return all_listings
}

func qtdListings(
	base_url string,
	query map[string]interface{},
	headers map[string]string,
) int {

	bytes_data := MakeRequest(base_url, query, headers)

	data := map[string]interface{}{}
	err := json.Unmarshal(bytes_data, &data)
	Check(err)

	data = data["search"].(map[string]interface{})

	qtd_listings := data["totalCount"].(float64)

	return int(qtd_listings)

}

func createQuery(
	origin string,
	location map[string]string,
	business_type string,
	listing_type string,
	size int,
) map[string]interface{} {
	return map[string]interface{}{
		"includeFields":       "search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount),expansion(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),nearby(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),page,fullUriFragments,developments(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),superPremium(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount)),owners(search(result(listings(listing(displayAddressType,amenities,usableAreas,constructionStatus,listingType,description,title,stamps,createdAt,floors,unitTypes,nonActivationReason,providerId,propertyType,unitSubTypes,unitsOnTheFloor,legacyId,id,portal,unitFloor,parkingSpaces,updatedAt,address,suites,publicationType,externalId,bathrooms,usageTypes,totalAreas,advertiserId,advertiserContact,whatsappNumber,bedrooms,acceptExchange,pricingInfos,showPrice,resale,buildings,capacityLimit,status),account(id,name,logoUrl,licenseNumber,showAddress,legacyVivarealId,legacyZapId,minisite),medias,accountLink,link)),totalCount))",
		"addressNeighborhood": location["neighborhood"],
		"addressLocationId":   location["locationId"],
		"addressState":        location["state"],
		"addressCity":         location["city"],
		"addressZone":         location["zone"],
		"listingType":         listing_type,
		"business":            business_type,
		"usageTypes":          "RESIDENTIAL",
		"categoryPage":        "RESULT",
		"size":                size,
		"from":                24,
	}
}

func getListings(bytes_data []byte) []Listing {
	var page_listings []Listing

	data := map[string]interface{}{}
	err := json.Unmarshal(bytes_data, &data)
	Check(err)

	// Interface to map and get listings
	data = data["search"].(map[string]interface{})
	data = data["result"].(map[string]interface{})
	listings_page := data["listings"].([]interface{})

	// Slice of listings to Struct
	jsonString, err := json.Marshal(listings_page)
	Check(err)
	// os.WriteFile("test.json", jsonString, 0666)

	json.Unmarshal(jsonString, &page_listings)

	return page_listings
}
