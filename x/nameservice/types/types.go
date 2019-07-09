package types

//我们要做的第一件事是定义一个结构，包含域名所有元数据。
//习惯上将模块相关的代码放在 ./x/

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)


// Whois is a struct that contains all the metadata of a name
type Whois struct {
	//域名解析出的值。这是任意字符串，
	//但将来您可以修改它以要求它适合特定格式，
	//例如IP地址，DNS区域文件或区块链地址。
	Value string         `json:"value"`
	//该域名当前所有者的地址
	Owner sdk.AccAddress `json:"owner"`
	//你需要为购买域名支付的费用
	Price sdk.Coins      `json:"price"`
}

// Initial Starting Price for a name that was never previously owned
//如果名称尚未有所有者，我们希望使用 MinPrice 对其进行初始化。
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

// Returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

// implement fmt.Stringer
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s`, w.Owner, w.Value, w.Price))
}
