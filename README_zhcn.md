# Video ASCII Player
## 如何使用
1. 将代码克隆到本地
2. 安装 **[Go](https://golang.google.cn/doc/install)** 和 **FFmpeg**并设置好Path
    - **中国用户可能需要设置代理 [goproxy.io](https://goproxy.io/zh/)**
3. 构建可执行文件
    - 打开终端, 执行以下命令
        ```
        go mod download
        go build
        ```
4. 转换视频到图片
    - 打开终端, 执行以下命令
        ```
        ffmpeg -i {视频文件的路径} {存放每一帧图片的文件夹}/%d.png
        ```
5. 运行程序
    - 参数
        ```
        -buffer int
                渲染缓冲区的大小 [建议保持在 120 以上] (默认为 120)
        -count int
                视频的总帧数 [运行时从0到设定值] (默认为 -1)
        -fps int
                视频播放帧率 (默认为 30)
        -path string
                存放视频每一帧图片的文件夹路径
        -pool int
                渲染线程最大数 (默认为 16)
        ```
    - 例如: 存放视频每一帧图片的文件夹路径为 `C:/images/`, 视频播放帧率为 `30`, 视频的总帧数为 `6300` 帧. 所以你需要在终端输入以下命令来运行 `./VideoAsciiPlayer.exe -path=C:/images/ -count=6300 -fps=30`