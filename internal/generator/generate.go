package generator

import (
	"embed"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/mohammadne/hello-world/internal/entities"
	"github.com/mohammadne/hello-world/pkg/cryptography"
)

type Generator struct {
	Protocol      entities.Protocol
	Version       string
	ServerAddress string
	ServerPort    int
	ClientPort    int

	serverOutputDirectory string
	clientOutputDirectory string
	keyPair               *cryptography.KeyPair
	uuid                  string
}

const (
	TemplatesDirectory = "templates/"
)

//go:embed templates
var templates embed.FS

func (generator Generator) Run(outputDirectory string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error findout pwd: %v\n", err)
	}
	outputDirectory = fmt.Sprintf("%s/%s/", wd, outputDirectory)

	if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
		log.Fatalf("error creating directory %s, %v\n", outputDirectory, err)
	}

	generator.serverOutputDirectory = outputDirectory + "/server/"
	generator.clientOutputDirectory = outputDirectory + "/client/"

	generator.uuid = cryptography.GenerateUUID()

	keyPair, err := cryptography.GenerateCurve25519Keys()
	if err != nil {
		log.Fatalf("error creating Curve25519 key pairs %v\n", err)
	}
	generator.keyPair = keyPair

	for _, machine := range []entities.Machine{entities.Server, entities.Client} {
		path := generator.outputPath(machine)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatalf("error creating directory %s, %v\n", path, err)
		}

		generator.generateDockerfile(machine)
		generator.generateConfig(machine)
		generator.generateCompose(machine)
	}
}

func (generator *Generator) outputPath(machine entities.Machine) string {
	if machine == entities.Server {
		return generator.serverOutputDirectory
	}
	return generator.clientOutputDirectory
}

const (
	DockerfileTemplateFile = "Dockerfile.tmpl"
	DockerfileOutputFile   = "Dockerfile"
)

func (generator *Generator) generateDockerfile(machine entities.Machine) {
	templateFile := TemplatesDirectory + DockerfileTemplateFile
	tmpl, err := template.ParseFS(templates, templateFile)
	if err != nil {
		panic(err)
	}

	var outputFile *os.File
	outputFile, err = os.Create(generator.outputPath(machine) + DockerfileOutputFile)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	values := struct{ XrayVersion string }{XrayVersion: generator.Version}

	if err = tmpl.Execute(outputFile, values); err != nil {
		panic(err)
	}
}

const (
	ServerTemplateConfigFile = "%s-server.json.tmpl"
	ClientTemplateConfigFile = "%s-client.json.tmpl"
	OutputConfigFile         = "%s.json"
)

func (g *Generator) generateConfig(machine entities.Machine) {
	var templateFile string
	if machine == entities.Server {
		templateFile = TemplatesDirectory + fmt.Sprintf(ServerTemplateConfigFile, g.Protocol)
	} else {
		templateFile = TemplatesDirectory + fmt.Sprintf(ClientTemplateConfigFile, g.Protocol)
	}

	tmpl, err := template.ParseFS(templates, templateFile)
	if err != nil {
		panic(err)
	}

	var outputFile *os.File
	outputFile, err = os.Create(g.outputPath(machine) + fmt.Sprintf(OutputConfigFile, g.Protocol))
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	var values any
	if machine == entities.Server {
		values = struct {
			Port       int
			UUID       string
			PrivateKey string
		}{
			Port:       g.ServerPort,
			UUID:       g.uuid,
			PrivateKey: g.keyPair.PrivateKey,
		}
	} else {
		values = struct {
			ClientPort    int
			ServerAddress string
			ServerPort    int
			UUID          string
			PublicKey     string
		}{
			ClientPort:    g.ClientPort,
			ServerAddress: g.ServerAddress,
			ServerPort:    g.ServerPort,
			UUID:          g.uuid,
			PublicKey:     g.keyPair.PublicKey,
		}
	}

	if err = tmpl.Execute(outputFile, values); err != nil {
		panic(err)
	}
}

const (
	ComposeTemplateFile = "compose.yml.tmpl"
	ComposeOutputFile   = "compose-%s.yml"
)

type ComposeValues struct {
	Port     int
	Machine  string
	Protocol string
}

func (g *Generator) generateCompose(machine entities.Machine) {
	templateFile := TemplatesDirectory + ComposeTemplateFile
	tmpl, err := template.ParseFS(templates, templateFile)
	if err != nil {
		panic(err)
	}

	var values ComposeValues
	values.Machine = string(machine)
	values.Protocol = string(g.Protocol)
	if machine == entities.Server {
		values.Port = g.ServerPort
	} else {
		values.Port = g.ClientPort
	}

	var outputFile *os.File
	outputFile, err = os.Create(g.outputPath(machine) + fmt.Sprintf(ComposeOutputFile, g.Protocol))
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	if err = tmpl.Execute(outputFile, values); err != nil {
		panic(err)
	}
}
