package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/hashicorp/go-envparse"
	"google.golang.org/api/option"
)

func main() {
	flags := parseFlags()
	envFile, err := os.ReadFile(flags.input)
	if err != nil {
		log.Fatalf("failed to read input file: %v", err)
	}
	envs, err := envparse.Parse(bytes.NewReader(envFile))
	if err != nil {
		log.Fatalf("failed to parse input file: %v", err)
	}

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx, option.WithCredentialsFile(flags.credential))
	if err != nil {
		log.Fatalf("failed to setup client (-credential may be required): %v", err)
	}
	defer client.Close()

	for key, value := range envs {
		if !strings.HasPrefix(value, "projects/") {
			continue
		}
		access, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
			Name: value,
		})
		if err != nil {
			log.Fatalf("failed to get secret (%s): %v", value, err)
		}

		escaped := string(access.Payload.Data)
		if flags.removeWhitespace {
			escaped = strings.ReplaceAll(escaped, "'", "'\"'\"'")
			escaped = strings.ReplaceAll(escaped, "\n", "\\n")
			escaped = strings.ReplaceAll(escaped, "\t", "\\t")
		}
		envs[key] = escaped
	}

	output := bytes.NewBuffer(nil)
	for key, value := range envs {
		output.WriteString(fmt.Sprintf("%s='%s'\n", key, value))
	}
	outputBytes := output.Bytes()

	file := os.Stdout
	if flags.output != "" {
		file, err = os.Create(flags.output)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
	}
	writer := bufio.NewWriter(file)
	_, err = writer.Write(outputBytes)
	if err != nil {
		log.Fatalf("failed to write output: %v", err)
	}
	writer.Flush()
}

type flags struct {
	output           string
	input            string
	credential       string
	removeWhitespace bool
}

func parseFlags() flags {
	defaultCredential := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	output := flag.String("output", "", "output file")
	help := flag.Bool("help", false, "show help")
	credential := flag.String("credential", defaultCredential, "gcp credential file")
	removeWhitespace := flag.Bool("remove-whitespace", false, "remove whitespaces {\\n,\\t}")

	flag.Parse()
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] <input-file>\n", os.Args[0])
		fmt.Println("Note: <input-file> is a required positional argument.")
		flag.PrintDefaults()
	}

	if flag.NArg() != 1 || *help {
		flag.Usage()
		os.Exit(1)
	}

	inputFilename := flag.Args()[0]
	return flags{
		output:           *output,
		input:            inputFilename,
		credential:       *credential,
		removeWhitespace: *removeWhitespace,
	}
}
