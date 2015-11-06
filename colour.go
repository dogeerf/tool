package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func CheckErr(err error, operating string) {
	//pc,file,line,ok = runtime.Caller(1)
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	if err != nil {
		log.Printf("@@@ ERROR: |%s| %s failed.\n", funcName, operating)
		log.Printf("  %s\n", err.Error())
		//panic(err)
		os.Exit(-1)
	}
	log.Printf("### OK: |%s| %s success.\n", funcName, operating)
}

type T struct {
	createTime time.Time
	startTime  time.Time
	endTime    time.Time
}

func NewT() *T {
	return &T{
		createTime: time.Now(),
	}
}
func (this *T) Show() {
	fmt.Printf("createTime:%v\n", this.createTime)
}

func init() {

}

func main() {
	fmt.Printf("==============================================================\n")
	fmt.Printf("\033[32;1m abc \033[0m")

	//字背景颜色范围:40----49
	//40:黑
	//41:深红
	//42:绿
	//43:黄色
	//44:蓝色
	//45:紫色
	//46:深绿
	//47:白色
	//字颜色:30-----------39
	//30:黑
	//31:红
	//32:绿
	//33:黄
	//34:蓝色
	//35:紫色
	//36:深绿
	//37:白色

	//  \033[0m 关闭所有属性
	//	\033[1m 设置高亮度
	//	\033[4m 下划线
	//	\033[5m 闪烁
	//	\033[7m 反显
	//	\033[8m 消隐
	//	\033[30m 至 \33[37m 设置前景色
	//	\033[40m 至 \33[47m 设置背景色
	//	\033[nA 光标上移n行
	//	\033[nB 光标下移n行
	//	\033[nC 光标右移n行
	//	\033[nD 光标左移n行
	//	\033[y;xH设置光标位置
	//	\033[2J 清屏
	//	\033[K 清除从光标到行尾的内容
	//	\033[s 保存光标位置
	//	\033[u 恢复光标位置
	//	\033[?25l 隐藏光标
	//	\033[?25h 显示光标<br>

	fmt.Printf("==============================================================\n")
}
