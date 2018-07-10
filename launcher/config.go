/*
Copyright (C) 2016-2018 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package launcher

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/andreas-jonsson/openwar/game"
	"github.com/andreas-jonsson/openwar/platform"
)

const configName = "config.json"

func createNewConfig(save bool) *game.Config {
	log.Println("Creating new config file.")

	cfg := &game.Config{
		ConfigVersion: game.ConfigVersion,
		Fullscreen:    false,
		Widescreen:    true,
		WC2Input:      true,
	}

	cfg.Debug.Map = "HUMAN01"
	cfg.Debug.Race = "Human"

	if save {
		saveConfig(cfg)
	}
	return cfg
}

func loadConfig() *game.Config {
	data, err := ioutil.ReadFile(platform.CfgRootJoin(configName))
	if err != nil {
		return createNewConfig(true)
	}

	var cfg game.Config
	if json.Unmarshal(data, &cfg) != nil || cfg.ConfigVersion != game.ConfigVersion {
		return createNewConfig(true)
	}

	return &cfg
}

func saveConfig(cfg *game.Config) {
	if data, err := json.Marshal(cfg); err == nil {
		ioutil.WriteFile(platform.CfgRootJoin(configName), data, 0755)
	} else {
		log.Println("Could not write config file.")
	}
}
