package main

import "time"

type Listing struct {
	Listing listing `json:"listing"`
	Medias  []media `json:"medias"`
}

type media struct {
	Url string `json:"url"`
}

type listing struct {
	Id              string    `json:"id"`
	Title           string    `json:"title"`
	TotalAreas      []int     `json:"totalAreas"`
	UsableAreas     []int     `json:"usableAreas"`
	Address         address   `json:"address"`
	Amenities       []string  `json:"amenities"`
	Bathrooms       []int     `json:"bathrooms"`
	Bedrooms        []int     `json:"bedrooms"`
	Suites          []int     `json:"suites"`
	Description     string    `json:"description"`
	Floors          []int     `json:"floors"`
	UnitFloor       []int     `json:"unitFloor"`
	ListingType     string    `json:"listingType"`
	UsageTypes      []string  `json:"usageTypes"`
	ParkingSpaces   []int     `json:"parkingSpaces"`
	Portal          string    `json:"portal"`
	PricingInfos    []pricing `json:"pricingInfos"`
	PropertyType    string    `json:"propertyType"`
	PublicationType string    `json:"publicationType"`
	UnitTypes       []string  `json:"unitTypes"`
	UnitsOnTheFloor int       `json:"unitsOnTheFloor"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type address struct {
	City         string  `json:"city"`
	Complement   string  `json:"complement"`
	LocationId   string  `json:"locationId"`
	Neighborhood string  `json:"neighborhood"`
	Point        latlong `json:"point"`
	State        string  `json:"state"`
	StateAcronym string  `json:"stateAcronym"`
	StreetNumber string  `json:"streetNumber"`
	ZipCode      string  `json:"zipCode"`
	Zone         string  `json:"zone"`
}

type latlong struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type pricing struct {
	BusinessType    string `json:"businessType"`
	MonthlyCondoFee string `json:"monthlyCondoFee"`
	Price           string `json:"price"`
	YearlyIptu      string `json:"yearlyIptu"`
}
