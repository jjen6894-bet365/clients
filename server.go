package main

import (
    "net"
    "fmt"
    "bufio"
    "strings"
    "time"
    "os"
    "log"
    "net/http"
    "encoding/json"
    "storage"
    "io/ioutil"
)
var store []string
var Storages []storage.Storage


func main() {
    fmt.Println("Server is starting... ")
    // passedStored := &store


    // storeWithInt := make(map[int]string)
    // passedStoreWithInt := &storeWithInt
    /* You can extract functionality like decision making to go routines
    instead of to a handler which wouldn't have worked in this example as
    the abstraction wouldnt know whether the connection was udp or tcp till
    it was made
    */
    go handleRequests() //http
    go serveTcp() //go routine for tcp connections
    go serveUDP()
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
//TCP
func serveTcp() {
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
        go handleTcp(tcpConnection)

    }

}

func handleTcp(tcpConnection net.Conn) {
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
            for _, v := range store {
                fmt.Println(v)
            }
        }
        outcome :=  appendToStore(strings.TrimSpace(string(tcpNetData)))
        if outcome {
            fmt.Print("Store ", store, "\n")
            fmt.Println("")
            t := time.Now()
            myTime := "Committed: " + t.Format(time.RFC3339) + "\n"
            fmt.Println([]byte(myTime))

            tcpConnection.Write([]byte(myTime))
        } else {
            response := "Key already exists: \n" + strings.TrimSpace(string(tcpNetData))
            byteResponse := []byte(response)
            // t := time.Now()
            // myTime := t.Format(time.RFC3339) + "\n"
            tcpConnection.Write(byteResponse)
            // tcpConnection.Write(badRequestMessage)
        }
    }
}
// UDP
func serveUDP() {
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
        handleUDP(buffer, udpConnection)
    }

}

func handleUDP(buffer []byte, udpConnection *net.UDPConn) {
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
        for _, v := range store {
            fmt.Println(v)
        }
    }
    outcome :=  appendToStore(strings.TrimSpace(string(buffer[0:udpNetData])))
    if outcome {
    // store = append(store, strings.TrimSpace(string(tcpNetData)))

        fmt.Print("Store ", store, "\n")
        fmt.Println("")
        t := time.Now()
        myTime := t.Format(time.RFC3339) + "\n"
        _, writeErr := udpConnection.WriteToUDP([]byte(myTime), udpAddr)
        if writeErr != nil {
            fmt.Println("Error : ", writeErr)
            return
        }
    } else {
        response := "BAD REQUEST 404"
        fmt.Println(response)
        byteResponse := []byte(response)
        // t := time.Now()
        // myTime := t.Format(time.RFC3339) + "\n"
        _, writeErr := udpConnection.WriteToUDP(byteResponse, udpAddr) //works but not for tcp?
        if writeErr != nil {
            fmt.Println("Error : ", writeErr)
            return
        }
    }
    // // *store = append(*store, strings.TrimSpace(string(buffer[0:udpNetData])))
    // fmt.Print("Store ", *store, "\n")
    // fmt.Println("")
    // t := time.Now()
    // myTime := t.Format(time.RFC3339) + "\n"
    // _, writeErr := udpConnection.WriteToUDP([]byte(myTime), udpAddr)
    // if writeErr != nil {
    //     fmt.Println("Error : ", writeErr)
    //     return
    // }
}

func appendToStore(data string) bool {
    if !(doesKeyExist(data)) {
        fmt.Println("Key doesnt exist")
        store = append(store, strings.TrimSpace(data))
        // passedMap[strings.TrimSpace(data)] = counter
        return true
    }
    fmt.Println("Key does exist")

    return false
}

func doesKeyExist(newKey string) bool {
    for _, key := range store {
        if key == newKey {
            fmt.Println("ERROR: Key already exists", newKey)
            // fmt.Fprintf(w, "That Key exists already")
            return true
        }
    }
    return false
}

func homePage(w http.ResponseWriter, r *http.Request){

    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func returnAllStorage(w http.ResponseWriter, r *http.Request) {
    // storages := make([]Storage, 2)

    fmt.Println("Endpoint Hit: stored")
    json.NewEncoder(w).Encode(Storages)

}

func returnOneKeyValueStore(w http.ResponseWriter, r *http.Request) {
    // variable := mux.Vars(r)
    for k, v := range r.URL.Query() {
        fmt.Printf("%s: %s\n", k, v)
    }
    fmt.Println("Endpoint Hit: get one key")

    reqBody, _ := ioutil.ReadAll(r.Body)
    var storeKeyValue storage.Storage
    json.Unmarshal(reqBody, &storeKeyValue)
    fmt.Println(storeKeyValue)
    key := storeKeyValue.Key
    // fmt.Fprintf(w, "Key: " + key)
    for _, storage := range Storages {
        if storage.Key == key {
            fmt.Println("Found KEY: ", key)

            json.NewEncoder(w).Encode(storage)
        }
    }
    fmt.Fprintf(w, "Key not found", key)
    fmt.Println("Key does not exist: ", key)
}

func createKeyValueStore(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: create")

    reqBody, _ := ioutil.ReadAll(r.Body)
    fmt.Fprintf(w, "%+v", string(reqBody))
}

func createNewKeyValueStore(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: create")

    // get the body of our POST request
    // unmarshal this into a new Article struct
    // append this to our Articles array.
    reqBody, _ := ioutil.ReadAll(r.Body)
    var storeKeyValue storage.Storage
    json.Unmarshal(reqBody, &storeKeyValue)
    fmt.Println(storeKeyValue)
    // update our global Articles array to include
    // our new Article
    exist := doesKeyExistHttp(Storages, storeKeyValue.Key)
    if exist {
        fmt.Fprintf(w, "That key already exists: %v", storeKeyValue.Key)
    } else {
        Storages = append(Storages, storeKeyValue)
    }

    json.NewEncoder(w).Encode(storeKeyValue)
}

func doesKeyExistHttp(storage []storage.Storage, key string) bool {
    for _, store := range storage {
        if store.Key == key {
            fmt.Println("ERROR: Key already exists", key)
            // fmt.Fprintf(w, "That Key exists already")
            return true
        }
    }
    return false
}

func handleRequestsToStorage(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("This is the request :")
    fmt.Println(r)
    fmt.Printf("This is the method of the request : ")
    fmt.Println(r.Method)

    reqBody, _ := ioutil.ReadAll(r.Body)
    var storeKeyValue storage.Storage
    json.Unmarshal(reqBody, &storeKeyValue)
    fmt.Printf("This is the decoded storage struct of the request : ")
    fmt.Println(storeKeyValue)
    fmt.Printf("This is the request body of the request : ")
    fmt.Println(string(reqBody))

    switch r.Method {
    case "GET":
        // returnOneKeyValueStore(w, r)
        fmt.Println("Endpoint Hit: get one key")
        key := storeKeyValue.Key
        // fmt.Fprintf(w, "Key: " + key)
        for _, storage := range Storages {
            if storage.Key == key {
                fmt.Println("Found KEY: ", key)

                json.NewEncoder(w).Encode(storage)
                break
            }
            fmt.Println("nothing so far....")
        }
        // fmt.Fprintf(w, "Key not found", key)
        // fmt.Println("Key does not exist: ", key)
    case "POST":
        exist := doesKeyExistHttp(Storages, storeKeyValue.Key)
        if exist {
            fmt.Fprintf(w, "That key already exists: ", storeKeyValue.Key)
        } else {
            Storages = append(Storages, storeKeyValue)
        }

        json.NewEncoder(w).Encode(storeKeyValue)
    case "DELETE":
        key := storeKeyValue.Key
        for index, storage := range Storages {
            if storage.Key == key {
                fmt.Println("Found KEY: ", key)
                fmt.Println("Deleting KEY: ", key)
                Storages = append(Storages[:index], Storages[index+1:]...)
                json.NewEncoder(w).Encode(Storages)

            }
            fmt.Println("nothing so far....")
        }
    case "PUT":
        key := storeKeyValue.Key
        for index, storage := range Storages {
            if storage.Key == key {
                fmt.Println("Found KEY: ", key)
                fmt.Println("Updating KEY: ", key)
                Storages[index] = storeKeyValue
                json.NewEncoder(w).Encode(Storages[index])
            }
            fmt.Println("nothing so far....")
        }
    }
}

func handleRequests() {
    //create a mew instance of a mux router
    // myRouter := mux.NewRouter().StrictSlash(true)
    http.HandleFunc("/", homePage)
    http.HandleFunc("/storages", returnAllStorage)
    http.HandleFunc("/storage", handleRequestsToStorage)

    // http.HandleFunc("/storage/", returnOneKeyValueStore)

    log.Fatal(http.ListenAndServe(":10000", nil))

}
