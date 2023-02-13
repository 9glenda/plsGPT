package main

import (
	"context"
	"encoding/json"

	//"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/alecthomas/chroma/quick"
)

type payload struct {
	Text string `json:"text"`
}

type response struct {
	EOF   bool   `json:"eof"`
	Error string `json:"error"`
	Text  string `json:"text"`
}

type Config struct {
	Commands map[string]string
}

const (
	blue  = "\033[34m"
	reset = "\033[0m"

	pink = "\033[95m"
)

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
				MaxTokens:   gpt3.IntPtr(10),
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

func Gpt(text string, apiKey string) string {
	client := gpt3.NewClient(apiKey)
	var j = false
	if j {
		log.Fatal(doJson(client, os.Stdin, os.Stdout))
	}

	Result := ""

	err := client.CompletionStreamWithEngine(
		context.Background(),
		gpt3.TextDavinci003Engine,
		gpt3.CompletionRequest{
			Prompt: []string{
				text,
			},
			MaxTokens:   gpt3.IntPtr(20),
			Temperature: gpt3.Float32Ptr(0),
		}, func(resp *gpt3.CompletionResponse) {
			//fmt.Print(resp.Choices[0].Text)
			if resp.Choices[0].Text != "" && resp.Choices[0].Text != "\n" {
				Result += resp.Choices[0].Text
			}
		})
	if err != nil {
		log.Panic(err)
	}
	return strings.TrimPrefix(strings.TrimLeft(Result, "\n"), "$ ")
}

func fakeGpt(_ string, _ string) string {
	return "fakegpt run \"hello\" <hello>"
}


func fillTemplateOld(template string) string {
  for {
    start := strings.Index(template, "<")
    end := strings.Index(template, ">")
    if start != -1 && end != -1 && end > start {
      placeholder := template[start+1 : end]
      var value string
      fmt.Printf("Please fill in the %s%s%s: ", pink, placeholder, reset)
      fmt.Scanln(&value)
      if value == "" {
        value = "<" + placeholder + ">"
      }
      template = strings.ReplaceAll(template, "<"+placeholder+">", value)
    } else {
      break
    }
  }
  return template
}

func fillTemplate(template string) string {
  for {
    start := strings.Index(template, "<")
    end := strings.Index(template, ">")
    if start != -1 && end != -1 && end > start {
      placeholder := template[start+1 : end]
      var value string
      fmt.Printf("Please fill in the %s%s%s: ", pink, placeholder, reset)
      fmt.Scanln(&value)
      if value == "" {
        value = "<" + placeholder + ">"
      }
      template = strings.ReplaceAll(template, "<"+placeholder+">", value)
    } else {
      break
    }
  }
  return template
}


func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing OPENAI_API_KEY env variable")
	}

	args := os.Args
	if len(args) == 2 {
		text := "plain text linux command, no explonation text, for text to fill in manually use <>, to" + args[1]
		resp := Gpt(text, apiKey)
		resp = fillTemplate(resp)
		fmt.Printf("%s$ %s", blue, reset)
		err := quick.Highlight(os.Stdout, fmt.Sprint(resp), "bash", "terminal16m", "")
		if err != nil {
			fmt.Println(err)
		}
	}
}
