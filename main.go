package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

type quizScore struct {
	attempted int
	correct   int
}

func main() {
	csvFlag := flag.String("csv", "./problems.csv", "a csv file in the format of 'question, answer'")
	limitFlag := flag.Int("limit", 10, "the time limit for the quiz in seconds")
	flag.Parse()

	f, err := os.Open(*csvFlag)
	defer f.Close()
	if err != nil {
		fmt.Println("err opening file")
		panic(err)
	}

	//Pres any key to start quiz
	readByteToStartQuiz()

	done := make(chan struct{})
	result := make(chan quizScore)
	go startQuiz(f, done, result)

	//waiting on timer vs completion of quix
	select {
	case <-time.Tick(time.Duration(*limitFlag) * time.Second):
		close(done)
		fmt.Println("\n\nYou ran out of time")
	case <-done:
		fmt.Println("\nCongratz you were able to attempt all the questions in time")
	}

	fmt.Println("Quiz is complete")

	//get quiz results
	score := <-result

	rewindFile(f)
	totalq := totalQuestions(f)
	fmt.Printf("\nResult\nTotal Questions: %v\nAttempted: %v\nCorrectly Answered: %v\n", totalq, score.attempted, score.correct)
}

func rewindFile(csvFile *os.File) {
	if _, err := csvFile.Seek(0, io.SeekStart); err != nil {
		fmt.Println("error moving file pointer to beginning")
		panic(err)
	}
}

func totalQuestions(csvFile *os.File) int {
	count := 0
	csvReader := csv.NewReader(csvFile)

	for {
		_, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error in Reading file")
			panic(err)
		}
		count += 1
	}
	return count
}

func startQuiz(csvFile *os.File, done chan struct{}, result chan<- quizScore) {
	csvReader := csv.NewReader(csvFile)
	attempt, score, i := 0, 0, 1

	userInput := make(chan string)
	go takeUserInput(userInput)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error in Reading file")
			panic(err)
		}
		fmt.Printf("Problem #%v: %s = ", i, record[0])
		i += 1
		select {
		case <-done:
			result <- quizScore{attempted: attempt, correct: score}
			return
		case input := <-userInput:
			if input != "" {
				attempt += 1
			}
			if record[1] == input {
				score += 1
			}
		}
	}
	fmt.Println("quiz is done man")
	done <- struct{}{}
	result <- quizScore{attempted: attempt, correct: score}
}

func takeUserInput(userInput chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		userInput <- scanner.Text()
	}
}

func readByteToStartQuiz() {
	fmt.Println("Press any key to start the quiz")
	os.Stdin.Read(make([]byte, 1))
}
