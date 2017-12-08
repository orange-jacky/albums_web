package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	. "github.com/orange-jacky/albums/data"
	. "github.com/orange-jacky/albums_web/util"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type UserPasswd struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *UserPasswd) Marshal() (string, error) {
	result, err := json.Marshal(s)
	return string(result), err
}

func signinHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	u := &UserPasswd{}
	u.Username = username
	u.Password = password

	postForm, _ := u.Marshal()
	buff := strings.NewReader(postForm)

	conf := GetConfigure()
	b := conf.Backend
	url := fmt.Sprintf("%s%s", b.Host, b.Url.Login)
	req, err := http.NewRequest("POST", url, buff)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%v", err))
		return
	}
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%v", err))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%v", err))
		return
	}
	c.String(http.StatusOK, "%v", string(body))
}

func signupHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	postForm := url.Values{}
	postForm.Add("username", username)
	postForm.Add("password", password)
	buff := strings.NewReader(postForm.Encode())

	conf := GetConfigure()
	b := conf.Backend
	url := fmt.Sprintf("%s%s", b.Host, b.Url.Signup)
	req, err := http.NewRequest("POST", url, buff)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%v", err))
		return
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%v", err))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%v", err))
		return
	}
	c.String(http.StatusOK, "%v", string(body))
}

func loginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func homeHandler(c *gin.Context) {
	conf := GetConfigure()
	c.Redirect(http.StatusMovedPermanently, conf.GinServer.Url+"/login")
}

func contentHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "content.html", gin.H{})
}

//curl -F "image"=@"IMAGEFILE" -F "key"="KEY" URL
//https://stackoverflow.com/questions/20205796/golang-post-data-using-the-content-type-multipart-form-data
func uploadFile(c *gin.Context, router string, params map[string]string) {
	//username, _ := c.Cookie("username")
	token, _ := c.Cookie("token")
	//fmt.Printf("odHandler{username:%v,token:%v}\n", username, token)
	r := c.Request
	//POST takes the uploaded file(s) and saves it to disk.
	//parse the multipart form in the request
	err := r.ParseMultipartForm(100000)
	if err != nil {
		c.String(http.StatusOK, "ParseMultipartForm %v", err)
		return
	}
	//get a ref to the parsed multipart form
	m := r.MultipartForm
	//get the *fileheaders
	files := m.File["upload_file"]
	//post 没有文件直接返回
	if len(files) == 0 {
		c.String(http.StatusOK, "upload_file empty")
		return
	}
	file := files[0]
	fp, err := file.Open()
	if err != nil {
		c.String(http.StatusOK, "files[0].Open() %v", err)
		return
	}
	defer fp.Close()

	var buff bytes.Buffer
	writer := multipart.NewWriter(&buff)
	fw, err := writer.CreateFormFile("image", file.Filename)
	if err != nil {
		// Don't forget to close the multipart writer.
		// If you don't close it, your request will be missing the terminating boundary.
		writer.Close()
		c.String(http.StatusOK, "writer.CreateFormFile %v", err)
		return
	}
	if _, err = io.Copy(fw, fp); err != nil {
		// Don't forget to close the multipart writer.
		// If you don't close it, your request will be missing the terminating boundary.
		writer.Close()
		c.String(http.StatusOK, "io.Copy %v", err)
		return
	}
	for key, val := range params {
		writer.WriteField(key, val)
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	writer.Close()

	conf := GetConfigure()
	b := conf.Backend
	url := fmt.Sprintf("%s%s", b.Host, router)
	req, err := http.NewRequest("POST", url, &buff)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("http.NewRequest %v", err))
		return
	}

	value := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", value)
	req.Header.Set("Content-type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("client.Do %v", err))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("ioutil.ReadAll %v", err))
		return
	}
	c.String(http.StatusOK, "%v", string(body))
}

func objectdectionHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "objectdection.html", gin.H{})
}

func odHandler(c *gin.Context) {
	conf := GetConfigure()

	b := conf.Backend
	router := b.Url.Objectdetection_dl
	uploadFile(c, router, nil)
}

func deeplearningHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "deeplearning.html", gin.H{})
}

func dpHandler(c *gin.Context) {
	conf := GetConfigure()
	b := conf.Backend
	router := b.Url.Deeplearning
	uploadFile(c, router, nil)
}

func uploadMultiFiles(c *gin.Context, router string, params map[string]string) {
	token, _ := c.Cookie("token")

	//fmt.Printf("odHandler{username:%v,token:%v}\n", username, token)
	r := c.Request
	//POST takes the uploaded file(s) and saves it to disk.
	//parse the multipart form in the request
	err := r.ParseMultipartForm(100000)
	if err != nil {
		c.String(http.StatusOK, "ParseMultipartForm %v", err)
		return
	}
	//get a ref to the parsed multipart form
	m := r.MultipartForm
	//get the *fileheaders
	files := m.File["upload_file"]
	//post 没有文件直接返回
	if len(files) == 0 {
		c.String(http.StatusOK, "upload_file empty")
		return
	}

	var buff bytes.Buffer
	writer := multipart.NewWriter(&buff)
	for _, file := range files {
		fp, err := file.Open()
		if err != nil {
			continue
		}
		//
		fw, err := writer.CreateFormFile("images", file.Filename)
		if err != nil {
			// Don't forget to close the multipart writer.
			// If you don't close it, your request will be missing the terminating boundary.
			fp.Close()
			continue
		}
		if _, err = io.Copy(fw, fp); err != nil {
			fp.Close()
			continue
		}
		fp.Close()
	}
	for key, val := range params {
		writer.WriteField(key, val)
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	writer.Close()

	conf := GetConfigure()
	b := conf.Backend
	url := fmt.Sprintf("%s%s", b.Host, router)
	req, err := http.NewRequest("POST", url, &buff)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("http.NewRequest %v", err))
		return
	}

	value := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", value)
	req.Header.Set("Content-type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("client.Do %v", err))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("ioutil.ReadAll %v", err))
		return
	}
	c.String(http.StatusOK, "%v", string(body))
}

func getAlbum(c *gin.Context) string {
	album, _ := c.Cookie("album")
	if album == "" {
		album = "default"
	}
	return album
}

func uploadHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{})
}

func ulHandler(c *gin.Context) {
	conf := GetConfigure()
	b := conf.Backend
	router := b.Url.Upload

	album := getAlbum(c)
	username, _ := c.Cookie("username")
	if album == "" || username == "" {
		c.String(http.StatusOK, "username or album empty")
		return
	}
	params := map[string]string{
		"username": username,
		"album":    album,
	}
	uploadMultiFiles(c, router, params)
}

func searchHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "search.html", gin.H{})
}

func sHandler(c *gin.Context) {
	conf := GetConfigure()
	b := conf.Backend
	router := b.Url.Search

	album := getAlbum(c)
	username, _ := c.Cookie("username")
	if album == "" || username == "" {
		c.String(http.StatusOK, "username or album empty")
		return
	}
	params := map[string]string{
		"username": username,
		"album":    album,
	}
	uploadFile(c, router, params)
}

func GetList(c *gin.Context, router string, params map[string]string) {
	token, _ := c.Cookie("token")

	var buff bytes.Buffer
	writer := multipart.NewWriter(&buff)
	for key, val := range params {
		writer.WriteField(key, val)
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	writer.Close()

	conf := GetConfigure()
	b := conf.Backend
	url := fmt.Sprintf("%s%s", b.Host, router)
	req, err := http.NewRequest("POST", url, &buff)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("http.NewRequest %v", err))
		return
	}

	value := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", value)
	req.Header.Set("Content-type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("client.Do %v", err))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("ioutil.ReadAll %v", err))
		return
	}
	c.String(http.StatusOK, "%v", string(body))
}

func albumlistHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "albumlist.html", gin.H{})
}

func alistHandler(c *gin.Context) {
	conf := GetConfigure()
	b := conf.Backend
	router := b.Url.Albummgt.Get
	username, _ := c.Cookie("username")
	if username == "" {
		c.String(http.StatusOK, "username empty")
		return
	}
	params := map[string]string{
		"username": username,
	}
	GetList(c, router, params)
}

func downloadHandler(c *gin.Context) {
	conf := GetConfigure()
	b := conf.Backend
	router := b.Url.Download

	album := getAlbum(c)
	username, _ := c.Cookie("username")
	if album == "" || username == "" {
		c.String(http.StatusOK, "username or album empty")
		return
	}
	params := map[string]string{
		"username": username,
		"album":    album,
	}
	GetList(c, router, params)
}

////////////////////////////////////////////////////////////////////////////

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func index2Handler(c *gin.Context) {
	c.HTML(http.StatusOK, "index2.html", gin.H{})
}

func index2newHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index2new.html", gin.H{})
}

func bHandler(c *gin.Context) {
	sli := make([]string, 0)
	sli = append(sli, "http://47.104.21.233:80/image/8c0ff2091096e52fbbf01d702f2d8f70")
	sli = append(sli, "http://47.104.21.233:80/image/e5a7251c1dcbf771d7dfeef5a15a15b9")
	sli = append(sli, "http://47.104.21.233:80/image/47ac1b3003e6903a211acb0950392114")
	sli = append(sli, "http://47.104.21.233:80/image/499a4d27d38e5cf8edbba501a9b000f2")
	c.HTML(http.StatusOK, "b.html", sli)
}

func b2Handler(c *gin.Context) {
	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	writer.WriteField("username", "admin")
	writer.WriteField("album", "default")
	writer.Close()
	req, err := http.NewRequest("POST", "http://47.104.21.233:9000/auth/download?page=1&size=20", buff)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%v", err))
		return
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTE5NTIwNjcsImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTUxMTk0ODQ2N30.L_tyRJQTkEc-vA37TJT2h9MvdrqXHwd6us0MoWS3-xg"
	value := fmt.Sprintf("Bearer %s", token)
	req.Header.Add("Authorization", value)
	req.Header.Set("Content-type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%v", err))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%v", err))
		return
	}
	ret := &Response{}
	if err := ret.Unmarshal(body); err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%v", err))
		return
	}
	sli, ok := ret.Data.([]interface{})
	if !ok {
		c.String(http.StatusOK, fmt.Sprintf("ret.Data.(ImageInfos) fail"))
		return
	}
	var imgs ImageInfos
	for _, s := range sli {
		//fmt.Printf("%T", s)
		img := &ImageInfo{}
		if err := mapstructure.Decode(s, img); err == nil {
			imgs = append(imgs, img)
		}
	}
	c.HTML(http.StatusOK, "b2.html", imgs)
}

func cHandler(c *gin.Context) {
	type A struct {
		A, B, C, D, E string
	}
	sli := make([]*A, 0)
	for i := 0; i < 10; i++ {
		a := &A{}
		a.A = fmt.Sprintf("%v", i)
		a.B = fmt.Sprintf("%v b", i)
		a.C = fmt.Sprintf("%v c", i)
		a.D = fmt.Sprintf("%v d", i)
		a.E = fmt.Sprintf("%v e", i)
		sli = append(sli, a)
	}
	c.HTML(http.StatusOK, "c.html", sli)
}
