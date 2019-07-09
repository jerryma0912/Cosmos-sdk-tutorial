package types
//构建允许用户购买域名和设置解析值的Msg
import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// Msg触发状态转变。Msgs被包裹在客户端提交至网络的Txs中
// Cosmos SDK从Txs中打包和解包来自Msgs，
// 这就意味着，作为一个应用开发者，你只需要去定义Msgs。
// Msgs必须要满足下面的接口:
/*
// Transactions messages must fulfill the Msg
type Msg interface {
    // Return the message type.
    // Must be alphanumeric or empty.
    Type() string

    // Returns a human-readable string for the message, intended for utilization
    // within tags
    Route() string

    // ValidateBasic does a simple validation check that
    // doesn't require access to any other information.
    ValidateBasic() Error

    // Get the canonical byte representation of the Msg.
    GetSignBytes() []byte

    // Signers returns the addrs of signers that must sign.
    // CONTRACT: All signatures must be present to be valid.
    // CONTRACT: Returns addrs in some deterministic order.
    GetSigners() []AccAddress
}
*/

//SDK中Msg的命令约束是 Msg{.Action}。
// 要实现的第一个操作是SetName，因此命名为MsgSetName。
// 此Msg允许域名的所有者设置该域名的解析返回值。
// MsgSetName defines a SetName message
type MsgSetName struct {
	Name  string         `json:"name"`
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgSetName(name string, value string, owner sdk.AccAddress) MsgSetName {
	return MsgSetName{
		// 所要设置的域名
		Name:  name,
		//要设置的域名解析值
		Value: value,
		//域名的所有者
		Owner: owner,
	}
}

// 实现Msg接口
// SDK使用上述函数将Msg路由至合适的模块进行处理。它们还为用于索引的数据库标签添加了可读性的名称。
// Route should return the name of the module
func (msg MsgSetName) Route() string { return RouterKey }

// 实现Msg接口
// SDK使用上述函数将Msg路由至合适的模块进行处理。它们还为用于索引的数据库标签添加了可读性的名称。
// Type should return the action
func (msg MsgSetName) Type() string { return "set_name" }

// ValidateBasic用于对Msg的有效性进行一些基本的无状态检查。在此情形下，请检查没有属性为空。
// 请注意这里使用sdk.Error类型。 SDK提供了一组应用开发人员经常遇到的错误类型。
// ValidateBasic runs stateless checks on the message
func (msg MsgSetName) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("Name and/or Value cannot be empty")
	}
	return nil
}

// GetSignBytes定义了如何编码Msg以进行签名。
// 在大多数情形下，要编码成排好序的JSON。不应修改输出。
// GetSignBytes encodes the message for signing
func (msg MsgSetName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}


// GetSigners定义一个Tx上需要哪些人的签名才能使其有效。
// 在这种情形下，MsgSetName要求域名所有者在尝试重置域名解析值时要对该交易签名。
// GetSigners defines whose signature is required
func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgBuyName defines the BuyName message
type MsgBuyName struct {
	Name  string         `json:"name"`
	Bid   sdk.Coins      `json:"bid"`
	Buyer sdk.AccAddress `json:"buyer"`
}

// NewMsgBuyName is the constructor function for MsgBuyName
func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		Name:  name,
		Bid:   bid,
		Buyer: buyer,
	}
}

// Route should return the name of the module
func (msg MsgBuyName) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBuyName) Type() string { return "buy_name" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBuyName) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name cannot be empty")
	}
	if !msg.Bid.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBuyName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}
