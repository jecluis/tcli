/*
 * tcli - A Trello CLI client
 * Copyright (C) 2022  Joao Eduardo Luis <joao@wipwd.dev>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"tcli/config"
	treq "tcli/trello"
	"tcli/trello/types"
)

func main() {

	config, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("error reading config: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(config)

	trello := treq.Trello(config.ID, config.Key, config.Token)

	// fmt.Println("--- Organizations ---")
	orgsEndpoint := treq.MakeEndpoint(
		fmt.Sprintf("/members/%s/organizations", config.ID),
		[]string{"id", "name", "displayName"},
	)
	orgsRaw, err := trello.ApiGet(orgsEndpoint)
	if err != nil {
		fmt.Printf("error obtaining orgs: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(orgsRaw))

	var orgs []types.Organization
	json.Unmarshal(orgsRaw, &orgs)
	fmt.Println(orgs)
}
