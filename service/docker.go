package service

import (
	"context"
	"fmt"
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
	_ = iota
	NGINX
	TOMCAT
	JAR
	SELF
)

type ImageSimpleBuildInfo struct {
	Name,
	Version,
	Prefix string
	Type int
}

var c *client.Client

func init() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("初始化错误 %s\r\n", e)
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
			fmt.Printf("docker build 错误 %s\r\n", e)
		}
	}()
	var reader *io.Reader
	switch simpleBuildInfo.Type {
	case NGINX:
		reader = nginxBuild(file, header)
	case JAR:
		reader = jarBuild(file, header)
	case TOMCAT:
	case SELF:
		reader = file
	default:
	}
	build(reader, simpleBuildInfo, ch)
}

func nginxBuild(file *io.Reader, header *multipart.FileHeader) *io.Reader {
	source := rand.NewSource(time.Now().UnixNano())
	rd := rand.New(source)
	tempDir := os.TempDir() + "\\docker" + strconv.Itoa(rd.Int())
	os.Mkdir(tempDir, os.ModePerm)
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
		var t io.Writer = target
		io.Copy(t, *file)
		if e != nil {
			panic(e)
		}
		target.Close()
	}
	dockerFileGen()
	//var w io.Writer = buffer
	//Tars(dir, &w)
	err := Tar("./", "docker.tar", false)
	if err != nil {
		panic(err)
	}
	//var reader io.Reader = bytes.NewReader(buffer.Bytes())
	//return &reader;
	f, err := os.Open("docker.tar")
	if err != nil {
		panic(err)
	}
	var t io.Reader = f
	return &t
}

func dockerFileGen() {
	i := template.New("Dockerfile")
	parse, e := i.Parse(`FROM  {{.BaseImage}}
	MAINTAINER  {{.UserName}}  {{.Email}}
	ADD  ./*   /usr/share/nginx/html/
	EXPOSE  80`)
	if e != nil {
		panic(e)
	}
	dockerfile, err := os.Create("Dockerfile")
	if err != nil {
		panic(err)
	}
	//var byteTemp []byte
	//buffer := bytes.NewBuffer(byteTemp)
	err = parse.Execute(dockerfile, struct {
		BaseImage, UserName, Email string
	}{"nginx", "wangzhu", "wangzhu@originaltek.com"})
	if err != nil {
		panic(err)
	}
	dockerfile.Close()
}

func tomcatBuild(file *io.Reader) *io.Reader {
	return nil
}

func jarBuild(file *io.Reader, header *multipart.FileHeader) *io.Reader {

	return nil
}

func build(file *io.Reader, simpleBuildInfo *ImageSimpleBuildInfo, ch chan<- struct{}) {
	buildResponse, e := c.ImageBuild(context.Background(), *file, types.ImageBuildOptions{
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
