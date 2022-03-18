package phoenix

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Init() {
	fmt.Println("Phoenix is waiting")
	filename := "phoenix/phoenix.txt"
	for {
		time.Sleep(time.Second)
		file, _ := os.Stat(filename)
		modifiedtime := file.ModTime()
		// fmt.Println("Last modified time : ", modifiedtime)

		// data, _ := os.ReadFile(filename)
		// fmt.Println("You recieved: ", string(data))

		if modifiedtime.Add(2 * time.Second).Before(time.Now()) {
			fmt.Println("Spawning new program ")
			break
		}
	}

	cmnd := exec.Command("gnome-terminal", "--", "go", "run", "./main.go")
	// if len(os.Args) > 1 {
	// 	cmnd.Args = append(cmnd.Args, os.Args[1])
	// }
	cmnd.Run()
}

func Phoenix() {
	for {
		filename := "phoenix/phoenix.txt"
		msg := fmt.Sprintln("Writing to file", time.Now())
		os.WriteFile(filename, []byte(msg), 0666)

		time.Sleep(time.Second)
	}
}
