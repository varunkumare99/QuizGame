# Go Timed Quiz

A simple command-line quiz application written in Go. It reads questions and answers from a CSV file, runs the quiz with a configurable time limit, and displays the final score.

## Features

- Read questions from a CSV file.
- Set a custom time limit.
- Track attempted and correctly answered questions.
- Use goroutines and channels for concurrency.

## Usage

Run the quiz with the default settings:

```bash
go run .
```

Set a custom time limit:

```bash
go run . -limit=30
```

Specify a CSV file:

```bash
go run . -csv=./problems.csv
```

## CSV Format

Create a `problems.csv` file with questions and answers separated by commas:

```csv
5+5,10
7+3,10
1+2,3
8-4,4
```

## Built With

- Go
- Go standard library


