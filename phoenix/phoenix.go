package phoenix

import (
	"fmt"
	"os/exec"
	"os"
	"time"
)

func Phoenix() {
	for {
		filename := "phoenix/phoenix.txt"		
		file, _ := os.Stat(filename)
		modifiedtime := file.ModTime()
		fmt.Println("Last modified time : ", modifiedtime)

		data, _ := os.ReadFile(filename)
		fmt.Println("You recieved: ", string(data))

		if modifiedtime.Add(2 * time.Second).Before(time.Now()) {
			fmt.Println("Spawning new program ")
			break
		}
		time.Sleep(time.Second)
	}
	
	//
	cmnd := exec.Command("gnome-terminal", "--", "go","run","./main.go")
	cmnd.Run()

	for {
		filename := "phoenix/phoenix.txt"
		msg := fmt.Sprintln("Writing to file", time.Now())
		os.WriteFile(filename, []byte(msg), 0666)

		time.Sleep(time.Second)
	}
}