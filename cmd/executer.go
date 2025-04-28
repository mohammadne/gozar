package cmd

import (
	"errors"
	"log"

	"github.com/manifoldco/promptui"
	"github.com/mohammadne/hello-world/internal/entities"
	"github.com/mohammadne/hello-world/internal/generator"
	"github.com/mohammadne/hello-world/pkg/validator"
	"github.com/spf13/cobra"
)

const (
	defaultXrayVersion     = "v25.3.6"
	defaultOutputDirectory = "outputs"
)

var defaultPorts = map[entities.Protocol]map[entities.Machine]int{
	entities.NoTLS: {
		entities.Client: 10808,
		entities.Server: 1443,
	},
	entities.Reality: {
		entities.Client: 10809,
		entities.Server: 443,
	},
}

type Executer struct {
	// arguments
	outputDirectory string
	version         string

	// ask from the user
	serverAddress string
	protocol      entities.Protocol
}

func (e Executer) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "executer",
		Short: "execute the program and generate the configuration files based on given prompts",
		Run: func(_ *cobra.Command, _ []string) {
			e.selectProtocol()
			e.promptServerDomainOrIP()

			generator.Generator{
				Protocol:      e.protocol,
				Version:       e.version,
				ServerAddress: e.serverAddress,
				ServerPort:    defaultPorts[e.protocol][entities.Server],
				ClientPort:    defaultPorts[e.protocol][entities.Client],
			}.Run(e.outputDirectory)
		},
	}

	cmd.Flags().StringVar(&e.version, "xray-version", defaultXrayVersion, "The xray-core version")
	cmd.Flags().StringVar(&e.outputDirectory, "output-directory", defaultOutputDirectory, "The output directory for generated files")
	// cmd.Flags().IntVar(&executer.serverPort, "server-port", defaultRealityServerPort, "The port server exposed from")
	// cmd.Flags().IntVar(&executer.clientPort, "client-port", defaultRealityClientPort, "The port client exposed from")

	return cmd
}

func (executer *Executer) selectProtocol() {
	prompt := promptui.Select{
		Label: "Select which Xray protocol you want for your configuration",
		Items: []entities.Protocol{
			entities.NoTLS,
			entities.Reality,
		},
	}

	_, rawProtocol, err := prompt.Run()
	if err != nil {
		log.Fatalf("protocol prompt has been failed %v\n", err)
	} else if entities.ValidateProtocol(rawProtocol) != nil {
		log.Fatalf("invalid protocol value has been given %s \n%v\n", rawProtocol, err)
	}

	executer.protocol = entities.Protocol(rawProtocol)
}

func (executer *Executer) promptServerDomainOrIP() {
	validate := func(input string) error {
		if validator.ValidateIP(input) != nil && validator.ValidateDomain(input) != nil {
			return errors.New("invalid Domain or IP has been given")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Enter Domain or IPv4 address of the server",
		Validate: validate,
	}

	address, err := prompt.Run()
	if err != nil {
		log.Fatalf("protocol prompt has been failed %v\n", err)
	}

	executer.serverAddress = address
}
