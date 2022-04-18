package main

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func FetchListings(
	DB *gorm.DB,
	origin string,
	location map[string]string,
	business_type string,
	listing_type string,
) (string, []error) {
	var all_listings []Property
	var page_listings Property
	var errors []error

	size := 24
	site_info := viper.Get("sites").(map[string]interface{})[origin].(map[string]interface{})
	base_url := fmt.Sprintf("https://%s/v2/listings", site_info["api"])

	max_page := viper.GetInt64("max_page")

	headers := CreateHeaders(origin)
	query := createQuery(origin, location, business_type, listing_type, size)

	qtd_listings, err := qtdListings(base_url, query, headers)
	if err != nil {
		return "", []error{err}
	}
	total_pages := int64(qtd_listings / size)
	if max_page <= 0 {
		max_page = total_pages
	} else {
		max_page = Min(max_page, total_pages)
	}

	log.Info(fmt.Sprintf("Getting %d/%d pages with %d listings from '%s'", max_page, total_pages, qtd_listings, origin))

	for page := 1; page <= int(max_page); page++ {
		log.Info(fmt.Sprintf("Getting page %d from '%s'", page, origin))
		query["from"] = page * query["size"].(int)

		bytes_data := MakeRequest(base_url, query, headers)

		listings, err := page_listings.Unmarshal(bytes_data, origin, business_type)
		if err != nil {
			log.Error(fmt.Sprintf("Erro ao parse da página %d do site '%s': %v", page, origin, err))
			errors = append(errors, err)
		}

		all_listings = append(all_listings, listings...)

		if page < int(max_page) {
			time.Sleep(300 * time.Millisecond)
		}
	}

	result := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "origin"}, {Name: "url"}, {Name: "business_type"}},
		UpdateAll: true,
	}).CreateInBatches(all_listings, 100)

	if result.Error != nil {
		log.Error(result.Error)
		errors = append(errors, result.Error)
	}

	return fmt.Sprintf("Saved %d pages from '%s'", max_page, origin), errors
}

func qtdListings(
	base_url string,
	query map[string]interface{},
	headers map[string]string,
) (int, error) {

	bytes_data := MakeRequest(base_url, query, headers)

	data := map[string]interface{}{}
	err := json.Unmarshal(bytes_data, &data)
	if err != nil {
		err := fmt.Errorf(fmt.Sprintf("erro ao buscar a quantidade de propriedades da página '%s' '%v' '%v': %v", base_url, query, bytes_data, err))
		log.Error(err)
		return 0, err
	}

	if !Contains(GetKeys(data), "search") {
		err := fmt.Errorf("not found search listings '%v' from '%s' '%v'", data, base_url, query)
		log.Error(err)
		return 0, err
	}

	data = data["search"].(map[string]interface{})

	qtd_listings := data["totalCount"].(float64)

	return int(qtd_listings), nil

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
