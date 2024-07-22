# 项目简介
#### 本项目是仿宠物森林小程序的实现并且增加了话题论坛功能的社交平台，可以在这里可以发布宠物交 友、宠物配对、宠物代溜、宠物领养。同时也提供了宠物话题社区，分享宠友们与爱宠的沙雕日常！整个系统分为小程序前端、系统后台前端、系统后端 3 个项目。

主要功能：
   用户模块：设计并实现用户注册、登录、信息管理等功能，宠物模块：实现添加、修改、删除宠物功能，
设计分类查询宠物、评论和点赞功能。话题社区模块：实现论坛基本功能，如发帖、删贴、评论和点赞。对
话功能：使用websocket实现用户之间的一对一消息发送功能，支持文字、图片和宠物信息的发送。
完整的项目在我github的仓库cute_pet的<完整项目_.zip>中（uniapp实现的小程序代码、vue实现的后台管理、数据库的sql文件、后端项目的代码）

#### 项目主要技术

- golang
- gin
- gorm
- MySQL
- JWT
- redis
- websocket

#### 使用方法
需要准备
- golang sdk 
- mysql 数据库
- redis 数据库

1. 将项目克隆到本地

   ```
   git clone https://github.com/zhangdonglin11/golang_pet_api.git
   ```

2. 在项目中初始化所需的依赖
   ```
   go mod tidy
   ```

3. 使用models文件中的golang_pet_db.sql导入数据库
4. RabbitMQ 需安装插件
   ```
   rabbitmq_delayed_message_exchange
   ```
5. 在命令窗口启动该项目

    ```
   PS .../golang-pet-api> go run main.go
   ```

#### 小程序端效果展示

   ![](https://github.com/zhangdonglin11/golang-pet-api/blob/main/img/index.png)

   ![](https://github.com/zhangdonglin11/golang-pet-api/blob/main/img/topic.png)

   ![](https://github.com/zhangdonglin11/golang-pet-api/blob/main/img/topicAdd.png)

   ![](https://github.com/zhangdonglin11/golang-pet-api/blob/main/img/message.png)

   ![](https://github.com/zhangdonglin11/golang-pet-api/blob/main/img/messageDetail.png)

   ![](https://github.com/zhangdonglin11/golang-pet-api/blob/main/img/user.png)

   ![](https://github.com/zhangdonglin11/golang-pet-api/blob/main/img/userUpdate.png)

   ![](https://github.com/zhangdonglin11/golang-pet-api/blob/main/img/petDetail.png)

   ![](https://github.com/zhangdonglin11/golang-pet-api/blob/main/img/petAdd.png)

#### 后台管理效果展示

   ![](https://github.com/zhangdonglin11/golang-pet-api/blob/main/img/adminIndex.png)

   ![](https://github.com/zhangdonglin11/golang-pet-api/blob/main/img/adminPet.png)
   
   ![](https://github.com/zhangdonglin11/golang-pet-api/blob/main/img/adminUser.png)
