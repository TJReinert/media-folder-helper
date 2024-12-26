package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Jisin0/filmigo/imdb"
	"github.com/tjreinert/media-folder-helper/omdb"
)

type MediaType string

const (
	Movie MediaType = "movie"
	TV    MediaType = "tvseries"
)

func (m MediaType) String() string {
	switch m {
	case Movie:
		return "Movie"
	case TV:
		return "Tv Series"
	default:
		panic("Unknown media type")
	}
}

func mediaTypeFromString(input string) MediaType {
	switch strings.ToLower(input) {
	case "movie":
		return Movie
	case "TV Series":
		return TV
	default:
		panic("Unknown media type")
	}
}

type Model struct {
	mediaType       MediaType
	title           string
	contentOptions  []*imdb.SearchResult
	selectedContent *imdb.SearchResult
	year            int
	seasons         int
}

func promptType(model *Model) {
	var selection string
	survey.AskOne(&survey.Select{
		Message: "Select the type:",
		Options: []string{Movie.String(), TV.String()},
	}, &selection)

	switch selection {
	case Movie.String():
		model.mediaType = Movie
	case TV.String():
		model.mediaType = TV
	default:
		panic("Unknown media type")
	}
}

func promptTitle(model *Model) string {
	survey.AskOne(&survey.Input{
		Message: "Enter the title of the movie or TV show:",
	}, &model.title)
	return model.title
}

func searchForContent(model *Model, imdbClient *imdb.ImdbClient, omdbClient *omdb.Client) {
	results, err := imdbClient.SearchTitles(model.title, &imdb.SearchConfigs{IncludeVideos: false})
	if err != nil || len(results.Results) == 0 {
		fmt.Println("[ERROR] Error while searching for titles, taking user input.")
	} else {
		optionsMap := make(map[string]*imdb.SearchResult)
		options := make([]string, 0)

		for _, entry := range results.Results {
			if strings.EqualFold(entry.Category, string(model.mediaType)) {
				optionsMap[entry.Title] = entry
				options = append(options, entry.Title)
			}
		}

		options = append(options, "None of the above")
		optionsMap["None of the above"] = nil

		var selected string
		survey.AskOne(&survey.Select{
			Message: "Select:",
			Options: options,
		}, &selected)

		selectedResult := optionsMap[selected]

		if selectedResult == nil {
			fmt.Println("[ERROR] Error while re-mapping title selection, taking user input.")
		} else {
			model.title = selectedResult.Title
			model.year = selectedResult.Year
			if model.mediaType == TV {
				if omdbClient != nil {
					resp, _ := omdbClient.GetById(selectedResult.ID)
					if resp != nil {
						parsed, _ := strconv.Atoi(resp.TotalSeasons)
						model.seasons = parsed
					}
				}
			}
		}
	}
}

func isNumber(val interface{}) error {
	strVal, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid input type")
	}

	_, err := strconv.Atoi(strVal)
	if err != nil {
		return fmt.Errorf("please enter a valid number")
	}
	return nil
}

func promptForYear(model *Model) {
	if model.year == 0 {
		survey.AskOne(&survey.Input{
			Message: "Year:",
		}, &model.year, survey.WithValidator(isNumber))
	} else {
		survey.AskOne(&survey.Input{
			Message: "Year:",
			Default: strconv.Itoa(model.year),
		}, &model.year, survey.WithValidator(isNumber))
	}
}

func promptForSeasons(model *Model) {
	if model.seasons == 0 {
		survey.AskOne(&survey.Input{
			Message: "Number of Seasons:",
		}, &model.seasons, survey.WithValidator(isNumber))
	} else {
		survey.AskOne(&survey.Input{
			Message: "Number of Seasons:",
			Default: strconv.Itoa(model.seasons),
		}, &model.seasons, survey.WithValidator(isNumber))
	}
}

func formatTitle(model *Model) string {
	var replacer = strings.NewReplacer(":", " - ")
	title := replacer.Replace(model.title)

	return title + " (" + strconv.Itoa(model.year) + ")"
}

func main() {
	model := Model{}

	promptType(&model)
	promptTitle(&model)
	imdbClient := imdb.NewClient()
	var omdbClient *omdb.Client
	omdbApiKey := os.Getenv("OMDB_API_KEY")
	if omdbApiKey != "" {
		omdbClient = omdb.NewClient(omdbApiKey)
	}
	searchForContent(&model, imdbClient, omdbClient)
	promptForYear(&model)
	if model.mediaType == TV {
		promptForSeasons(&model)
	}

	if model.title == "" || model.year == 0 || model.mediaType == "" {
		fmt.Println("Error: Missing or invalid required arguments.")
		os.Exit(1)
	}

	baseDir := filepath.Join(".", formatTitle(&model))

	if model.mediaType == Movie {
		createFolder(baseDir)
	} else if model.mediaType == TV {
		createFolder(baseDir)
		for i := 1; i <= model.seasons; i++ {
			seasonFolder := filepath.Join(baseDir, fmt.Sprintf("Season %02d", i))
			createFolder(seasonFolder)
		}
	}

	fmt.Println("Folders created successfully!")
}

func createFolder(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Printf("Error creating folder '%s': %v\n", path, err)
		os.Exit(1)
	}
}
