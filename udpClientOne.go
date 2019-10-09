package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
        "strings"
        "math/rand"
        "time"
)

func main() {
        // arguments := os.Args //command line arguments
        // if len(arguments) == 1 { //check they are given
        //         fmt.Println("Please provide a host:port string")
        //         return
        // }
        // CONNECT := arguments[1] // set the host:port string to const connect
        // stringsArray := []string {
        //     "Hello\n",
        //     "Bonjour\n",
        //     "Hola\n",
        //     "World\n",
        //     "Monde\n",
        //     "Mundo\n",
        // }
        CONNECT := "127.0.0.1:5000"
        s, err := net.ResolveUDPAddr("udp4", CONNECT) //returns a edp end point address of type UDPAddres
        c, err := net.DialUDP("udp4", nil, s) //connection to the udp server is established here
        if err != nil {
                fmt.Println(err)
                return
        } //if connection isnt working

        fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
        //confirmation that the server is connected to the specified
        defer c.Close() // defer the closure of the connection. wait untill you can. When main retuns
        rand.Seed(time.Now().UnixNano())
        overallmessages := 0

        for { // infinite loop reading untill a return which happens on reading
            // a STOP

                reader := bufio.NewReader(os.Stdin) //reads input from command line
                fmt.Print(">> ")
                text, _ := reader.ReadString('\n') //read until there is a new line.
                // n := rand.Int() % len(stringsArray)

                // text := stringsArray[rand.Intn(len(stringsArray))]
                // fmt.Println(text)
                data := []byte(text) //send with a new line
                _, err = c.Write(data)
                if err != nil {
                        fmt.Println(err)
                        return
                }
                if strings.TrimSpace(string(text)) == "STOP" {
                       fmt.Println("Exiting UDP client!")
                       return
                }
                buffer := make([]byte, 1024)
                n, _, err := c.ReadFromUDP(buffer)
                if err != nil {
                        fmt.Println(err)
                        return
                }
                if strings.TrimSpace(string(data)) == "STOP" {
                       fmt.Println("Exiting UDP client!")
                       return
               }
               if strings.TrimSpace(string(data)) == "SWITCH" {
                      fmt.Println("Switch from UDP client!")
                      continue
              }
                fmt.Printf("Reply: %s\n", string(buffer[0:n]))
                overallmessages ++
                fmt.Println("We have made: ", overallmessages)
                if overallmessages > 10000 {
                    fmt.Printf("Reply: %s\n", string(buffer[0:n]))
                    c.Write([]byte("STOP"))
                    fmt.Println("Exiting UDP client!")
                    return
                }
        }
}
