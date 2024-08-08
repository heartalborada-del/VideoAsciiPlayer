# Video ASCII Player [中文文档](./README_zhcn.md)
## How to use 
1. Clone the code
2. Make sure you have installed **[Go](https://go.dev/doc/install)** and **FFmpeg**
    - **For China Users, you may need set the proxy [goproxy.io](https://goproxy.io/zh/)**
3. Build the executable file
    - Open terminal at the code folder, then run this script
        ```
        go mod download
        go build
        ```
4. Convert Video to Image 
    - Open terminal at the code folder, then run this script
        ```
        ffmpeg -i {Video path} {Image folder path}/%d.png
        ```
5. Running the program
    - Parameters
        ```
        -buffer int
                The Size of frame buffer [Recommended above 120] (default 120)
        -count int
                The Count of Frame Image [From 0 to Your setting] (default -1)
        -fps int
                Target FPS (default 30)
        -path string
                The Path of Stroage Video Every Frame Images
        -pool int
                Render Pool Size (default 16)
        ```
    - For example: Images store at `C:/images/`, video frame rate is `30`, the count of video frames is `6300`. So you should run this  script `./VideoAsciiPlayer.exe -path=C:/images/ -count=6300 -fps=30`