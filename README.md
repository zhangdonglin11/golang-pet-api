# 项目简介

#### 这是一个仿宠物森林小程序和添加了论坛功能的后端项目。该项目是我曾使用SpringBoot和UniApp实现的毕业设计，
#### 现在该项目使用gin、gorm、、mysql等相关的技术栈重写了部分的接口。实现用户注册登录、发布宠物、宠物评论等功能，对增删改查的操作进行了实践。

完整的项目在我github的仓库cute_pet的<完整项目_.zip>中（uniapp实现的小程序代码、vue实现的后台管理、数据库的sql文件、后端项目的代码）

#### 项目主要技术

- golang
- gin
- gorm
- MySQL
- JWT
- redis
- rabbitMQ

#### 文件介绍

```
golang_pet_api 
├──/api/                      控制层API
├──/common/                   一些公共变量代码等
   ├── /config/               项目配置隐射
   ├── /constant/             一些前缀信息
   ├── /global/               全局公共变量
   ├── /initialize/           初始化
   ├── /result/               响应封装
   ├── /utils/                工具类
├──/img/                      项目展示的图片
├──/middleware/               中间件
├──/model/                    数据模型
├──/router/                   路由
├──/service/                  服务层
├──/static/                   项目静态文件
go.mod                        项目依赖
main.go                       程序主入口
settings.yaml                 项目配置文件
```

#### 使用方法
需要准备
- golang sdk 
- mysql 数据库
- redis 数据库
- RabbitMQ 服务

1. 将项目克隆到本地

   ```
   git clone https://github.com/JasonZhang1124/golang_pet_api.git
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

   ![](https://github.com/JasonZhang1124/golang-pet-api/blob/main/img/index.png)

   ![](https://github.com/JasonZhang1124/golang-pet-api/blob/main/img/topic.png)

   ![](https://github.com/JasonZhang1124/golang-pet-api/blob/main/img/topicAdd.png)

   ![](https://github.com/JasonZhang1124/golang-pet-api/blob/main/img/message.png)

   ![](https://github.com/JasonZhang1124/golang-pet-api/blob/main/img/messageDetail.png)

   ![](https://github.com/JasonZhang1124/golang-pet-api/blob/main/img/user.png)

   ![](https://github.com/JasonZhang1124/golang-pet-api/blob/main/img/userUpdate.png)

   ![](https://github.com/JasonZhang1124/golang-pet-api/blob/main/img/petDetail.png)

   ![](https://github.com/JasonZhang1124/golang-pet-api/blob/main/img/petAdd.png)

#### 后台管理效果展示

   ![](https://github.com/JasonZhang1124/golang-pet-api/blob/main/img/adminIndex.png)

   ![](https://github.com/JasonZhang1124/golang-pet-api/blob/main/img/adminPet.png)
   
   ![](https://github.com/JasonZhang1124/golang-pet-api/blob/main/img/adminUser.png)