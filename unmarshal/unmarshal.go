package unmarshal

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/data"
	"io"
	"net/http"
	"strings"
)

func Fetch(link string, id int) {
	artist, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
	}

	artistData, err := io.ReadAll(artist.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(artistData, &data.Artists)
	if err != nil {
		fmt.Println(err)
	}
	if id != 0 {
		artist2 := &data.Artists[id-1]
		var relations data.Relations
		rRelations, err := http.Get(artist2.RelationsURL)
		if err != nil {
			fmt.Println(err)
		}
		relationsData, err := io.ReadAll(rRelations.Body)
		if err != nil {
			fmt.Println(err)
		}
		err = json.Unmarshal(relationsData, &relations)
		if err != nil {
			fmt.Println(err)
		}

		newMap := make(map[string][]string)
		for key, val := range relations.DatesLocations {
			newKey := key
			newKey = strings.Replace(newKey, "-", ", ", -1)
			newKey = strings.Replace(newKey, "_", " ", -1)
			newKey = strings.Title(strings.ToLower(newKey))
			newKey = strings.Replace(newKey, "Usa", "USA", -1)
			newKey = strings.Replace(newKey, "Uk", "UK", -1)
			newMap[newKey] = val
		}
		artist2.Relations = newMap
	}
	err = artist.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
}
