//go:generate stringer -type ErrCode -linecomment
/*
基于插件stringer来生成代码，自动生成注释
使用go generate命令
*/

package code

type ErrCode int32

const (
	ERR_CODE_OK      ErrCode = iota + 1 //成功
	ERR_CODE_INVALID                    //参数错误
	ERR_CODE_TIMEOUT                    //超时
)
