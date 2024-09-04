package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
);

type problem struct{
	question string
	answer string
}

func problemFetcher(filename string)([] problem,error){
	// 1. open the file 
	file, err := os.Open(filename)
	if err != nil { 
		return nil, fmt.Errorf("error while opening %s file; %s", filename, err.Error()) 
		} 
	defer file.Close()
	// 2. read the file 
	reader := csv.NewReader(file) 
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error while reading data from csv file"+ "format from %s file; %s", filename, err.Error())
	}
	// 3. parse the file and return the data 
	return parseProblem(lines), nil
}

func parseProblem(lines [][]string) ([] problem){
	problems := make([]problem, len(lines))
	for i, eachrecord := range lines { 
		problems[i] = problem{question: eachrecord[0], answer: eachrecord[1]}
    }
	return problems
}

func exit(msg string){
	fmt.Println(msg)
	os.Exit(1)
}

func main(){
	// 1. Input the csv file name 
	filename := "quiz.csv"
	// 2. set the duration of the timer
	timeduration := time.Duration(40) * time.Second
	// 3. pull the data from the csv file
	problems, err := problemFetcher(filename)
	// 4. handle the error 
	if err != nil {
		exit(fmt.Sprintf("error while fetching the data: %s", err.Error()))
	}
	// 5. create variable to count correct answers
	correctAnswers := 0
	// 6. using the duration of the timer, we want to initialize the timer
	timer := time.NewTimer(timeduration)
	// 7. create a for loop to iterate through the questions and print questions and accept the answer
	problemLoop:
		for i, problem := range problems {
			fmt.Printf("Problem %d: %s = ", i+1, problem.question)
			answerCh := make(chan string)
			go func() {
				var answer string
				fmt.Scanf("%s", &answer)
				answerCh <- answer
			}()
			select {
			case <-timer.C:
				// 8. print the final score
				fmt.Printf("\nYou scored %d out of %d.\n", correctAnswers, len(problems))
				break problemLoop
			case answer := <-answerCh:
				if answer == problem.answer {
					correctAnswers++
				}
			}
	}
}
