// generates the fixtures/fonts.html for testing the fonts in docker.
// Use the google translate to translate "test" into all the languages, print the result into a html page.
// By reviewing the generated pdf we can find out what font is missing for a specific language.

package main

import (
	"fmt"
	"strings"

	"github.com/ysmood/kit"
	"github.com/ysmood/rod"
	"github.com/ysmood/rod/lib/launcher"
)

func main() {
	url := launcher.New().Headless(true).Launch()
	b := rod.New().ControlURL(url).Connect()
	defer b.Close()

	p := b.Page("https://translate.google.com/")

	p.Element("#source").Input("Test the google translate.")

	if p.Has(".tlid-dismiss-button") {
		p.Element(".tlid-dismiss-button").Click()
	}

	showList := p.Element(".tlid-open-target-language-list")
	list := p.Elements(".language-list:nth-child(2) .language_list_section:nth-child(2) .language_list_item_language_name")

	html := ""

	for _, lang := range list {
		showList.Click()
		wait := p.WaitRequestIdle()
		lang.Click()
		wait()
		name := lang.Text()
		result := p.Element(".tlid-translation").Text()
		for strings.Contains(result, "...") {
			kit.Sleep(0.1)
			result = p.Element(".tlid-translation").Text()
		}
		kit.Log(name, result)
		html += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", name, result)
	}

	html = fmt.Sprintf(`<html>
		<p style="font-family: serif;">
			This file is generated by <code>"fixtures/cmd/gen-fonts"</code>
		</p>
		<table>
		%s
		<table></html>`,
		html,
	)

	kit.E(kit.OutputFile("fixtures/fonts.html", html, nil))
}
