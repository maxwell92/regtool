// This program is a strengthen version of reglist.sh
// Maintainer: liyao.miao@yeepay.com
// Date: 2016-01-27


package main
import (
    "fmt"
    "net/http"
    "io/ioutil"
    "crypto/tls"
    "encoding/json"
)

type repoInfo struct {
    Repositories []string `json: "repositories"`
}

var defaultProtocol string
var defaultRepoDomain string
var defaultRepoPort string
var repoInfoPath string


func regGet() (body []byte, err error) {

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    client := &http.Client{Transport: tr}
    resp, err := client.Get(defaultProtocol + defaultRepoDomain + defaultRepoPort + repoInfoPath)
    if err != nil {
        //panic(err)
        return nil, err
    }
   
    defer resp.Body.Close()
   
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        //panic(err)
        return nil, err
    } 
    fmt.Println("------------------original response------------------") 
    fmt.Println(string(body))
    return body, nil
}

func resolveToStruct(response []byte) (s *repoInfo, err error) {
    json.Unmarshal(response, &s)     
    return s, nil

}   

func showResult(rs *repoInfo) {

    fmt.Println("------------------structure response------------------") 
    for i := 0 ; i < len(rs.Repositories); i++ {
        fmt.Println(rs.Repositories[i])
    }
    //fmt.Println(*rs) 
}

func main() {
    defaultProtocol = "https://"
    defaultRepoDomain = "registry.test.com"
    defaultRepoPort = ":5000"
    repoInfoPath = "/v2/_catalog"
    
    var reponse []byte  
 
    reponse, err := regGet()
    if err != nil {
        panic(err)
    }   

    
    rs := new(repoInfo)   
    rs, err = resolveToStruct(reponse)      

    showResult(rs)     
     
}
