# cloudNote API

标签： cludeNote

---

## 用户

---

## 分组

### 1.获取所有分组
 - URL：/V1.0/Category/All
 - METHOD：GET
 - **接口地址：** /V1.0/Category/All
 - **支持格式：** Json
 - **请求方式：** GET
 - **接口备注：**
 - **调用样例及调试工具：**
 - **请求head参数说明：**

    >| 参数名     | 类型 | 是否必传   | 说明 |
    >| :------- | ----: | ---: |:--- |
    >|AccountID|    String	|Y	|用户唯一标识|
    >|Token|    String	|Y	|认证令牌|


 - **请求body参数说明：**
    >无
    

 - **请求示例：**
    ```
    curl -X GET --header 'Accept: application/json' --header 'AccountID:s837dh2jud' --header 'Token:kd7w38sjf7sk28dus920dkjs8d0wmf' 'http://10.0.12.104:3000/V1.0/Category/All'
    ```

 - **返回参数说明**
    >| 参数名     | 类型 |  说明 |
    >| :------- | :----: | :--- |
    >|Code | Number | 返回错误代码 0：返回正常 |
    >|Message| String | 错误消息|
    >|Data |Array | 返回对象|
    >|CategoryPath |String | 分组路径|
    >|CategoryName |String | 分组名称|
    >|CreateTime |String | 创建时间|
    >|UpdateTime |String | 更新时间|

 - **成功示例**
    ```json
    {
        "Code": 0,
        "Message": "",
        "Data": [
            {
                "Path": "/28dhw2j0wjd8wh2",
                "Name": "test",
                "CreateTime": "2016-01-02 00:00:00",
                "UpdateTime": "2016-01-02 00:00:00"
            },
            {
                "Path": "/29d8h2jd7giwnhs",
                "Name": "是的方法",
                "CreateTime": "2016-01-02 00:00:00",
                "UpdateTime": "2016-01-02 00:00:00"
            },
            {
                "Path": "/dGVzdG15Tm90ZWJvb2s=",
                "Name": "myNotebook",
                "CreateTime": "2016-01-02 00:00:00",
                "UpdateTime": "2016-01-02 00:00:00"
            }
        ]
    }
    ```

 - **失败示例**
    ```
    {
        "Code": 100,
        "Message": "错误消息",
        "Data": []
    }
    ```

---



### 2.添加分组
### 3.删除分组
### 4.修改分组

---

## 笔记

---

## 附件

---

## 模板

---

## 背景

---