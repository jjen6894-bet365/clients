package main
//Connect udp
import (
    "net"
    "fmt"
)

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:5000")
    if err != nil {
        return
    }
    defer conn.Close()

    //simple Read
    buffer := make([]byte, 1024)
    input, err := conn.Read(buffer)
    if err != nil {
        return
    }
    fmt.Println(input)

    //simple write
    conn.Write([]byte("Hello from client"))
}
