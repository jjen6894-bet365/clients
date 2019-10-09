package main

import (
    "net"
    "fmt"
    "bufio"
    "strings"
    "time"
    "os"
    "sync/atomic"

)

func main() {
    fmt.Println("Server is starting... ")
    storeWithInt := make(map[uint64]string)
    passedStoreWithInt := &storeWithInt
    var counter uint64
    pointToCounter := &counter
    /* You can extract functionality like decision making to go routines
    instead of to a handler which wouldn't have worked in this example as
    the abstraction wouldnt know whether the connection was udp or tcp till
    it was made
    */
    go serveTcp(passedStoreWithInt, pointToCounter) //go routine for tcp connections
    go serveUDP(passedStoreWithInt, pointToCounter)
    //go routine for udp connections
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

func serveTcp(store *map[uint64]string, counter *uint64) {
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
        go handleTcp(tcpConnection, store, counter)

    }

}

func handleTcp(tcpConnection net.Conn, store *map[uint64]string, counter *uint64) {
    for {
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

        // *store = append(*store, strings.TrimSpace(string(tcpNetData)))
        getStore := *store
        getStore[uint64(*counter)] = strings.TrimSpace(string(tcpNetData))
        atomic.AddUint64(counter, 1)

        fmt.Print("Store ", *store, "\n")
        fmt.Println("")
        t := time.Now()
        myTime := t.Format(time.RFC3339) + "\n"
        tcpConnection.Write([]byte(myTime))
    }
}

func serveUDP(store *map[uint64]string, counter *uint64) {
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
        handleUDP(buffer, udpConnection, store, counter)
    }

}

func handleUDP(buffer []byte, udpConnection *net.UDPConn, store *map[uint64]string, counter *uint64) {
    udpNetData, udpAddr, udpErr := udpConnection.ReadFromUDP(buffer)
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
    }
    // *store = append(*store, strings.TrimSpace(string(buffer[0:udpNetData])))
    getStore := *store

    getStore[uint64(*counter)] = strings.TrimSpace(string(buffer[0:udpNetData]))
    atomic.AddUint64(counter, 1)

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
