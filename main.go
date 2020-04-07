package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	csvFileName := flag.String("csv", "problems.csv", "the problems file, it should be csv")
	limit := flag.Int("limit", 10, "time limit for quiz in seconds")

	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file %s\n", *csvFileName))
	}

	r := csv.NewReader(file)

	lines, err := r.ReadAll()

	if err != nil {
		exit("Failed!!")
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*limit) * time.Second)

	correct := 0

problemloop:
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, problem.q)

		answerChan := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChan <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d\n", correct, len(problems))
			break problemloop
		case ans := <-answerChan:
			if ans == problem.a {
				correct++
				fmt.Println("Correct!")
			} else {
				fmt.Println("Wrong!!")
			}
		}
	}

	fmt.Printf("\nYou scored %d out of %d\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	problems := make([]problem, len(lines))

	for i, line := range lines {
		problems[i] = problem{
			q: strings.TrimSpace(line[0]),
			a: strings.TrimSpace(line[1]),
		}
	}

	return problems
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
