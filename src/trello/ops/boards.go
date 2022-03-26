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
	treq "tcli/trello"
	"tcli/trello/types"
)

func GetBoards(
	trello *treq.TrelloCtx,
	workspace types.Organization,
) ([]types.Board, error) {

	boardsEndpoint := treq.MakeEndpoint(
		fmt.Sprintf("/organizations/%s/boards", workspace.ID),
		[]string{"id", "name", "desc", "descData", "closed"},
	)
	boardsRaw, err := trello.ApiGet(boardsEndpoint)
	if err != nil {
		log.Printf("error obtaining orgs: %s\n", err)
		return nil, err
	}

	fmt.Println(string(boardsRaw))

	var boards []types.Board
	json.Unmarshal(boardsRaw, &boards)
	return boards, nil
}
