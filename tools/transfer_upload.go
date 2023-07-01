package tools

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"leaf/util"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
)

func UploadFile() (error, string) {

	cmd := exec.Command("curl", "--upload-file", util.OutputJson, "https://transfer.sh/output.json")
	out, err := cmd.CombinedOutput()
	return err, string(out)
}

func main() {
	// 文件路径
	filePath := "./output_items.json"

	// 目标 URL
	uploadURL := "https://transfer.sh/output_items.json"

	// 创建一个 buffer 用于存储文件内容
	fileBuffer := &bytes.Buffer{}

	// 创建一个 multipart writer
	multipartWriter := multipart.NewWriter(fileBuffer)

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建一个 form 文件字段
	fileWriter, err := multipartWriter.CreateFormFile("file", filePath)
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// 将文件内容复制到 form 文件字段中
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		fmt.Println("Error copying file content:", err)
		return
	}

	// 完成 multipart 写入
	multipartWriter.Close()

	// 创建 HTTP 请求
	request, err := http.NewRequest("POST", uploadURL, fileBuffer)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	// 设置请求头，指定 Content-Type 为 multipart/form-data
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer response.Body.Close()

	// 读取响应内容
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 输出响应内容
	fmt.Println("Response:", string(responseBody))
}
