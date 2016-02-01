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
    "os"
)

type repoInfo struct {
    Repositories []string `json: "repositories"`
}

type tagList struct {
    Tags []string `json: "tags"`
    Name string `json: "name"`
}


var tagsMap map[string][]string

var repoProtocol string
var repoDomain string
var repoPort string
var repoInfoPath string

func setDefaultHost() {
    repoProtocol = "https://"
    repoDomain = "registry.test.com"
    repoPort = ":5000"
    repoInfoPath = "/v2/_catalog"
}

func setHost(args string) {
    repoDomain = args
}


func regGet() (body []byte, err error) {

    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    client := &http.Client{Transport: tr}
    resp, err := client.Get(repoProtocol + repoDomain + repoPort + repoInfoPath)
    if err != nil {
        return nil, err
    }
   
    defer resp.Body.Close()
   
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    } 
    return body, nil
}

func resolveToStruct(response []byte) (s *repoInfo, err error) {
    json.Unmarshal(response, &s)     
    return s, nil

}   

func getTags(rs *repoInfo, tagsMap map[string][]string) {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    } 
    client := &http.Client{Transport: tr}
   

    for i := 0; i < len(rs.Repositories); i++ {
        resp, err := client.Get(repoProtocol + repoDomain + repoPort + "/v2/" + rs.Repositories[i] + "/tags/list")
        if err != nil {
            fmt.Println("client.Get")
            panic(err)
        }
   
        defer resp.Body.Close()
        var body []byte 
        body, err = ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Println("ioutil.ReadAll")
            panic(err)
        }

        temp := new(tagList) 
        json.Unmarshal(body, &temp)                
        tagsMap[temp.Name] = temp.Tags
    } 
} 

func showTags(tagsMap map[string][]string, rs *repoInfo) {
    fmt.Printf("Total: %d\n\n", len(rs.Repositories))
    for i := 0; i < len(tagsMap); i++ { 
        fmt.Println(rs.Repositories[i])
        for j := 0; j < len(tagsMap[rs.Repositories[i]]); j++ {
            fmt.Printf("\t\t: %s\n",tagsMap[rs.Repositories[i]][j])
        } 
    }
}

func main() {

    setDefaultHost() 
     
    if len(os.Args) == 3 {
       if os.Args[1] == "-h" {
            setHost(os.Args[2])
       } 
        
    }
    
    var reponse []byte  
 
    reponse, err := regGet()
    if err != nil {
        panic(err)
    }   

    
    rs := new(repoInfo)   
    rs, err = resolveToStruct(reponse)      
   
    tagsMap := make(map[string][]string)
   // tagsMap = getTags(rs)
    getTags(rs, tagsMap)
    showTags(tagsMap, rs)      
}
