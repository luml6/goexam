# 附件云存储

标签（空格分隔）： 云笔记
---


## 七牛云存储凭证管理

### 1.获取上传凭证
 - URL：/V1.0/Cloud/GetUploadToken
 - METHOD：GET
 - **接口地址：** /V1.0/Cloud/GetUploadToken
 - **支持格式：** Json
 - **请求方式：** POST
 - **接口备注：**
 - **调用样例及调试工具：**
 - **请求head参数说明：**

    >| 参数名     | 类型 | 是否必传   | 说明 |
    >| :------- | ----: | ---: |:--- |
    >|AccountID|    String	|Y	|用户唯一标识|
    >|Token|    String	|Y	|认证令牌|


 - **请求body参数说明：**
    >| 参数名     | 类型 | 是否必传   | 说明 |
    >| :------- | ----: | ---: |:--- |
    >|FileKey|    String	|Y	|云存储文件为唯一key|
    >|Expires|    Number	|Y	|获取Upload token有效时间(单位：秒)|
    >|NoteID|    String	    |Y	|笔记ID|
    

 - **请求示例：**
    ```
    curl -d '{"FileKey":"d01dea53-f72e-43c1-a0dd-16751ff0dbbd", "Expires":3600, "NoteID":"24d85bea-841a-4c4d-afae-31ff70f0539e"}' --header 'Accept: application/json' --header 'AccountID:s837dh2jud' --header 'Token:kd7w38sjf7sk28dus920dkjs8d0wmf' 'http://127.0.0.1:3000/V1.0/Cloud/GetUploadToken'
    ```

 - **返回参数说明**
    >| 参数名     | 类型 |  说明 |
    >| :------- | :----: | :--- |
    >|Code | Number | 返回错误代码 0：返回正常 |
    >|Message| String | 错误消息|
    >|Data |Object | 返回对象|
    >|UploadToken |String | 云存储上传凭证|
    >|Domain|     String	|云存储域名 |
    >|Expires|    Number	|获取Upload token有效时间(单位：秒)|

 - **成功示例**
    ```json
    {
        "Code": 0,
        "Message": "",
        "Data": {
            "Domain": "oe3mg81er.bkt.clouddn.com",
            "Expires": 3600,
            "UploadToken": "Yj9DZSWMsd3RZ9gNql2YpX_eyq9CzZxl-qTDHoFj:saKlwmWZfEPIzovbHx9b4ggTH98=:eyJzY29wZSI6Im5vdGVwYWQvZDAxZGVhNTMtZjcyZS00M2MxLWEwZGQtMTY3NTFmZjBkYmJkIiwiZGVhZGxpbmUiOjE0NzYyNjU1NzUsInVwaG9zdHMiOlsiaHR0cDovL3VwLnFpbml1LmNvbSIsImh0dHA6Ly91cGxvYWQucWluaXUuY29tIiwiLUggdXAucWluaXUuY29tIGh0dHA6Ly8xODMuMTM2LjEzOS4xNiJdfQ=="
        }
    }
    ```

 - **失败示例**
    ```
    {
        "Code": 100,
        "Message": "错误消息"
    }
    ```

---

### 2.获取云回调上传凭证
 - URL：/V1.0/Cloud/GetUploadTokenWithCB
 - METHOD：GET
 - **接口地址：** /V1.0/Cloud/GetUploadToken
 - **支持格式：** Json
 - **请求方式：** POST
 - **接口备注：**
 - **调用样例及调试工具：**
 - **请求head参数说明：**

    >| 参数名     | 类型 | 是否必传   | 说明 |
    >| :------- | ----: | ---: |:--- |
    >|AccountID|    String	|Y	|用户唯一标识|
    >|Token|    String	|Y	|认证令牌|


 - **请求body参数说明：**
    >| 参数名     | 类型 | 是否必传   | 说明 |
    >| :------- | ----: | ---: |:--- |
    >|FileKey|    String	|Y	|云存储文件为唯一key|
    >|Expires|    Number	|Y	|获取Upload token有效时间(单位：秒)|
    >|NoteID|    String	    |Y	|笔记ID|
    >|FileType|    String	|Y	|文件类型|
    >|FileName|    String	|Y	|文件名称|
    

 - **请求示例：**
    ```
    curl -d '{"FileKey":"d01dea53-f72e-43c1-a0dd-16751ff0dbbd", "Expires":3600, "NoteID":"d946f08d-8d84-4d16-9773-0e7d720900d3", "FileName":"myfile", "FileType":"png"}' --header 'Accept: application/json' --header 'AccountID:s837dh2jud' --header 'Token:kd7w38sjf7sk28dus920dkjs8d0wmf' 'http://127.0.0.1:3000/V1.0/Cloud/GetUploadToken'
    ```

 - **返回参数说明**
    >| 参数名     | 类型 |  说明 |
    >| :------- | :----: | :--- |
    >|Code | Number | 返回错误代码 0：返回正常 |
    >|Message| String | 错误消息|
    >|Data |Object | 返回对象|
    >|UploadToken |String | 云存储上传凭证|
    >|Domain|     String	|云存储域名 |
    >|Expires|    Number	|获取Upload token有效时间(单位：秒)|

 - **回调返回参数说明**
    >| 参数名     | 类型 |  说明 |
    >| :------- | :----: | :--- |
    >|Success  | Bool | 成功返回ture，失败返回false |
    >|Key|     String | 传入的FileKey.即附件UUID|
    >|Message |String | 返回对象|

 - **成功示例**
    ```json
    {
        "Code": 0,
        "Message": "",
        "Data": {
            "Domain": "oe3mg81er.bkt.clouddn.com",
            "Expires": 3600,
            "UploadToken": "Yj9DZSWMsd3RZ9gNql2YpX_eyq9CzZxl-qTDHoFj:QqSTSoZuDqbX6Om4LTaC47LcaOY=:eyJzY29wZSI6Im5vdGVwYWQvZDAxZGVhNTMtZjcyZS00M2MxLWEwZGQtMTY3NTFmZjBkYmJkIiwiZGVhZGxpbmUiOjE0NzYyNjY4NzAsImNhbGxiYWNrVXJsIjoiaHR0cDovL25vdGVwYWQuemV1c2lzLmNvbTo1MDA0Ny9WMS4wL0Nsb3VkL0NhbGxCYWNrLyIsImNhbGxiYWNrQm9keSI6IntcIlV1aWRcIjokKGtleSksIFwiRmlsZVR5cGVcIjpcInBuZ1wiLCBcIkZpbGVOYW1lXCI6XCJteWZpbGVcIiwgXCJOb3RlSWRcIjpcImQ5NDZmMDhkLThkODQtNGQxNi05NzczLTBlN2Q3MjA5MDBkM1wifSIsImNhbGxiYWNrQm9keVR5cGUiOiJhcHBsaWNhdGlvbi9qc29uIiwidXBob3N0cyI6WyJodHRwOi8vdXAucWluaXUuY29tIiwiaHR0cDovL3VwbG9hZC5xaW5pdS5jb20iLCItSCB1cC5xaW5pdS5jb20gaHR0cDovLzE4My4xMzYuMTM5LjE2Il19"
        }
    }
    ```

 - **失败示例**
    ```
    {
        "Code": 100,
        "Message": "错误消息"
    }
    ```

---

### 3.获取文件下载URL
 - URL：/V1.0/Cloud/GetDownloadURL
 - METHOD：GET
 - **接口地址：** /V1.0/Cloud/GetDownloadURL
 - **支持格式：** Json
 - **请求方式：** POST
 - **接口备注：**
 - **调用样例及调试工具：**
 - **请求head参数说明：**

    >| 参数名     | 类型 | 是否必传   | 说明 |
    >| :------- | ----: | ---: |:--- |
    >|AccountID|    String	|Y	|用户唯一标识|
    >|Token|    String	|Y	|认证令牌|


 - **请求body参数说明：**
    >|FileKey|    String	|Y	|云存储文件为唯一key|
    >|Expires|    Number	|Y	|获取Download URL有效时间(单位：秒)|
    

 - **请求示例：**
    ```
    curl -d '{"FileKey":"d01dea53-f72e-43c1-a0dd-16751ff0dbbd", "Expires":3600}' --header 'Accept: application/json' --header 'AccountID:s837dh2jud' --header 'Token:kd7w38sjf7sk28dus920dkjs8d0wmf' 'http://127.0.0.1:3000//V1.0/Cloud/GetDownloadURL'
    ```

 - **返回参数说明**
    >| 参数名     | 类型 |  说明 |
    >| :------- | :----: | :--- |
    >|Code | Number | 返回错误代码 0：返回正常 |
    >|Message| String | 错误消息|
    >|Data |Object | 返回对象|
    >|DownloadURL |String | 文件下载URL|
    >|Expires|    Number	|获取Download URL有效时间(单位：秒)|

 - **成功示例**
    ```json
    {
        "Code": 0,
        "Message": "",
        "Data": {
            "DownloadURL":"http://oe3mg81er.bkt.clouddn.com/d01dea53-f72e-43c1-a0dd-16751ff0dbbd?e=1474885461&token=Yj9DZSWMsd3RZ9gNql2YpX_eyq9CzZxl:aKDVpzBO67vsjKbHcEBIp3J5HxE",
            "Expires":3600
        }
    }
    ```

 - **失败示例**
    ```
    {
        "Code": 100,
        "Message": "错误消息"
    }
    ```

---





