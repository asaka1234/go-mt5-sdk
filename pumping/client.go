package pumping

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

// TCP 客户端
func main() {
	//local
	conn, err := net.Dial("tcp", "127.0.0.1:8355")
	if err != nil {
		fmt.Println("err : ", err)
		return
	}
	fmt.Printf("---start-----\n")
	defer conn.Close() // 关闭TCP连接
	//inputReader := bufio.NewReader(os.Stdin)

	go readResponses(conn)

	sendSubTickRequest(conn)

	fmt.Scanln()
}

func readResponses(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Read response error: %v\n", err)
			return
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(line, &resp); err != nil {
			fmt.Printf("Invalid JSON response: %v\n", err)
			continue
		}

		fmt.Printf("Server response: %+v\n", resp)
	}
}

func sendSubTickRequest(conn net.Conn) {
	request := map[string]interface{}{
		"env":  "CPTA",
		"type": "tick",
		"payload": map[string]string{
			"symbols": "", //空就是全部
		},
	}
	sendJSONRequest(conn, request)
}

func sendJSONRequest(conn net.Conn, request map[string]interface{}) {
	data, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		return
	}

	data = append(data, '\n')
	if _, err := conn.Write(data); err != nil {
		fmt.Printf("Error sending request: %v\n", err)
	}
}
