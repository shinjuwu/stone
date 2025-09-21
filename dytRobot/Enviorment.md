要執行該程式需先安裝Fyne
以下為安裝步驟
1.下載MinGW-w64 for gcc
(1)至下載頁面https://sourceforge.net/projects/mingw-w64/files/mingw-w64/
(2)下載x86_64-posix-seh
(3)解壓縮檔案
(4)把路徑bin添加至環境變數
(5)輸入gcc -v測試

2.安裝MSYS2 
(1)至下載頁面https://www.msys2.org/
(2)下載檔案"msys2-x86_64-20221028.exe"並安裝
   (https://mirrors.cloud.tencent.com/msys2/distrib/x86_64/msys2-x86_64-20221028.exe)
(3)安裝完成後開啟其終端機
(4)輸入pacman -Syu
(5)輸入pacman -S git mingw-w64-x86_64-toolchain
(6)有選項時就選all

3.撰寫程式碼測試
(1)建立main.go並添加以下程式碼
//==================================
package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}
//==================================
(2)編譯成功並產生一視窗就代表成功