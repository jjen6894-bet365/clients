package main

import (
    "net"
    "fmt"
    "bufio"
    "strings"
    "time"
    "os"
    "storage"
    "reflect"
    "encoding/json"
)

func main() {
    fmt.Println("Server is starting... ")
    // store := make([]string, 0)
    // passedStored := &store
    var Storages []storage.Storage
    passedStorage := &Storages
    /* You can extract functionality like decision making to go routines
     instead of to a handler which wouldn't have worked in this example as
     the abstraction wouldnt know whether the connection was udp or tcp till
     it was made
    */
    go serveTcp(passedStorage) //go routine for tcp connections
    go serveUDP(passedStorage) //go routine for udp connections
    // to keep main alive to allow connections in.
    for {
        buf := bufio.NewReader(os.Stdin)
        fmt.Print("> ")
        sentence, err := buf.ReadBytes('\n')
        if err != nil {
            fmt.Println("Error : ", err)
        }
        if strings.TrimSpace(string(sentence)) == "STOP" {
            return
        }
    }
}

func serveTcp(store *[]storage.Storage) {
    tcpListener, tcpErr := net.Listen("tcp4", ":5000")
    if tcpErr != nil {
        fmt.Println("Error : ", tcpErr)
        return
    }
    fmt.Println("TCP Server is running... ")
    for {
        tcpConnection, tcpConnectionErr := tcpListener.Accept()
        if tcpConnectionErr != nil {
            fmt.Println("Error : ", tcpConnectionErr)
            return
        }
        fmt.Println("TCP connection made: ")
        defer tcpListener.Close()
        go handleTcp(tcpConnection, store)
    }
}

func handleTcp(tcpConnection net.Conn, store *[]storage.Storage) {
    tcpNetData, err := bufio.NewReader(tcpConnection).ReadString('\n')

    if err != nil {
        fmt.Println("Error : ", err)
        return
    }
    if strings.TrimSpace(string(tcpNetData)) == "STOP" {
        fmt.Println("Exiting TCP server!")
        return
    }
    if strings.TrimSpace(string(tcpNetData)) == "STORE" {
        fmt.Println("Accessing the store!")
        for _, v := range *store {
            fmt.Println(v)
        }
    }
    storeKeyValue := storage.Storage{}
    // emptyKeyValue := storage.Storage{}

    json.Unmarshal([]byte(tcpNetData), &storeKeyValue)
    fmt.Println(storeKeyValue)
    fmt.Println(storeKeyValue.Key)
    fmt.Println(storeKeyValue.Value)
    // if (reflect.DeepEqual(emptyKeyValue, storeKeyValue)) {
    //     fmt.Println("Error: Unproccesable entitiy")
    //     tcpConnection.Write([]byte("Error: Unproccesable entitiy"))
    //
    // }
    // if doesKeyExist(*store, storeKeyValue.Key) {
    //     tcpConnection.Write([]byte("Error: Key already exists"))
    // }
    *store = append(*store, storeKeyValue)
    // *store = append(*store, strings.TrimSpace(string(tcpNetData)))

    fmt.Print("Store ", *store, "\n")
    fmt.Println("")
    t := time.Now()
    myTime := t.Format(time.RFC3339) + "\n"
    tcpConnection.Write([]byte(myTime))

}

func serveUDP(store *[]storage.Storage) {
    udpAddr, udpErr := net.ResolveUDPAddr("udp4", ":5000")
    if udpErr != nil {
        fmt.Println("Error : ", udpErr)
        return
    }
    fmt.Println("UDP Server is running... ")
    udpConnection, udpConnectionErr := net.ListenUDP("udp4", udpAddr)
    if udpConnectionErr != nil {
        fmt.Println("Error : ", udpConnectionErr)
        return
    }
    fmt.Println("UDP connection made: ")
    defer udpConnection.Close()
    buffer := make([]byte, 1024)
    for {
        handleUDP(buffer, udpConnection, store)
    }

}

func handleUDP(buffer []byte, udpConnection *net.UDPConn, store *[]storage.Storage) {
    udpNetData, udpAddr, udpErr := udpConnection.ReadFromUDP(buffer)
    fmt.Println(buffer[0:udpNetData])
    fmt.Println(string(buffer[0:udpNetData]))

    if udpErr != nil {
        fmt.Println("Error : ", udpErr)
        return
    }
    if strings.TrimSpace(string(buffer[0:udpNetData])) == "STOP" {
        fmt.Println("Exiting UDP server!")
        return
    }
    if strings.TrimSpace(string(buffer[0:udpNetData])) == "STORE" {
        fmt.Println("Accessing the store!")
        for _, v := range *store {
            fmt.Println(v)
        }
        storeToSend, _ := json.Marshal(*store)
        fmt.Println("Store to send: ", storeToSend)
        udpConnection.WriteToUDP(storeToSend, udpAddr)
        return
    }
    storeKeyValue := storage.Storage{}
    emptyKeyValue := storage.Storage{}

    json.Unmarshal(buffer[0:udpNetData], &storeKeyValue)
    fmt.Println(storeKeyValue)
    fmt.Println(storeKeyValue.Key)
    fmt.Println(storeKeyValue.Value)
    if (reflect.DeepEqual(emptyKeyValue, storeKeyValue)) {
        fmt.Println("Error: Unproccesable entitiy")
        udpConnection.WriteToUDP([]byte("Error: Unproccesable entitiy"), udpAddr)
        return
    }
    if doesKeyExist(*store, storeKeyValue.Key) {
        udpConnection.WriteToUDP([]byte("Error: Key already exists"), udpAddr)
        return
    }
    *store = append(*store, storeKeyValue)
    fmt.Print("Store ", *store, "\n")
    fmt.Println("")
    t := time.Now()
    myTime := t.Format(time.RFC3339) + "\n"
    _, writeErr := udpConnection.WriteToUDP([]byte(myTime), udpAddr)
    if writeErr != nil {
        fmt.Println("Error : ", writeErr)
        return
    }
}

func doesKeyExist(storage []storage.Storage, key string) bool {
    for _, store := range storage {
        if store.Key == key {
            fmt.Println("ERROR: Key already exists", key)
            // fmt.Fprintf(w, "That Key exists already")
            return true
        }
    }
    return false
}
