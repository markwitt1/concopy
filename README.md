# concopy

concopy (concatenate-copy) is a compact Command Line Interface (CLI) tool designed to facilitate the seamless pasting of your codebase into your preferred Generative AI assistant. Whether you're using ChatGPT, Bard, or any company-internal Language Model, concopy enables you to write code efficiently by allowing you to select the specific context about your codebase to supply and iterate on the code using the AI assistants output.

![demo](images/demo.gif)

concopy proves particularly useful when you need to generate code for a specific segment of your codebase without sharing the entire codebase with the AI. It can sometimes outperform GitHub Copilot by enabling you to select the context of the code you wish to generate. Copilot's RAG occasionally focuses on irrelevant parts of the codebase, leading to undesired results.

![copilot wrong result](images/copilot_screenshot.png)

As illustrated in the screenshot, Copilot only injects lines 1 to 32 into the prompt, despite the file containing many more lines.

## Installation

### Using Homebrew

You can install `concopy` using Homebrew by tapping into our repository and then installing the formula:

```sh
brew tap markwitt1/tap https://github.com/markwitt1/homebrew-tap
brew install concopy
```

### Other Methods

Download the binary compatible with your platform and place it somewhere in your PATH.

Run `chmod +x path/to/binary` to make the binary executable.

## Usage

### By Command Line Arguments

Run `concopy <file1> <file2> <folder1>` to copy the content of the specified files to the clipboard.

### Using .concopyuse

Add a `.concopyuse` file at your project root and populate it with the file/directory paths separated by line breaks.

For example:

```
go.mod
.github
main.go
```

NOTE: On the first run on macOS, the system may block the execution. You will need to navigate to `System Settings -> Privacy & Security -> Security` and allow the execution of the binary.
