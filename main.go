package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/PullRequestInc/go-gpt3"
)

type payload struct {
	Text string `json:"text"`
}

type response struct {
	EOF   bool   `json:"eof"`
	Error string `json:"error"`
	Text  string `json:"text"`
}

func doJson(client gpt3.Client, r io.Reader, w io.Writer) error {
	enc := json.NewEncoder(w)
	dec := json.NewDecoder(r)
	for {
		var p payload
		err := dec.Decode(&p)
		if err != nil {
			return err
		}
		err = client.CompletionStreamWithEngine(
			context.Background(),
			gpt3.TextDavinci003Engine,
			gpt3.CompletionRequest{
				Prompt: []string{
					p.Text,
				},
				MaxTokens:   gpt3.IntPtr(3000),
				Temperature: gpt3.Float32Ptr(0),
			}, func(resp *gpt3.CompletionResponse) {
				enc.Encode(response{EOF: false, Text: resp.Choices[0].Text})
			},
		)
		if err != nil {
			err = enc.Encode(response{Error: err.Error()})
			if err != nil {
				return err
			}
			continue
		}
		err = enc.Encode(response{EOF: true})
		if err != nil {
			return err
		}
	}
}

func main() {
	var j bool
	flag.BoolVar(&j, "json", false, "json input/output")
	flag.Parse()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing OPENAI_API_KEY env variable")
	}

	client := gpt3.NewClient(apiKey)

	if j {
		log.Fatal(doJson(client, os.Stdin, os.Stdout))
	}

  args := os.Args
  if len(args) == 2 {
		text := "plain text command to" + args[1]
		err := client.CompletionStreamWithEngine(
			context.Background(),
			gpt3.TextDavinci003Engine,
			gpt3.CompletionRequest{
				Prompt: []string{
					text,
				},
				MaxTokens:   gpt3.IntPtr(10),
				Temperature: gpt3.Float32Ptr(0),
			}, func(resp *gpt3.CompletionResponse) {
				fmt.Print(resp.Choices[0].Text)
			})
		if err != nil {
			fmt.Println(err)
		}
  }
}

