package main

import (
        // "bufio"
        "fmt"
        "net"
        // "os"
        "strings"
        "math/rand"
        "time"
        "encoding/json"
        "storage"
)

func main() {
        // arguments := os.Args //command line arguments
        // if len(arguments) == 1 { //check they are given
        //         fmt.Println("Please provide a host:port string")
        //         return
        // }
        // CONNECT := arguments[1] // set the host:port string to const connect
        stringsArray := []string {
            "POST,Hello,World",
            "POST,Goodbye,cruelWorld",
            "GET,Hello, ",
            "PUT,Hello, Monde",
            "DELETE,goodbye, ",
        }
        CONNECT := "127.0.0.1:5001"
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
                // storeKeyValue := storage.Storage{}
                // fmt.Print(">> ")
                // err := json.NewDecoder(os.Stdin).Decode(&storeKeyValue)
                // if err != nil {
                //
                //     fmt.Println("ERRROR reading standard input: ",err)
                //     return
                // }
                text := stringsArray[rand.Intn(len(stringsArray))]

                // reader := bufio.NewReader(os.Stdin)
                // fmt.Print(">> ")
                // text, err := reader.ReadString('\n')
                // if err != nil {
                //
                //     fmt.Println("ERRROR reading standard input: ",err)
                //     return
                // }
                if strings.TrimSpace(string(text)) == "STOP" {
                    fmt.Printf("Reply: %s\n", string(text))
                    fmt.Println("Exiting UDP client!")
                    return
                }
                // reader := bufio.NewReader(os.Stdin) //reads input from command line
                // text, _ := reader.ReadString('\n') //read until there is a new line.
                // n := rand.Int() % len(stringsArray)
                // fmt.Println(text)
                // fmt.Println(storeKeyValue.Key)
                // fmt.Println(storeKeyValue.Value)
                // json.Unmarshal([]byte(text), &storeKeyValue)
                withoutSpaces := strings.Replace(text, " ", "", -1)
                fmt.Println(withoutSpaces)

                removedNewLine := strings.TrimSuffix(withoutSpaces, "\n")
                fmt.Println(removedNewLine)
                requestString := strings.Split(removedNewLine, ",")
                fmt.Println(len(requestString))
                // fmt.Println("command ", requestString[len(requestString)-len(requestString)])
                // fmt.Println("Key ", requestString[-1])
                // fmt.Println("value ", requestString[len(requestString)])
                if len(requestString) < 3 {
                    fmt.Println("error, not enough arguments passed in")
                    continue
                }
                storeKeyValue := storage.Storage{Command: requestString[len(requestString)-len(requestString)], Key: requestString[len(requestString)-2], Value: requestString[len(requestString)-1]}
                fmt.Println(storeKeyValue)
                fmt.Println(storeKeyValue.Key)
                fmt.Println(storeKeyValue.Value)
                fmt.Println(storeKeyValue.Command)

                // text := stringsArray[rand.Intn(len(stringsArray))]
                // fmt.Println(text)
                dataStruct, err := json.Marshal(&storeKeyValue)
                if err != nil {
                    fmt.Println("ERror: ", err)
                    return
                }
                // data := []byte(storeKeyValue) //send with a new line
                _, err = c.Write(dataStruct)
                if err != nil {
                        fmt.Println(err)
                        return
                }
                // if strings.TrimSpace(string(text)) == "STORE" {
                //         // fmt.Printf("Reply: %s\n", string(text))
                //         // fmt.Println("Exiting UDP client!")
                //         // return
                //         _, err = c.Write([]byte(text))
                //
                // }

                buffer := make([]byte, 1024)
                n, _, err := c.ReadFromUDP(buffer)
                if err != nil {
                        fmt.Println(err)
                        return
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
