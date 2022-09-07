/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"black-key/config"

	"black-key/pkg/router"

	"github.com/bamzi/jobrunner"
	"github.com/gin-gonic/gin"
)

func main() {
	// cmd.Execute()
	Run()
}

func Run() error {
	jobrunner.Start()
	if env := config.Cfg.GetString("APP_ENV"); env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	route := gin.Default()
	router.SetRouter(route)
	port := config.Cfg.GetString("PORT")
	if port != "" {
		return route.Run(":" + port)
	}
	return route.Run()
}
