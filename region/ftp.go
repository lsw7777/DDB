package region

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jlaffaye/ftp"

	. "Distributed-MiniSQL/common"
)

//定义Ftp结构体变量
type FtpUtils struct {
	port      string		//端口号
	username  string		//用户名
	password  string		//密码
	ftpClient *ftp.ServerConn	//客户端对象
}

//用于设置 FTP 连接参数
func (fu *FtpUtils) Construct() {
	fu.port = "21"
	fu.username = "lsw"
	fu.password = "lsw"
}

//用于登陆 FTP 服务器
func (fu *FtpUtils) Login(IP string) {
	c, err := ftp.Dial(IP+":"+fu.port, ftp.DialWithTimeout(5*time.Second), ftp.DialWithDisabledEPSV(true))
	if err != nil {
		fmt.Printf("[from ftputils]ftp连接出现问题: %v\n", err)
	}
	err = c.Login(fu.username, fu.password)
	if err != nil {
		fmt.Printf("[from ftputils]登录出现问题: %v\n", err)
	}
	fu.ftpClient = c
}

//用于关闭 FTP 连接
func (fu *FtpUtils) CloseConnect() {
	err := fu.ftpClient.Quit()
	if err != nil {
		fmt.Printf("[from ftpUtils]ftp断开连接出现问题: %v\n", err)
	}
}

//用于在FTP服务器下载文件
func (fu *FtpUtils) DownloadFile(remoteFileName string, localFileName string, appendOrNot bool, IP string) bool {
	fu.Login(IP)
	fetchfile, _ := fu.ftpClient.Retr(remoteFileName)
	defer fetchfile.Close()
	var localfile *os.File
	var err error
	if appendOrNot {
		localfile, err = os.OpenFile(localFileName, os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("%v", err)
		} else {
			fmt.Println("[hint]append ok")
		}
	} else {
		localfile, err = os.OpenFile(localFileName, os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			fmt.Printf("%v", err)
		} else {
			fmt.Println("[hint]create new file ok")
		}
	}
	defer localfile.Close()
	io.Copy(localfile, fetchfile)
	fu.CloseConnect()
	return true
}

//用于在FTP服务器上传文件
func (fu *FtpUtils) UploadFile(localFileName string, remoteFileName string, IP string) bool {
	fu.Login(IP)
	file, err := os.Open(localFileName)
	if err != nil {
		fmt.Printf("[from ftputils]read local file failed: %v\n", err)
		return false
	}
	defer file.Close()
	err = fu.ftpClient.Stor(remoteFileName, file)
	if err != nil {
		fmt.Printf("[from ftputils]uploading file failed: %v\n", err)
		return false
	}
	
	fu.CloseConnect()
	return true
}

//用于在FTP服务器删除文件
func (fu *FtpUtils) DeleteFile(remoteFileName string, IP string) bool {
	fu.Login(IP)

	cur, err := fu.ftpClient.CurrentDir()
	if err != nil {
		log.Printf("current dir fail")
	}
	log.Printf("[hint]current path in ftp: %v", cur)
	err = fu.ftpClient.Delete(remoteFileName)
	if err != nil {
		fmt.Printf("[from ftputils]delete file failed: %v\n", err)
		return false
	}

	fu.CloseConnect()
	return true
}

//用于在FTP服务器下载目录
func (fu *FtpUtils) DownloadDir(remoteDir, localDir, ip string) {
	CleanDir(localDir)
	// download everything from backup's sql dir to local sql dir
	fu.Login(ip)
	err := fu.ftpClient.ChangeDir(remoteDir)
	if err != nil {
		fmt.Println("[from ftputils]ftpPath not exist")
	}
	ftpFiles, e := fu.ftpClient.List("./") 
	if e != nil {
		fmt.Printf("[from backup ftp]ftpfiles list fail: %v\n", e)
	}
	if ftpFiles == nil {
		fmt.Println("[from backup ftp]list下无文件")
	}
	var fileNameArray []string
	fileNameArray = make([]string, 0)
	for _, file := range ftpFiles {
		fileNameArray = append(fileNameArray, file.Name)
	}
	fu.CloseConnect()

	for _, fileName := range fileNameArray {
		fu.AnotherDownload(ip, remoteDir, localDir, fileName)
	}
}

//用于在FTP服务器下载文件
func (fu *FtpUtils) AnotherDownload(IP string, remoteDir string, localDir string, fileName string) {
	fu.Login(IP)
	err := fu.ftpClient.ChangeDir(remoteDir)
	if err != nil {
		fmt.Println("[from ftputils]ftpPath not exist")
	}
	fetchfile, ferr := fu.ftpClient.Retr(fileName)
	if ferr != nil {
		log.Printf("%v", ferr)
	}
	localfile, err := os.Create(localDir + fileName)
	if err != nil {
		log.Printf("%v", err)
	}
	_, err = io.Copy(localfile, fetchfile)
	if err != nil {
		log.Printf("fail to copy")
	}
	defer fetchfile.Close()
	defer localfile.Close()
	fu.CloseConnect()
}
