package main

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ListLocations(local string, origin string) map[string]map[string]string {

	final_locations := map[string]map[string]string{}

	site_info := viper.Get("sites").(map[string]interface{})[origin].(map[string]interface{})
	portal := site_info["portal"].(string)
	api := site_info["api"]

	base_url := fmt.Sprintf("https://%s/v3/locations", api)

	headers := CreateHeaders(origin)

	query := map[string]interface{}{
		"q": local, "fields": "neighborhood", "portal": portal, "size": "6",
	}

	bytes_data := MakeRequest(base_url, query, headers)

	data := map[string]interface{}{}
	err := json.Unmarshal(bytes_data, &data)
	if err != nil {
		log.Error(fmt.Sprintf("Erro ao listar localizações %s: %v", local, err))
	}

	// Interface to map and get listings
	data = data["neighborhood"].(map[string]interface{})
	data = data["result"].(map[string]interface{})
	locations := data["locations"].([]interface{})

	for _, location := range locations {
		address := location.(map[string]interface{})["address"].(map[string]interface{})

		final_locations[address["locationId"].(string)] = map[string]string{
			"city":         address["city"].(string),
			"zone":         address["zone"].(string),
			"state":        address["state"].(string),
			"locationId":   address["locationId"].(string),
			"neighborhood": address["neighborhood"].(string),
			"stateAcronym": address["stateAcronym"].(string),
		}
	}

	return final_locations
}
