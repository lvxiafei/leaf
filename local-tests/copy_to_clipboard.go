package main

import "fmt"
import "os/exec"

func main() {

	//cat your_file.txt | yank
	cmdStr := fmt.Sprintf("cat your_file.txt | yank")
	cmd := exec.Command("sh", "-c", cmdStr)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Errorf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: %s\n%s\n", cmdStr, out)
	}
}
