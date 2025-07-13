package lang

import (
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/noppikinatta/ebitenginegamejam2025/asset"
)

func Switch() string {
	return txtProv.Switch()
}

func Text(key string) string {
	return ExecuteTemplate(key, nil)
}

func ExecuteTemplate(key string, data map[string]any) string {
	return txtProv.ExecuteTemplate(key, data)
}

const defaultLanguage = "english"

var txtProv *textProvider

type textProvider struct {
	CurrentLanguageIdx int
	Languages          []string
	Templates          map[string]map[string]*cachedTemplate
}

func (p *textProvider) SetDefault() {
	for i := range p.Languages {
		if p.Languages[i] == defaultLanguage {
			p.CurrentLanguageIdx = i
			return
		}
	}

	// fallback
	p.CurrentLanguageIdx = 0
}

func (p *textProvider) Switch() string {
	l := len(p.Languages)
	p.CurrentLanguageIdx = (p.CurrentLanguageIdx + 1) % l

	langName := p.Languages[p.CurrentLanguageIdx]
	return strings.ToUpper(langName[:1]) + langName[1:]
}

func (p *textProvider) ExecuteTemplate(key string, data map[string]any) string {
	lang := p.Languages[p.CurrentLanguageIdx]
	tmpl, ok := p.Templates[lang][key]
	if !ok {
		return fmt.Sprintf("NO_TMPL: %s, %v", key, data)
	}

	if data == nil {
		return tmpl.Text
	}

	return tmpl.Execute(data)
}

type cachedTemplate struct {
	Text string
	tmpl *template.Template
}

func (c *cachedTemplate) Execute(data map[string]any) string {
	t := c.template()

	buf := strings.Builder{}
	err := t.Execute(&buf, data)
	if err != nil {
		return fmt.Sprintf("ERR: %s, %v", c.Text, data)
	}

	return buf.String()
}

func (c *cachedTemplate) template() *template.Template {
	if c.tmpl != nil {
		return c.tmpl
	}

	t, err := template.New("t").Parse(c.Text)
	if err != nil {
		t = template.Must(template.New("fallback").Delims("[[", "]]").Parse("ERR:" + c.Text))
	}
	c.tmpl = t
	return t
}

func init() {
	langData := asset.LoadTemplates()
	langs := make([]string, 0, len(langData))
	langTmpls := make(map[string]map[string]*cachedTemplate, len(langData))

	for lang, dict := range langData {
		langs = append(langs, lang)
		tmpls := make(map[string]*cachedTemplate, len(dict))
		for key, tmplTxt := range dict {
			tmpls[key] = &cachedTemplate{Text: tmplTxt}
		}

		langTmpls[lang] = tmpls
	}

	sort.Strings(langs)

	txtProv = &textProvider{
		Languages: langs,
		Templates: langTmpls,
	}

	txtProv.SetDefault()
}
