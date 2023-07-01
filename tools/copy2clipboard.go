package tools

import "fmt"
import "os/exec"

func Copy2Clipboard() error {

	cmdStr := fmt.Sprintf("cat output_items.json | yank")
	cmd := exec.Command("sh", "-c", cmdStr)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Error: %v\n", err)
	} else {
		fmt.Printf("Success: %s\n%s\n", cmdStr, out)
		return nil
	}
}
