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
    "strings"
)

type repoInfo struct {
    Repositories []string `json: "repositories"`
}

type tagList struct {
    Tags []string `json: "tags"`
    Name string `json: "name"`
}


var tagsMap map[string][]string

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
        resp, err := client.Get(defaultProtocol + defaultRepoDomain + defaultRepoPort + "/v2/" + rs.Repositories[i] + "/tags/list")
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

func searchItem(image string, tag string, tagsMap map[string][]string) {
    if tag != "" {
        var i int
        for i = 0; i < len(tagsMap[image]); i++ {
            if strings.Compare(tagsMap[image][i], tag) == 0 {
                fmt.Println("found!") 
                break
            }
        }
        if i == len(tagsMap[image]) {
            fmt.Println("Not found!") 
            fmt.Println(image + ":")
            for i := 0; i < len(tagsMap[image]); i++ {
                fmt.Println(tagsMap[image][i])
            }
            fmt.Printf("Total: %d\n", len(tagsMap[image]))
        }
    } else {
        fmt.Println(image + ":")
        for i := 0; i < len(tagsMap[image]); i++ {
            fmt.Println(tagsMap[image][i])
        }
        fmt.Printf("Total: %d\n", len(tagsMap[image]))
    }
}


func main() {
    if len(os.Args) == 1 {
        fmt.Println("Usage:")
        fmt.Println("regsearch IMAGE[:TAG]")
        os.Exit(-1)
    }
    var image string
    var tag string 
    if strings.Contains(os.Args[1], ":") {
        args := strings.Split(os.Args[1], ":")
        image = args[0]
        tag = args[1]
    } else {
        image = os.Args[1]
        tag = "" 
    } 

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
   
    tagsMap := make(map[string][]string)
    getTags(rs, tagsMap)
    if tagsMap[image] == nil {
        fmt.Println("Not Found")
        os.Exit(-2)
    }
    
    searchItem(image, tag, tagsMap)    

}
