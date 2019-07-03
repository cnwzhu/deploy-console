package service

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	NO = iota
	NGINX
	TOMCAT
	JAR
	DEFAULT
)

type ImageSimpleBuildInfo struct {
	Name,
	Version,
	Prefix string
	Type int
}

type DockerfileInfo struct {
	BaseImage, UserName, Email, Location, Port, Cmd string
}

var c *client.Client

func init() {
	defer func() {
		if e := recover(); e != nil {
			logs.Error("初始化错误 %s\r\n", e)
		}
	}()
	httpHead := make(map[string]string)
	dockerClient, e := client.NewClient("http://192.168.31.188:2375", "1.39", nil, httpHead)
	if e != nil {
		panic(e)
	}
	c = dockerClient
}

func DockerImageBuild(file *io.Reader, simpleBuildInfo *ImageSimpleBuildInfo, ch chan<- struct{}, header *multipart.FileHeader) {
	defer func() {
		if e := recover(); e != nil {
			close(ch)
			logs.Error("docker build 错误 %s\r\n", e)
		}
	}()
	if simpleBuildInfo.Type == NO {
		doBuild(file, simpleBuildInfo, ch)
	} else {
		reader := preBuild(file, header, simpleBuildInfo.Type)
		doBuild(reader, simpleBuildInfo, ch)
	}
}

func preBuild(file *io.Reader, header *multipart.FileHeader, ty int) *io.Reader {
	source := rand.NewSource(time.Now().UnixNano())
	rd := rand.New(source)
	tempDir := os.TempDir() + "\\docker" + strconv.Itoa(rd.Int())
	err := os.Mkdir(tempDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	e := os.Chdir(tempDir)
	if e != nil {
		panic(e)
	}
	if strings.LastIndex(header.Filename, ".tar") > 0 {
		err := UnTarFiles(file, "./")
		if err != nil {
			panic(err)
		}
	} else {
		target, e := os.Create(header.Filename)
		if e != nil {
			panic(e)
		}
		var t io.Writer = target
		_, e = io.Copy(t, *file)
		if e != nil {
			panic(e)
		}
		e = target.Close()
		if e != nil {
			panic(e)
		}
	}
	dockerFileGen(ty, header.Filename)
	err = Tar("./", "docker.tar", false)
	if err != nil {
		panic(err)
	}
	f, err := os.Open("docker.tar")
	if err != nil {
		panic(err)
	}
	var t io.Reader = f
	return &t
}

func doBuild(reader *io.Reader, simpleBuildInfo *ImageSimpleBuildInfo, ch chan<- struct{}) {
	buildResponse, e := c.ImageBuild(context.Background(), *reader, types.ImageBuildOptions{
		Tags:           []string{simpleBuildInfo.Prefix + "/" + simpleBuildInfo.Name + ":" + simpleBuildInfo.Version},
		SuppressOutput: true,
		NoCache:        false,
		Remove:         true,
		ForceRemove:    false,
		PullParent:     true,
		Labels:         map[string]string{},
	})
	if e != nil {
		panic(e)
	}
	body := buildResponse.Body
	io.Copy(os.Stdout, body)
	defer body.Close()
	ch <- struct{}{}
}
func IsNotNil(item interface{}) bool {
	return item != nil
}
func dockerFileGen(ty int, name string) {
	username := "wangzhu"
	email := "wang-zhu@live.com"
	logs.Info("username: %v ,email: %v", username, email)
	dockerfile, err := os.Create("Dockerfile")
	if err != nil {
		panic(err)
	}
	i := template.New("Dockerfile")
	fm := template.FuncMap{"IsNotNil": IsNotNil}
	i.Funcs(fm)
	parse, e := i.Parse(`FROM  {{.BaseImage}}
MAINTAINER  {{.UserName}}  {{.Email}}
ADD  ./*   {{.Location}}
EXPOSE  {{.Port}} {{if IsNotNil .Cmd }}
CMD  {{.Cmd}}{{end}}`)
	if e != nil {
		panic(e)
	}
	var info DockerfileInfo
	switch ty {
	case NGINX:
		info = DockerfileInfo{BaseImage: "nginx", UserName: username, Email: email, Location: "/usr/share/nginx/html/", Port: "80",}
	case JAR:
		info = DockerfileInfo{BaseImage: "openjdk:8", UserName: username, Email: email, Location: "/root/", Port: "80", Cmd: "java -jar /root/" + name,}
	case TOMCAT:
		info = DockerfileInfo{BaseImage: "tomcat:jdk8", UserName: username, Email: email, Location: " /usr/local/tomcat/webapps/", Port: "8080",}
	case DEFAULT:
		info = DockerfileInfo{BaseImage: "alpine", UserName: username, Email: email, Location: " /root/", Port: "80", Cmd: "sh  /root/" + name,}
	default:
		panic("不支持的打包类型")
	}
	err = parse.Execute(dockerfile, info)
	if err != nil {
		panic(err)
	}
	dockerfile.Close()
}

func DockerImagePull(image string) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("docker pull 错误 %s\r\n", e)
		}
	}()
	rep, e := c.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if e != nil {
		panic(e)
	}
	io.Copy(os.Stdout, rep)
	defer rep.Close()
}

func DockerImageList() {

}

func DockerImageQuery() {

}

func DockerImagePush() {

}
