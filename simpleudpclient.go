package main
//Connect udp
import (
    "net"
    "fmt"
)

func main() {
    conn, err := net.Dial("udp", "127.0.0.1:5000")
    if err != nil {
        return err
    }
    defer conn.Close()

    //simple Read
    buffer := make([]byte, 1024)
    input := conn.Read(buffer)
    // fmt.Println(input)

    //simple write
    conn.Write([]byte("Hello from client"))
}
