# plsGPT

plsGPT is a Go (golang) project that utilizes GPT-3 (Generative Pretrained Transformer 3) to generate shell commands from natural language descriptions.
## Installation

To install plsGPT, you need to have Go (golang) installed on your system. Once you have Go installed, simply run the following command to download and install the package:


```sh
go install github.com/9glenda/plsGPT
```
### NixOS
```sh
nix shell github:9glenda/plsGPT 
```

## Setup

Before you can use plsGPT, you need to set the following environment variable with your OpenAI API key:

```sh
export OPENAI_API_KEY="key"
```
To get an OpenAI API key, you need to sign up for an OpenAI account and subscribe to the GPT-3 API. Once you have subscribed, you will be able to generate an API key from the OpenAI dashboard.
Usage

Using plsGPT is easy. Simply run the following command in your terminal, filling in the description of the desired shell command in place of "create a folder with a given name":

```sh
plsgpt "create a folder with a given name"
```
### Templating
The program will prompt you for additional information, such as the name of the folder in the example above:

```raw
Please fill in the name: hello
```
Once you have provided the necessary information, plsGPT will generate the corresponding shell command:

```raw
$ mkdir hello
```
## Contributions

We welcome contributions to plsGPT! If you find a bug or have an idea for a new feature, please open an issue or submit a pull request.
