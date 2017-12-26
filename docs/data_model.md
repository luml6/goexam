## mongodb collection设计
- dbname: cloudnote
    
    ### 测试服数据库信息:

    -
    
	>|数据库名称|IP地址|端口|用户名|密码|
	>|:----| ----: | :---- | :---- | :---- |
    >|cloudnote | 192.168.60.75 | 27020 |cloudnote_app |cloudnote_app



    ---
    ### user
    - collection: user
    - 用户信息
    
    >|字段|类型|说明|
    >|:----| :---- | :---- |
    >|UID   |string |   用户唯一标识 |
    >|SyncTime   |string |   上次同步时间 |
    >|UpdateTime   |string |   上次更新时间 |
    >|Usn   |int |   同步系数 |

    ---
    ### category
    - collection: category
    - 分组信息
    
    >|字段|类型|说明|
    >|:----| :---- | :---- |
    >|Path   |string |   笔记分组唯一标识 |
    >|Name   |string |   笔记分组名称 |
    >|UID   |string |   用户唯一标识 |
    >|Status   |string |   分组状态: 0-状态常常；1-已删除 |
    >|CreateTime   |string |   创建时间 |
    >|UpdateTime   |string |   上次更新时间 |
	>|Source   |int |   笔记本来源 1-历史短信|
    >|Usn   |int |   同步系数 |

    ---
    ### note
    - collection: note
    - 笔记信息

    >|字段|类型|说明|
    >|:----| :---- | :---- |
    >|Path   |string |   笔记唯一标识 |
    >|CategoryPath   |string |   笔记分组唯一标识 |
    >|CategoryName   |string |   笔记分组名称 |
    >|UID   |string |   用户唯一标识 |
    >|Title   |string |   笔记标题 |
    >|Author   |string |   笔记作者 |
    >|Summary   |string |   笔记概述 |
    >|Content   |string |   笔记内容 |
	>|CheckSum   |string |   笔记checkSum |
    >|Status   |string |   笔记状态: 0-状态常常；1-已删除,2-预创建,3-回收站删除 |
    >|IsPinned   |uint8 |   笔记是否置顶: 0-不置顶；1-置顶 |
	>|IsShare   |uint8 |   笔记是否被分享: 0-不分享；1-分享 |
	>|IsVote   |uint8 |   笔记是否投票: 0-没有投票；1-投票 |
	>|IsPassword   |uint8 |   笔记是否加密: 0-没有加密；1-已加密 |
	>|ConflictCount   |int |   笔记冲突次数 |
    >|CreateTime   |string |   创建时间 |
    >|UpdateTime   |string |   上次更新时间 |
    >|Usn   |int |   同步系数 |
    
    ---
    ### attachment
    - collection: attachment
    - 附件信息

    >|字段|类型|说明|
    >|:----| :---- | :---- |
    >|ID   |string |   附件唯一标识 |
    >|NotePath   |string |   笔记唯一标识 |
    >|OT   |int |   附件存储类型： 0-第三方云；1-服务器 |
    >|OK   |string |   附件key,若OT=0，则OK为第三方云的key；若OT=1；则OK为附件URL |
    >|Name   |string |   附件文件名称 |
    >|Type   |string |   附件文件类型 |
   
    ---
    ### share
    - collection: attachment
    - 分享表

    >|字段|类型|说明|
    >|:----| :---- | :---- |
    >|ShareID   |string |   分享表唯一标识 |
    >|NoteID   |string |   笔记唯一标识 |
	>|Question   |string |   分享问题 |
	>|Answer   |string |   分享问题答案 |
	>|Password   |string |   分享答案 |
	>|Type   |uint8 |   分享类型 0表示普通分享，1表示加密分享 |
	>|Status   |uint8 |   分享状态 0表示正常，1表示失效 |
    >|FailureTime   |int |   失效时间，暂定单位为分钟（预留字段暂时没用） |
    >|CreateTime   |string |   创建时间 |
    >|OpenNum   |int |   分享链接打开次数 | 
    
	---
	### vote
    - collection: vote
    - 投票表

    >|字段|类型|说明|
    >|:----| :---- | :---- |
    >|VoteID   |string |   投票表唯一标识 |
    >|NoteID   |string |   笔记唯一标识 |
	>|Question   |string |   投票问题|
	>|Type   |uint8 |   投票类型 0表示单选投票，1表示多选投票 |
    >|InvalidTime   |int |   失效时间，暂定单位为秒 |
    >|CreateTime   |string |   创建时间 |
	>|ClickNum   |int |   投票次数 |
	
	---
	### vote_option
    - collection: vote_option
    - 投票选项表

    >|字段|类型|说明|
    >|:----| :---- | :---- |
    >|VoteID   |string |   分享表唯一标识 |
    >|ID   |string |   投票选项表唯一标识 |
	>|Option   |string |   投票选项|
	>|OptionCount   |int |   投票被选次数 |
   
	---
	
	### public attachment
    - collection: public_attach
    - 公共附件

    >|字段|类型|说明|
    >|:----| :---- | :---- |
    >|ID   |string |   附件唯一标识 |
    >|Key  |string |   附件key|
    >|Type |number |   附件类型： 0-背景图片;1-控件;2-模板|
    >|Status|number|   附件状态: 0-状态正常；1-不可用|
    >|CreateTime   |string |   创建时间 |
    >|UpdateTime   |string |   上次更新时间 |

    --- 
    ### background
    - collection: background
    - 背景图片
    
    >|字段|类型|说明|
    >|:----| :---- | :---- |
    >|ID   |string |   背景唯一标识 |
    >|AttachID   |string |  public_attach唯一ID |
    >|Name |string |背景图片名称|
    >|CateName|string|背景图片分类名称|
    >|Price| nunbe | 定价 |
    >|Status|number|   状态: 0-状态正常；1-不可用|
    >|CreateTime   |string |   创建时间 |
    >|UpdateTime   |string |   上次更新时间 |

    --- 
    ### widget
    - collection: widget
    - 控件
    
    >|字段|类型|说明|
    >|:----| :---- | :---- |
    >|ID   |string |   唯一标识 |
    >|Name |string |控件名称|
    >|CateName|string|控件分类名称|
    >|Content| string| 控件内容|
    >|Price| nunbe | 定价 |
    >|Status|number|   状态: 0-状态正常；1-不可用|
    >|CreateTime   |string |   创建时间 |
    >|UpdateTime   |string |   上次更新时间 |

    --- 
    ### tempate
    - collection: template
    - 模板
    
    >|字段|类型|说明|
    >|:----| :---- | :---- |
    >|ID   |string |   唯一标识 |
    >|Name |string |模板名称|
    >|CateName|string|模板分类名称|
    >|Content| string| 模板内容|
    >|Price| nunbe | 定价 |
    >|Status|number|   状态: 0-状态正常；1-不可用|
    >|CreateTime   |string |   创建时间 |
    >|UpdateTime   |string |   上次更新时间 

	 --- 
    ### conflict
    - collection: conflict
    - 冲突历史表
    
    >|字段|类型|说明|
    >|:----| :---- | :---- |
    >|ID   |string |   唯一标识 |
    >|FatherPath |string |原笔记path|
    >|SonPath|string|新生成笔记path|
    >|CreateTime   |string |   创建时间 |
	