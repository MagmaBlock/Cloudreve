package main

import (
    "fmt"
    "encoding/json"
    "crypto/hmac"
    "crypto/sha1"
    "encoding/hex"
    "log"
    "net/http"
    "net/url"
    "strings"
    "io/ioutil"
//    "reflect"
    
    model "github.com/cloudreve/Cloudreve/v3/models"
	"github.com/cloudreve/Cloudreve/v3/pkg/filesystem/fsctx"
	"github.com/cloudreve/Cloudreve/v3/pkg/filesystem/response"
	"github.com/cloudreve/Cloudreve/v3/pkg/request"
	"github.com/cloudreve/Cloudreve/v3/pkg/serializer"
)

// 调用 DogeCloud 的 API
// apiPath：是调用的 API 接口地址，包含 URL 请求参数 QueryString，例如：/console/vfetch/add.json?url=xxx&a=1&b=2
// data：POST 的数据，对象，例如 {a: 1, b: 2}，传递此参数表示不是 GET 请求而是 POST 请求
// jsonMode：数据 data 是否以 JSON 格式请求，默认为 false 则使用表单形式（a=1&b=2）
// 返回值 ret 是一个 map[string]，其中 ret["code"] 为 200 表示 api 请求成功
func DogeCloudAPI(apiPath string, data map[string]interface{}, jsonMode bool) (ret map[string]interface{}) {

    // 这里替换为你的 DogeCloud 永久 AccessKey 和 SecretKey，可在用户中心 - 密钥管理中查看
    AccessKey := handler.Policy.AccessKey
    SecretKey := handler.Policy.SecretKey

    body := ""
    mime := ""
    if jsonMode {
        _body, err := json.Marshal(data)
        if err != nil{ log.Fatalln(err) }
        body = string(_body)
        mime = "application/json"
    } else {
        values := url.Values{}
        for k, v := range data {
            values.Set(k, v.(string))
        }
        body = values.Encode()
        mime = "application/x-www-form-urlencoded"
    }

    signStr := apiPath + "\n" + body
    hmacObj := hmac.New(sha1.New, []byte(SecretKey))
    hmacObj.Write([]byte(signStr))
    sign := hex.EncodeToString(hmacObj.Sum(nil))
    Authorization := "TOKEN " + AccessKey + ":" + sign

    req, err := http.NewRequest("POST", "https://api.dogecloud.com" + apiPath, strings.NewReader(body))
    req.Header.Add("Content-Type", mime)
    req.Header.Add("Authorization", Authorization)
    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil{ log.Fatalln(err) } // 网络错误
    defer resp.Body.Close()
    r, err := ioutil.ReadAll(resp.Body)

    json.Unmarshal([]byte(r), &ret)

    // Debug，正式使用时可以注释掉
    fmt.Printf("[DogeCloudAPI] code: %d, msg: %s, data: %s\n", int(ret["code"].(float64)), ret["msg"], ret["data"])
    return
}
