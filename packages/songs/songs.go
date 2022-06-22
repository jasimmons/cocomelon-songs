package main

type SongList struct {
	season  int      `json:"season"`
	episode int      `json:"episode"`
	songs   []string `json:"songs"`
}

var (
	allSongs = map[string]map[string]SongList{
		"1": map[string]SongList{
			"1": SongList{
				season:  1,
				episode: 1,
				songs: []string{
					"song1",
					"song2",
					"song3",
				},
			},
		},
	}
)

func Main(args map[string]interface{}) map[string]interface{} {
	season, ok := args["season"].(string)
	if !ok {
		season = "1"
	}
	episode, ok := args["episode"].(string)
	if !ok {
		episode = "1"
	}

	if s, ok := allSongs[season]; ok {
		if list, ok := s[episode]; ok {
			return map[string]interface{}{
				"body": list.songs,
			}
		}
	}

	return map[string]interface{}{
		"errors": "did not work",
	}
}
