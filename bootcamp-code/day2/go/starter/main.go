package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

// Constants required to talk to your local bitcoind node
const (
    rpcURL      = "http://127.0.0.1:18443/"
    rpcUser     = "bjj"
    rpcPassword = "smile123"
)

type rpcRequest struct {
    JSONRPC string `json:"jsonrpc"`
    ID      string `json:"id"`
    Method  string `json:"method"`
    Params  []any  `json:"params"`
}

type rpcResponse struct {
    Result json.RawMessage `json:"result"`
    Error  *struct {
        Message string `json:"message"`
    } `json:"error"`
}

// The RPC Helper provided in slides 5 & 6
func rpc(method string, params []any, wallet string, out any) error {
    url := rpcURL
    if wallet != "" {
        url += "wallet/" + wallet
    }

    body, err := json.Marshal(rpcRequest{
        JSONRPC: "1.0",
        ID:      "explorer",
        Method:  method,
        Params:  params,
    })
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", url, bytes.NewReader(body))
    if err != nil {
        return err
    }
    req.SetBasicAuth(rpcUser, rpcPassword)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    var parsed rpcResponse
    if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
        return err
    }

    if parsed.Error != nil {
        return fmt.Errorf("RPC error: %s", parsed.Error.Message)
    }

    return json.Unmarshal(parsed.Result, out)
}

func main() {
    fmt.Println("Starting Bitcoin Explorer...")

    // Execute challenges sequentially
    if err := showBlockchainInfo(); err != nil {
        fmt.Println("Error:", err)
    }
    
    if err := showWalletBalance("alice"); err != nil {
        fmt.Println("Error:", err)
    }

    if err := listTransactions("alice", 5); err != nil {
        fmt.Println("Error:", err)
    }

    if err := showBlock(""); err != nil {
        fmt.Println("Error:", err)
    }
}
