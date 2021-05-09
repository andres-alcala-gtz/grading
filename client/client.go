package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"

	"../env"
	"../obj"
)

func main() {
	connection, err := rpc.Dial(env.ConnectionType, env.ConnectionPort)
	if err != nil {
		fmt.Println("CONNECTION ERROR:", err)
		return
	}
	defer connection.Close()

	scanner := bufio.NewScanner(os.Stdin)

	var result string
	var student string
	var subject string
	var grade float64

	var option string

	for option != "0" {
		fmt.Println("\n\nGRADING SYSTEM")
		fmt.Println("[1] Add a grade     (student)")
		fmt.Println("[2] Get the average (student)")
		fmt.Println("[3] Get the average (subject)")
		fmt.Println("[4] Get the average (general)")
		fmt.Println("[0] Exit")

		fmt.Print("Option: ")
		fmt.Scanln(&option)

		fmt.Println()

		switch option {
		case "1":
			fmt.Print("Student: ")
			scanner.Scan()
			student = scanner.Text()

			fmt.Print("Subject: ")
			scanner.Scan()
			subject = scanner.Text()

			fmt.Print("Grade  : ")
			fmt.Scanln(&grade)

			temp := obj.ContainerFull{Student: student, Subject: subject, Grade: grade}

			err := connection.Call("Server.SetGradeStudent", temp, &result)
			if err != nil {
				fmt.Println("ERROR:", err)
			} else {
				fmt.Println(result)
			}

		case "2":
			fmt.Print("Student: ")
			scanner.Scan()
			student = scanner.Text()

			temp := obj.ContainerStudent{Student: student}

			err := connection.Call("Server.GetAverageStudent", temp, &result)
			if err != nil {
				fmt.Println("ERROR:", err)
			} else {
				fmt.Println("Average:", result)
			}

		case "3":
			fmt.Print("Subject: ")
			scanner.Scan()
			subject = scanner.Text()

			temp := obj.ContainerSubject{Subject: subject}

			err := connection.Call("Server.GetAverageSubject", temp, &result)
			if err != nil {
				fmt.Println("ERROR:", err)
			} else {
				fmt.Println("Average:", result)
			}

		case "4":
			temp := obj.ContainerEmpty{}

			err := connection.Call("Server.GetAverageGeneral", temp, &result)
			if err != nil {
				fmt.Println("ERROR:", err)
			} else {
				fmt.Println("Average:", result)
			}

		case "0":
			fmt.Println("Goodbye")

		default:
			fmt.Println("Invalid option")
		}
	}
}
