package main

import (
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/krujos/cfcurl"
)

type AllServicesPlugin struct{}

func (c *AllServicesPlugin) Run(cliConnection plugin.CliConnection, args []string) {

	if args[0] == "all-services" {

		c.getRoutes(cliConnection)
	}

}

func (c *AllServicesPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "all-services",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 1,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "all-services",
				HelpText: "cf all-services",

				UsageDetails: plugin.Usage{
					Usage: "all-services\n   cf all-services",
				},
			},
		},
	}
}

func (c *AllServicesPlugin) getRoutes(cliConnection plugin.CliConnection, args ...string) {

	header := "organization,space_name,service_name,service_type,bound,guid,activity_last_30_days"
	fmt.Println(header)
	var nextURL interface{}
	nextURL = "/v2/organizations"

	for nextURL != nil {
		json, _ := cfcurl.Curl(cliConnection, nextURL.(string))
		resources := toJSONArray(json["resources"])

		for _, j := range resources {
			res := toJSONObject(j)
			entity := toJSONObject(res["entity"])
			spaceURL := entity["spaces_url"].(string)
			organization := entity["name"].(string)
			json, _ := cfcurl.Curl(cliConnection, spaceURL)
			resources := toJSONArray(json["resources"])

			for _, k := range resources {
				res := toJSONObject(k)
				entity := toJSONObject(res["entity"])
				spacename := entity["name"].(string)
				service_instances_url := entity["service_instances_url"].(string)

				json, _ := cfcurl.Curl(cliConnection, service_instances_url)
				resources := toJSONArray(json["resources"])

				for _, l := range resources {

					res := toJSONObject(l)
					entity := toJSONObject(res["entity"])

					service_name := entity["name"].(string)
					service_bindings_url := entity["service_bindings_url"].(string)
					service_url := entity["service_url"].(string)

					req, _ := cfcurl.Curl(cliConnection, service_url)
					entity1 := toJSONObject(req["entity"])
					service_type := entity1["label"].(string)

					json, _ := cfcurl.Curl(cliConnection, service_bindings_url)
					resources := toJSONArray(json["resources"])

					var guid string
					var bound string
					var event string
					bound = "no"
					for _, m := range resources {

						res := toJSONObject(m)
						entity := toJSONObject(res["entity"])
						app_url := entity["app_url"].(string)

						json, _ = cfcurl.Curl(cliConnection, app_url)
						metadata := toJSONObject(json["metadata"])
						guid = metadata["guid"].(string)
						bound = "yes"
						json1, _ := cfcurl.Curl(cliConnection, "/v2/events?q=actee:"+guid)
						total_results := fmt.Sprint(json1["total_results"])
						event = "no"
						if total_results != "0" {
							event = "yes"
							break
						}
					}

					var record1 interface{}
					record1 = organization + "," + spacename + "," + service_name + "," + service_type + "," + bound + "," + guid + "," + event
					fmt.Println(record1)
				}
			}
			nextURL = json["next_url"]
		}
	}
}

func main() {
	plugin.Start(new(AllServicesPlugin))
}

func toJSONArray(obj interface{}) []interface{} {
	return obj.([]interface{})
}

func toJSONObject(obj interface{}) map[string]interface{} {
	return obj.(map[string]interface{})
}
