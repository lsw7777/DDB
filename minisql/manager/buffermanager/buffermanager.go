package buffermanager

import (
	"Distributed-MiniSQL/common"
	"fmt"
	"io"
	"os"
	"unsafe"
)

const BLOCKSIZE = 4096

//定义数据块结构体变量
type Block struct {
	LRUCount    int
	BlockOffset int
	IsDirty     bool
	IsValid     bool
	IsLocked    bool
	FileName    string
	BlockData   []byte
}

//分配新的数据块
func NewBlock() *Block {
	return &Block{LRUCount: 0, BlockOffset: 0, IsDirty: false, IsValid: false, IsLocked: false, BlockData: make([]byte, BLOCKSIZE)}
}

//写入数据
func (b *Block) WriteData(offset int, data []byte) bool {
	if offset+len(data) > BLOCKSIZE {
		return false
	}
	for i := 0; i < len(data); i++ {
		b.BlockData[i+offset] = data[i]
	}
	b.IsDirty = true
	b.LRUCount++
	return true
}

//将模式各变量设为0
func (b *Block) ResetModes() {
	b.IsDirty = false
	b.IsLocked = false
	b.IsValid = false
	b.LRUCount = 0
}

//读取整数
func (b *Block) ReadInteger(offset int) int32 {
	if offset+4 > BLOCKSIZE {
		return 0
	}
	var b0 int32 = int32(b.BlockData[offset] & 0xFF)
	var b1 int32 = int32(b.BlockData[offset+1] & 0xFF)
	var b2 int32 = int32(b.BlockData[offset+2] & 0xFF)
	var b3 int32 = int32(b.BlockData[offset+3] & 0xFF)
	b.LRUCount++
	return (b0 << 24) | (b1 << 16) | (b2 << 8) | b3
}

//写入字节
func (b *Block) WriteByte(offset int, value byte) bool {
	if offset+1 > BLOCKSIZE {
		return false
	}
	b.BlockData[offset] = value
	b.IsDirty = true
	return true
}

//写入整数
func (b *Block) WriteInteger(offset int, value int32) bool {
	if offset+4 > BLOCKSIZE {
		return false
	}
	b.BlockData[offset] = (byte)((value >> 24) & 0xFF)
	b.BlockData[offset+1] = (byte)((value >> 16) & 0xFF)
	b.BlockData[offset+2] = (byte)((value >> 8) & 0xFF)
	b.BlockData[offset+3] = (byte)(value & 0xFF)
	b.LRUCount++
	b.IsDirty = true
	return true
}

//读取浮点数
func (b *Block) ReadFloat(offset int) float32 {
	var dat int32 = b.ReadInteger(offset)
	return *(*float32)(unsafe.Pointer(&dat))
}

//写入浮点数
func (b *Block) WriteFloat(offset int, value float32) bool {
	var dat int32 = *(*int32)(unsafe.Pointer(&value))
	b.IsDirty = true
	return b.WriteInteger(offset, dat)
}

//读取字符串
func (b *Block) ReadString(offset int, length int) string {
	var buf []byte = make([]byte, length)
	for i := 0; i < length && (i < BLOCKSIZE-offset); i++ {
		buf[i] = b.BlockData[offset+i]
	}
	b.LRUCount++
	return string(buf)
}

//写入字符串
func (b *Block) WriteString(offset int, str string) bool {
	var buf []byte = []byte(str)
	if offset+len(buf) > BLOCKSIZE {
		return false
	}
	for i := 0; i < len(buf); i++ {
		b.BlockData[offset+i] = buf[i]
	}
	b.LRUCount++
	b.IsDirty = true
	return true
}

//给数据块赋初值
func (b *Block) SetBlockData() {
	for i := 0; i < len(b.BlockData); i++ {
		b.BlockData[i] = 0
	}
}

var MAXBLOCKNUM int = 50
var EOF int = -1
var buffer []Block = make([]Block, MAXBLOCKNUM)

//缓冲器初始化
func BufferInit() {
	for i := 0; i < MAXBLOCKNUM; i++ {
		buffer[i] = *NewBlock()
	}
}

//测试接口
func TestInterFace() {
	b := NewBlock()
	b.WriteInteger(1200, 2245)
	b.WriteFloat(76, 2232.14)
	b.WriteString(492, "hwggnb!!")
	b.FileName = "buffer_test"
	b.BlockOffset = 15
	buffer[1] = *b
	WriteBlockToDisk(1)
}

//销毁缓冲管理器
func DestructBufferManager() {
	for i := 0; i < MAXBLOCKNUM; i++ {
		if buffer[i].IsValid {
			WriteBlockToDisk(i)
		}
	}
}

//使某文件不可用（锁定）
func MakeInvalid(filename string) {
	for i := 0; i < MAXBLOCKNUM; i++ {
		if buffer[i].FileName == filename {
			buffer[i].IsValid = false
		}
	}
}

//从磁盘读出数据块
func ReadBlockFromDisk(filename string, ofs int) int {
	for i := 0; i < MAXBLOCKNUM; i++ {
		if buffer[i].IsValid && buffer[i].FileName == filename && buffer[i].BlockOffset == ofs {
			return i
		}
	}
	bid := GetFreeBlockId()
	if bid == EOF {
		return EOF
	}
	f, err := os.OpenFile(common.DIR+filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Print(err)
		return EOF
	}
	if !ReadBlockFromDisk1(filename, ofs, bid) {
		return EOF
	}
	defer f.Close()
	return bid
}

//从磁盘队列读出数据块
func ReadBlockFromDiskQuote(filename string, ofs int) *Block {
	i := 0
	for ; i < MAXBLOCKNUM; i++ {
		if buffer[i].IsValid && buffer[i].FileName == filename && buffer[i].BlockOffset == ofs {
			break
		}
	}
	if i < MAXBLOCKNUM {
		return &buffer[i]
	} else {
		bid := GetFreeBlockId()
		if bid == EOF {
			return nil
		}
		f, err := os.OpenFile(common.DIR+filename, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			fmt.Print(err)
			return nil
		}
		if !ReadBlockFromDisk1(filename, ofs, bid) {
			return nil
		}
		defer f.Close()
		return &buffer[bid]
	}
}

//也是从磁盘读出数据块。bid表示数据块序号值。
func ReadBlockFromDisk1(filename string, ofs int, bid int) bool {
	flag := false
	data := make([]byte, BLOCKSIZE)
	file, err := os.Open(common.DIR + filename)
	fileinfo, _ := os.Stat(common.DIR + filename)
	length := fileinfo.Size()
	if err != nil {
		fmt.Print(err)
		return false
	}
	defer file.Close()
	if int64((ofs+1)*BLOCKSIZE) <= length {
		_, err := file.Seek(int64(ofs*BLOCKSIZE), 0)
		if err != nil {
			fmt.Print(err)
			return false
		}
		_, err1 := io.ReadFull(file, data)
		if err1 != nil {
			fmt.Print(err1)
			return false
		}
	} else {
		for i := 0; i < len(data); i++ {
			data[i] = 0
		}
	}
	flag = true
	if flag {
		buffer[bid].ResetModes()
		buffer[bid].BlockData = data
		buffer[bid].FileName = filename
		buffer[bid].BlockOffset = ofs
		buffer[bid].IsValid = true
	}
	return flag
}

//向磁盘写入数据块
func WriteBlockToDisk(bid int) bool {
	if !buffer[bid].IsDirty {
		buffer[bid].IsValid = false
		return true
	}
	file, err := os.OpenFile(common.DIR+buffer[bid].FileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Print(err)
		return false
	}
	defer file.Close()
	_, err1 := file.Seek(int64(buffer[bid].BlockOffset*BLOCKSIZE), 0)
	if err1 != nil {
		fmt.Print(err1)
		return false
	}
	_, err2 := file.Write(buffer[bid].BlockData)
	if err2 != nil {
		fmt.Print(err2)
		return false
	}
	buffer[bid].IsValid = false
	buffer[bid].IsDirty = false
	return true

}

//获取一个空闲的数据块的序号
func GetFreeBlockId() int {
	index := EOF
	var mincount int = 0x7FFFFFFF
	for i := 0; i < MAXBLOCKNUM; i++ {
		if !buffer[i].IsLocked && buffer[i].LRUCount < mincount {
			index = i
			mincount = buffer[i].LRUCount
		}
	}
	if index != EOF && buffer[index].IsDirty {
		WriteBlockToDisk(index)
	}
	return index
}
