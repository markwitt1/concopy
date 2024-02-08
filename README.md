# concopy

concopy (concat-copy) is a small CLI tool thata allows you to quickly paste your codebase into your preferred Generative AI assistant. Use ChatGPT, Bard or any company-internal LLM to write code for you, while having a lot of context about your codebase.

## Installation

Download the binary for your platform and put it somewhere in your PATH.

## Usage

### By command line arguments

Run `concopy <file1> <file2> <folder1>` to copy the content of the specified files to the clipboard.

### Using .concopyuse

Add a `.concopyuse` file at your project root and populate it with the file / directory paths separated with line breaks.

For example:

```
go.mod
.github
main.go
```
