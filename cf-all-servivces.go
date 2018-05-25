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
				HelpText: "cf all-routes",

				UsageDetails: plugin.Usage{
					Usage: "all-routes\n   cf all-routes",
				},
			},
		},
	}
}

func (c *AllServicesPlugin) getRoutes(cliConnection plugin.CliConnection, args ...string) {

	header := "space_name,service_name,app_name,guid,activity"
	fmt.Println(header)

	var nextURL interface{}
	nextURL = "/v2/organizations/6c3d5da1-4e4b-48cb-8d6f-d2cc89708f64/spaces"
	for nextURL != nil {

		json, _ := cfcurl.Curl(cliConnection, nextURL.(string))
		resources := toJSONArray(json["resources"])

		// Loop in all space for service instances url in order to list the services for each space

		for _, i := range resources {
			res := toJSONObject(i)
			entity := toJSONObject(res["entity"])

			spacename := entity["name"].(string)
			service_instances_url := entity["service_instances_url"].(string)

			json, _ := cfcurl.Curl(cliConnection, service_instances_url)
			resources := toJSONArray(json["resources"])

			// loop in each services of all space
			for _, i = range resources {

				res := toJSONObject(i)
				entity := toJSONObject(res["entity"])
				//metadata := toJSONObject(res["metadata"])

				service_name := entity["name"].(string)
				//service_url := metadata["url"].(string)
				service_bindings_url := entity["service_bindings_url"].(string)

				json, _ := cfcurl.Curl(cliConnection, service_bindings_url)
				total_results := fmt.Sprint(json["total_results"])

				if total_results == "0" {

					var record interface{}
					record = spacename + "," + service_name
					fmt.Println(record)
				}

				resources := toJSONArray(json["resources"])

				for _, i = range resources {

					res := toJSONObject(i)
					entity := toJSONObject(res["entity"])
					app_url := entity["app_url"].(string)

					json, _ = cfcurl.Curl(cliConnection, app_url)
					entity1 := toJSONObject(json["entity"])
					metadata := toJSONObject(json["metadata"])
					guid := metadata["guid"].(string)
					app_name := entity1["name"].(string)
					json1, _ := cfcurl.Curl(cliConnection, "/v2/events?q=actee:"+guid)
					total_results := fmt.Sprint(json1["total_results"])
					event := "yes"
					if total_results == "0" {
						event = "No"
					}

					var record1 interface{}
					record1 = spacename + "," + service_name + "," + app_name + "," + guid + "," + event
					fmt.Println(record1)
				}
			}
		}
		nextURL = json["next_url"]

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
