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
	"fmt"
	"log"
	"os"
	"strings"

	"tcli/config"
	"tcli/trello"

	"github.com/AlecAivazis/survey/v2"
)

type TrelloPath struct {
	workspace trello.Workspace
	board     trello.Board
	card      string
}

func (p *TrelloPath) GetPath() string {
	var path []string
	if (trello.Workspace{}) != p.workspace {
		path = append(path, p.workspace.DisplayName)
	}
	if (trello.Board{}) != p.board {
		path = append(path, p.board.Name)
	}
	if p.card != "" {
		path = append(path, p.card)
	}
	return strings.Join(path, "/")
}

func (p *TrelloPath) HasWorkspace() bool {
	return (trello.Workspace{}) != p.workspace
}

func (p *TrelloPath) HasBoard() bool {
	return (trello.Board{}) != p.board
}

func (p *TrelloPath) HasCard() bool {
	return p.card != ""
}

var CmdsAvailable = []string{
	"help",
	"exit",
	"workspace",
	"board",
	"ls",
}

func CmdSuggestions(toComplete string) []string {

	log.Printf("cmd suggestions > input: %s\n", toComplete)

	if len(toComplete) == 0 {
		return CmdsAvailable
	}
	candidates := []string{}
	for _, cmd := range CmdsAvailable {
		if strings.Contains(cmd, toComplete) {
			candidates = append(candidates, cmd)
		}
	}
	return candidates
}

func main() {

	logFile, err := os.OpenFile(
		"tcli.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666,
	)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("test log")

	config, err := config.ReadConfig()
	if err != nil {
		fmt.Printf("error reading config: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(config)

	ctx := trello.Trello(config.ID, config.Key, config.Token)

	path := new(TrelloPath)

	exit := false
	for !exit {
		cmd := ""
		prompt := &survey.Input{
			Suggest: CmdSuggestions,
		}
		cliPath := path.GetPath()
		cliPathStr := ""
		if len(cliPath) > 0 {
			cliPathStr = fmt.Sprintf(" [ %s ]", cliPath)
		}
		survey.AskOne(
			prompt,
			&cmd,
			survey.WithIcons(func(icons *survey.IconSet) {
				icons.Question.Text = fmt.Sprintf("tcli%s >", cliPathStr)
			}),
			survey.WithValidator(survey.Required),
		)

		if cmd == "exit" {
			break
		} else if cmd == "workspace" {
			orgs, err := trello.GetWorkspaces(ctx, config)
			if err != nil {
				fmt.Printf("error: unable to obtain workspaces: %s\n", err)
				continue
			}
			var orgsNames []string
			nameToOrg := map[string]trello.Workspace{}
			for _, o := range orgs {
				orgsNames = append(orgsNames, o.DisplayName)
				nameToOrg[o.DisplayName] = o
			}
			prompt := &survey.Select{
				Options: orgsNames,
			}
			var selectedOrg string
			survey.AskOne(prompt, &selectedOrg)
			fmt.Printf("selected org: %s\n", nameToOrg[selectedOrg])
			path.workspace = nameToOrg[selectedOrg]

		} else if cmd == "board" {
			if !path.HasWorkspace() {
				fmt.Println("error: Workspace not selected")
				continue
			}
			boards, err := path.workspace.GetBoards(ctx)
			if err != nil {
				fmt.Printf(
					"error: unable to obtain boards for %s: %s",
					path.workspace.DisplayName,
					err,
				)
				continue
			}
			var boardNames []string
			nameToBoard := map[string]trello.Board{}
			for _, b := range boards {
				boardNames = append(boardNames, b.Name)
				nameToBoard[b.Name] = b
			}
			prompt := &survey.Select{Options: boardNames}
			var selectedBoard string
			survey.AskOne(prompt, &selectedBoard)
			fmt.Printf("selected board: %s\n", selectedBoard)
			board := nameToBoard[selectedBoard]
			if board.Closed {
				fmt.Printf("Board %s is closed", selectedBoard)
			} else {
				path.board = nameToBoard[selectedBoard]
			}

		} else if cmd == "ls" {
			if !(path.HasWorkspace() && path.HasBoard()) {
				fmt.Println("error: workspce and board not selected")
				continue
			}

			lists, err := path.board.GetLists(ctx)
			if err != nil {
				fmt.Printf(
					"error: unable to obtain lists for %s/%s: %s\n",
					path.workspace.DisplayName,
					path.board.Name,
					err,
				)
				continue
			}
			var listNames []string
			for _, l := range lists {
				listNames = append(listNames, l.Name)
			}
			prompt := &survey.Select{Options: listNames}
			var selectedList string
			survey.AskOne(prompt, &selectedList)

			fmt.Printf(
				"list cards on %s/%s/%s\n",
				path.workspace.DisplayName,
				path.board.Name,
				selectedList,
			)

			path.board.GetCards(ctx)
		}
	}
}
