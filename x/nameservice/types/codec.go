package types

// 在Amino中注册你的数据类型使得它们能够被编码/解码
import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// 创建的任何接口和实现接口的任何结构都需要在RegisterCodec函数中声明。
// 在此模块中，需要注册两个Msg的实现（SetName和BuyName），
// 但你的Whois查询返回的类型不需要
// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetName{}, "nameservice/SetName", nil)
	cdc.RegisterConcrete(MsgBuyName{}, "nameservice/BuyName", nil)
}
