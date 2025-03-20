# DyParser

DyParser 是一个使用 Go 语言开发的工具，用于从抖音和 TikTok 的分享链接中提取信息。

## 功能

- 从抖音分享链接中提取视频信息
- 从 TikTok 分享链接中提取视频信息
- 支持通过代理访问

## 安装

1. 克隆仓库：

    ```sh
    git clone https://github.com/seanchan/dyparser.git
    ```

2. 进入项目目录：

    ```sh
    cd dyparser
    ```

3. 安装依赖：

    ```sh
    go mod tidy
    ```

## 配置

在项目根目录下创建 `config.yaml` 文件，并添加以下内容：

```yaml
server:
  port: "8081"
proxy:
  url: "http://127.0.0.1:7890"
```

## 使用
运行项目
```
go run main.go
```

服务器启动后，可以通过以下 API 进行访问：

```
GET /parse?query=<query>&source=<source> - 解析分享链接，返回视频信息
```
## 示例
请求
```
curl "http://localhost:8081/parse?query=<query>&source=<source>"
curl "http://localhost:8081/parse?source=tiktok&query=https://www.tiktok.com/@naploes/video/7480531425259834642?is_from_webapp=1&sender_device=pc"
```
响应
```
{"code":200,"msg":"","data":{"author":{"uid":"","name":"","avatar":"","source_data":null},"author_src":{"data":null},"video":{"title":"","video_url":"","music_url":"","cover_url":"","images":null,"video_src":{"data":null}},"video_src":"省略"}}}
```

## 贡献
欢迎提交 Issue 和 Pull Request 来帮助改进这个项目。

## 许可证
本项目使用 MIT 许可证，详情请参阅 LICENSE 文件。