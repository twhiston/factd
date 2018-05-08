package cmd

import (
	"github.com/spf13/cobra"
	"github.com/twhiston/factd/lib/common/logging"
	"os"
	"strings"
	"text/template"
)

// tmplCmd represents the tmpl command
var tmplCmd = &cobra.Command{
	Use:   "tmpl",
	Short: "Create a template skeleton to extend factd",
	Long: `Subcommands of tmpl render various templates for items that can be used to extend factd.
As factd itself does not manage the plugins they must be added via the config.
In the standard factd setup this means adding them to the factdPlugins slice in cmd/root.go
Once this is done they will become available via the --include or --exclude options.
All plugins are on by default`,
}

// formatterCmd represents the formatter command
var formatterCmd = &cobra.Command{
	Use:   "formatter",
	Short: "Generate a skeleton formatter",
	Run: func(cmd *cobra.Command, args []string) {
		name, location := getNameAndLocation(cmd, "formatter")
		logging.Fatal(os.MkdirAll(location, os.ModePerm))
		render(name, location, formatterTmpl)
	},
}

// pluginCmd represents the plugin command
var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Generate a skeleton plugin",
	Run: func(cmd *cobra.Command, args []string) {
		name, location := getNameAndLocation(cmd, "plugin")
		location += name + "/"
		logging.Fatal(os.MkdirAll(location, os.ModePerm))
		render(name, location, pluginTmpl)
	},
}

func init() {
	rootCmd.AddCommand(tmplCmd)
	tmplCmd.AddCommand(formatterCmd)
	tmplCmd.AddCommand(pluginCmd)

	tmplCmd.PersistentFlags().String("name", "My", "The name of the Formatter/Plugin")
	tmplCmd.PersistentFlags().String("location", "", "Path to render the template. Defaults to ./lib/{{type}}")
}

func getNameAndLocation(cmd *cobra.Command, lType string) (string, string) {
	name, err := cmd.Flags().GetString("name")
	logging.Fatal(err)
	name = sanitizeName(name)

	location, err := cmd.Flags().GetString("location")
	logging.Fatal(err)
	location = sanitizeLocation(location, lType)
	return name, location
}

func sanitizeName(name string) string {
	return strings.Replace(name, " ", "", -1)
}

func sanitizeLocation(location, lType string) string {
	if location == "" {
		location = "lib/" + lType + "/"
	}
	if !strings.HasSuffix(location, "/") {
		location += "/"
	}
	return location
}

func render(name, location, tmplData string) {
	tmpl := template.New("factdTemplate")
	tmpl, err := tmpl.Parse(tmplData)
	logging.Fatal(err)

	f, err := os.Create(location + name + ".go")
	logging.Fatal(err)

	defer func(r *os.File) {
		err = r.Close()
		logging.Fatal(err)
	}(f)

	data := make(map[string]string)
	data["fname"] = strings.Title(name)
	err = tmpl.Execute(f, data)
	logging.HandleError(err)
}

var formatterTmpl = `package formatter

import (
	"bytes"
	"github.com/twhiston/factd/lib/common"
)

// {{ .fname }}Formatter prints-out facts in {{ .fname }} format
type {{ .fname }}Formatter struct {}

// New{{ .fname }}Formatter returns new {{ .fname }} formatter
func New{{ .fname }}Formatter() *{{ .fname }}Formatter {
	return &{{ .fname }}Formatter{}
}

// Format returns a formatted bytes buffer of facts in {{ .fname }} format
func (f *{{ .fname }}Formatter) Format(facts map[string]common.FactList) (*bytes.Buffer, error) {
	return bytes.NewBuffer([]byte{}), nil
}
`

var pluginTmpl = `package {{ .fname }}

import (
	"github.com/twhiston/factd/lib/common"
	"github.com/twhiston/factd/lib/plugin"
)

// The {{ .fname }} plugin provides information about SOMETHING
type {{ .fname }} struct{}

// Name returns the plugin printable name, also used as the map key in the master fact list
func (p *{{ .fname }}) Name() string {
	return plugin.GetPluginName(&p)
}

// Report writes a set (or subset) of facts to a channel
func (p *{{ .fname }}) Report(facts chan<- ReportedFact) {
	plugin.OneShotReport(p, facts)
}

// Facts gathers the actual fact data related to the plugin type
func (p *{{ .fname }}) Facts() (common.FactList, error) {
	facts := common.FactList{}
	return facts, nil
}
`
