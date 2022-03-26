/*
 * tcli - A Trello CLI client
 * Copyright (C) 2022  Joao Eduardo Luis <joao@wipwd.dev>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */

package ops

import (
	"encoding/json"
	"fmt"
	"log"
	"tcli/config"
	treq "tcli/trello"
	"tcli/trello/types"
)

func GetWorkspaces(
	trello *treq.TrelloCtx,
	config *config.Config,
) ([]types.Organization, error) {

	orgsEndpoint := treq.MakeEndpoint(
		fmt.Sprintf("/members/%s/organizations", config.ID),
		[]string{"id", "name", "displayName"},
	)
	orgsRaw, err := trello.ApiGet(orgsEndpoint)
	if err != nil {
		log.Printf("error obtaining orgs: %s\n", err)
		return nil, err
	}

	fmt.Println(string(orgsRaw))

	var orgs []types.Organization
	json.Unmarshal(orgsRaw, &orgs)
	return orgs, nil
}
