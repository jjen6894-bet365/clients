package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
    "encoding/json"
    "storage"
)

func main() {

    arguments := os.Args
    if len(arguments) == 1 {
        fmt.Println("Please provide host:port.")
        return
    }

    CONNECT := arguments[1]
    c, err := net.Dial("tcp4", CONNECT)
    if err != nil {
        fmt.Println(err)
        return
    }

    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print(">> ")
        text, err := reader.ReadString('\n')
        if err != nil {
            return
        }
        text = strings.TrimSuffix(text, "\n")

        requestMessage := strings.TrimSpace(text)
        if string(requestMessage) == "STOP" {
            fmt.Println("TCP client exiting...")
            return
        }
        withoutSpaces := strings.Replace(text, " ", "", -1)
        requestString := strings.Split(withoutSpaces, ",")
        if len(requestString) < 3 {
            fmt.Println("error, not enough arguments passed in")
            continue
        }
        fmt.Println(requestString)
        upperCommand := strings.ToUpper(requestString[0])

        storeKeyValue := storage.Storage{Command: upperCommand, Key: requestString[1], Value: requestString[2]}
        // storeKeyValue := storage.Storage{}
        fmt.Println(storeKeyValue)
        // storeKeyValue.Command = s[0]
        // storeKeyValue.Key = s[1]
        // storeKeyValue.Value = s[2]
        // fmt.Print(">> ")
        // err := json.NewDecoder(os.Stdin).Decode(&storeKeyValue)

        if err != nil {
            fmt.Println("ERRROR reading standard input: ",err)
            return
        }

        marshalKVStore, _ := json.Marshal(&storeKeyValue)

        fmt.Fprintf(c, string(marshalKVStore) +"\n")


        message, err := bufio.NewReader(c).ReadString('\n')
        if err != nil {
            fmt.Println("ERROR", err)
            return
        }
        fmt.Print("->: " + message)
    }
}
