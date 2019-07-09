package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/astaxie/beego/config/env"
	"github.com/astaxie/beego/logs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//build镜像类型
const (
	//自己提供Dockerfile
	NO = iota
	//基于nginx
	NGINX
	//基于tomcat
	TOMCAT
	//基于jdk
	JAR
	//基于基础镜像
	DEFAULT
)

//镜像构建信息
type ImageSimpleBuildInfo struct {
	Name,
	Version,
	Prefix string
	Type int
}

//json类型
type Json string

//dockerfile信息
type DockerfileInfo struct {
	BaseImage, UserName, Email, Location, Port, Cmd string
}

//创建docker客户端
func NewDockerClient() *client.Client {
	defer func() {
		if e := recover(); e != nil {
			logs.Error("初始化错误 %s\r\n", e)
		}
	}()
	env := env.Get("ENV", "remote")
	if env != "remote" {
		dockerClient, e := client.NewEnvClient()
		if e != nil {
			panic(e)
		}
		logs.Info("docker.sock 模式")
		return dockerClient;
	} else {
		httpHead := make(map[string]string)
		dockerClient, e := client.NewClient("http://192.168.31.185:2375", "1.39", nil, httpHead)
		if e != nil {
			panic(e)
		}
		logs.Info("远程模式")
		return dockerClient
	}
}

//build镜像
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

//镜像构建前置函数
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

//构建
func doBuild(reader *io.Reader, simpleBuildInfo *ImageSimpleBuildInfo, ch chan<- struct{}) {
	c := NewDockerClient()
	defer c.Close()
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
	defer body.Close()
	bytes, e := ioutil.ReadAll(body)
	if e != nil {
		panic(e)
	}
	logs.Info(string(bytes))
	ch <- struct{}{}
}

//判断是否为空
func IsNotNil(item interface{}) bool {
	return item != nil
}

//生成dockerfile模板
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

//拉取镜像
func DockerImagePull(image string) {
	c := NewDockerClient()
	defer c.Close()
	defer func() {
		if e := recover(); e != nil {
			logs.Error("docker pull 错误 %s\r\n", e)
		}
	}()
	authConfig := types.AuthConfig{
		Username:      "admin",
		Password:      "Harbor12345",
		ServerAddress: "harbor.self.com",
		Email:         "admin@example.com",
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	rep, e := c.ImagePull(context.Background(), image, types.ImagePullOptions{
		RegistryAuth: authStr,
	})
	if e != nil {
		panic(e)
	}
	defer rep.Close()
	bytes, e := ioutil.ReadAll(rep)
	if e != nil {
		panic(e)
	}
	logs.Info(string(bytes))
}

//查询所有镜像
func DockerImageList() []types.ImageSummary {
	defer func() {
		if e := recover(); e != nil {
			logs.Error("docker list 错误 %s\r\n", e)
		}
	}()
	dockerClient := NewDockerClient()
	defer dockerClient.Close()
	imageList, e := dockerClient.ImageList(context.Background(), types.ImageListOptions{})
	if e != nil {
		panic(e)
	}
	return imageList
}

func DockerImageQuery() {
	/*	defer func() {
		if e := recover(); e != nil {
			close(ch)
			logs.Error("docker push 错误 %s\r\n", e)
		}
	}()*/
}

//删除镜像
func DockerImageDelete(id string) []types.ImageDelete {
	defer func() {
		if e := recover(); e != nil {
			logs.Error("docker push 错误 %s\r\n", e)
		}
	}()
	dockerClient := NewDockerClient()
	defer dockerClient.Close()
	remove, e := dockerClient.ImageRemove(context.Background(), id, types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: false,
	})
	if e != nil {
		panic(e)
	}
	return remove
}

func DockerImagePush(image string, ch chan<- struct{}) {
	defer func() {
		if e := recover(); e != nil {
			close(ch)
			logs.Error("docker push 错误 %s\r\n", e)
		}
	}()
	c := NewDockerClient()
	defer c.Close()
	authConfig := types.AuthConfig{
		Username:      "admin",
		Password:      "Harbor12345",
		ServerAddress: "harbor.self.com",
		Email:         "admin@example.com",
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	rsp, err := c.ImagePush(context.Background(), image, types.ImagePushOptions{RegistryAuth: authStr,})
	if err != nil {
		panic(err)
	}
	defer rsp.Close()
	bytes, e := ioutil.ReadAll(rsp)
	if e != nil {
		panic(e)
	}
	logs.Info(string(bytes))
	ch <- struct{}{}
}
