package main


import (
	"github.com/BurntSushi/toml"
	"html/template"
	"path/filepath"
	"strings"
	"os"
	"fmt"
)

type Config struct {
	StaticPath string
	TemplatePath string
	Pages []Page
	Layout template.Template
}

func (config *Config) ParseLayout() error {
	fmt.Println(config.ResolveTemplatePath("layout"))
	tmpl, err := template.
		New("layout").
		Funcs(config.Helpers()).
		ParseFiles(config.ResolveTemplatePath("layout"))


	if err != nil {
		return err
	}
	
	config.Layout = *tmpl

	return nil
}

func (config *Config) Helpers() template.FuncMap {
	return map[string]any {
		"imagePath": func (imageName string) string {
			return filepath.Join(config.StaticPath, "img", imageName)
		},
	}
}

func (config *Config) ResolveTemplatePath(name string) string {
	return filepath.Join(config.TemplatePath, name + ".html.tmpl")
}

type Page interface {
	TemplateName() string
	Title() []string
}

type GamePageEntry struct {
	Img string
	Title string
}

type GamePage struct {
	Entries []GamePageEntry
}

func (page *GamePage) Title() []string {
	return []string {"Game Page"}
}

func (page *GamePage) TemplateName() string {
	return "game_page"
}

type IndexPageEntry struct {
	Img string
	Width string
}

type IndexPage struct {
	Entries []IndexPageEntry
}

func (page *IndexPage) Title() []string {
	return []string {
		"It's Summer v",
		"s Winter",
	}
}

func (page *IndexPage) TemplateName() string {
	return "index"
}

type FanMusicPageEntry struct {
	Title string
	Img string
	Width string
	Tracks []string
}


type FanMusicPage struct {
	Entries []FanMusicPageEntry
}

func (page *FanMusicPage) Title() []string {
	return []string { "Fan Music Page" }
}

func (page *FanMusicPage) TemplateName() string {
	return "fan_music_page"
}

type UpdateEntryType string

const (
	UpdateEntryImage UpdateEntryType = "image"
	UpdateEntryPage UpdateEntryType = "page"
	UpdateEntryText UpdateEntryType = "text"
)

type UpdateLogPageEntry struct {
	Img string
	Page string
	Text string
	Type UpdateEntryType
}

type UpdateLogPageDateBlock struct {
	Date string
	Entries []UpdateLogPageEntry
}

type UpdateLogPage struct {
	Entries []UpdateLogPageDateBlock
}

func (page *UpdateLogPage) Title() []string {
	return []string { "Update Log" }
}

func (page *UpdateLogPage) TemplateName() string {
	return "update_log"
}

func getTemplateByName(config Config, name string) (string, error) {
	finalPath := config.ResolveTemplatePath(name)
	content, err := os.ReadFile(finalPath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func renderPage(config Config, page Page) (string, error) {
	builder := strings.Builder{}
	tmpl := template.Must(template.Must(config.Layout.Clone()). ParseFiles(config.ResolveTemplatePath(page.TemplateName())))

	err := tmpl.Execute(&builder, page)

	 if err != nil {
		 return "", err
	 }

	 return builder.String(), nil
 }

 func parseConfig(config Page) error {
	 configContent, err := os.ReadFile(config.TemplateName() + ".toml")
	 if err != nil {
		 return fmt.Errorf("Couldn't read config file %s\n", err)
	 }

	 _, err = toml.Decode(string(configContent), config)

	 if err != nil {
		 return err
	 }

	 return nil
 }

 func main() {
	 config := Config {
		 Pages: []Page {
			 &IndexPage{},
			 &UpdateLogPage{},
			 &FanMusicPage{},
			 &GamePage{},
		 },
	 }

	 configContent, err := os.ReadFile("config.toml")
	 if err != nil {
		 fmt.Printf("Abort: Couldn't read config file %s\n", err)
		 return
	 }

	 _, err = toml.Decode(string(configContent), &config)

	 if err != nil {
		 fmt.Printf("Abort: Couldn't parse config file %s\n", err)
		 return
	 }

	err = config.ParseLayout()

	if err != nil {
		fmt.Printf("Abort: Couldn't parse layout templates %s\n", err)
		return
	}

	 for _, page := range config.Pages {
		err := parseConfig(page)

		if err != nil {
			fmt.Printf("Abort: Couldn't parse config file %s %s", page.TemplateName() + ".toml", err)
			return
		}

		content, err := renderPage(config, page)

		if (err != nil) {
			fmt.Printf("Abort: Couldn't render %s %s\n", page.Title(), err)
			return
		}

		 err = os.WriteFile(page.TemplateName() + ".html", []byte(content), 0666)
	}
}
