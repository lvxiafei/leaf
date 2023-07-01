package main

import (
	"fmt"
	"log"
	"os/exec"
)

func execWrapper() {
	cmd := exec.Command("ls", "-l", "/var/log/")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}

func UploadFile() (error, string) {

	cmd := exec.Command("curl", "--upload-file", "./output_items.json", "https://transfer.sh/output_items.json")
	out, err := cmd.CombinedOutput()
	return err, string(out)
}
func main() {
	err, outStr := UploadFile()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", outStr)
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	log.Printf("output res: %v", outStr)
}
