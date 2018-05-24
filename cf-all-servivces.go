package main

import (
	"fmt"
	"github.com/krujos/cfcurl"
	"github.com/cloudfoundry/cli/plugin"
)

type AllRoutesPlugin struct{}

func (c *AllRoutesPlugin) Run(cliConnection plugin.CliConnection, args []string) {

	if args[0] == "all-routes" {

    c.getRoutes(cliConnection)

}

}

func (c *AllRoutesPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "all-routes",
		Version: plugin.VersionType{
			Major: 2,
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
				Name:     "all-routes",
				HelpText: "cf all-routes",

				UsageDetails: plugin.Usage{
					Usage: "all-routes\n   cf all-routes",
				},
			},
		},
	}
}

func (c *AllRoutesPlugin) getRoutes(cliConnection plugin.CliConnection, args ...string) {

header:="space_name,service_name,service_url,routes_url,service_created_at,service_updated_at,app_name,package_updated_at"
fmt.Println(header)

var nextURL interface{}
	nextURL = "/v2/organizations/6c3d5da1-4e4b-48cb-8d6f-d2cc89708f64/spaces"
	for nextURL != nil {

  json, _ := cfcurl.Curl(cliConnection, nextURL.(string))
  resources := toJSONArray(json["resources"])

	  for _, i := range resources {
	    res := toJSONObject(i)
	    entity := toJSONObject(res["entity"])
		space_name := entity["name"].(string)
		service_instances_url := entity["service_instances_url"].(string)

		json, _ = cfcurl.Curl(cliConnection, service_instances_url)
		resources := toJSONArray(json["resources"])
			for _, j := range resources {
			 res := toJSONObject(j)
			 entity := toJSONObject(res["entity"])
			 service_name := entity["name"].(string)

			 metadata := toJSONObject(res["metadata"])
			 created_at := metadata["created_at"].(string)
			 updated_at := metadata["updated_at"].(string)
			 service_url := metadata["url"].(string)
			 routes_url := entity["routes_url"].(string)

			 service_bindings_url := entity["service_bindings_url"]
			 
			 json, _ = cfcurl.Curl(cliConnection, service_bindings_url.(string))
			 total_results := fmt.Sprint(json["total_results"])
			 
			if total_results == "0" {

			 var record1 interface{}
			 record1 = space_name+","+service_name+","+service_url+","+routes_url+","+created_at+","+updated_at
			 fmt.Println(record1)
			}


			 resources = toJSONArray(json["resources"])
			 for _, k := range resources {
			 	res := toJSONObject(k)			 
				entity := toJSONObject(res["entity"])
			 	app_url := entity["app_url"].(string)

					json, _ = cfcurl.Curl(cliConnection, app_url)		 
					entity1 := toJSONObject(json["entity"])
					app_name := entity1["name"].(string)
					package_updated_at := entity1["package_updated_at"].(string)


			 var record interface{}
			 record = space_name+","+service_name+","+service_url+","+routes_url+","+created_at+","+updated_at+","+app_name+","+package_updated_at
			 fmt.Println(record)
			 }


			 //json = cfcurl.Curl(cliConnection, service_instance_url)
			 //resources := toJSONArray(json["resources"])

			}



		}
	nextURL = json["next_url"]
}
}

func main() {
	plugin.Start(new(AllRoutesPlugin))
}

func toJSONArray(obj interface{}) []interface{} {
	return obj.([]interface{})
}

func toJSONObject(obj interface{}) map[string]interface{} {
	return obj.(map[string]interface{})
}
