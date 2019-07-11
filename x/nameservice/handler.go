package nameservice

//Handler定义了在接收到一个特定Msg时，需要采取的操作（哪些存储需要更新，怎样更新及要满足什么条件）。
//在此模块中，你有两种类型的Msg，用户可以发送这些Msg来和应用程序状态进行交互：
// SetName和BuyName。它们各自同其Handler关联。
// test
import (
	"fmt"
	"github.com/jerryma0912/Cosmos-sdk-tutorial/x/nameservice/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler本质上是一个子路由，它将进入该模块的msg路由到正确的handler做处理。
// NewHandler returns a handler for "nameservice" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case types.MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// 定义处理MsgSetName消息的实际逻辑
// Handle a message to set name
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg MsgSetName) sdk.Result {
	//检查Msg的发送者是否就是域名的所有者(keeper.GetOwner)
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
		//如果不是，则抛出错误并返回给用户。
		return sdk.ErrUnauthorized("Incorrect Owner").Result() // If not, throw an error
	}
	//用Keeper里的函数来设置域名
	keeper.SetName(ctx, msg.Name, msg.Value) // If so, set the name to the value specified in the msg.
	return sdk.Result{}                      // return
}

// 定义BuyName的handler，该函数执行由msg触发的状态转换。
// 此时msg已运行其ValidateBasic函数，因此已进行了一些输入验证。
// 但是，ValidateBasic无法查询应用程序状态。
// 应在handler中执行依赖于网络状态（例如帐户余额）的验证逻辑。
// Handle a message to buy name
func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg MsgBuyName) sdk.Result {
	// 首先确保出价高于当前价格。然后，检查域名是否已有所有者。如果有，之前的所有者将会收到Buyer的钱。
	if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) { // Checks if the the bid price is greater than the price paid by the current owner
		return sdk.ErrInsufficientCoins("Bid not high enough").Result() // If not, throw an error
	}
	// 如果没有所有者，你的nameservice模块会把Buyer的资金“燃烧”（即发送到不可恢复的地址）。
	if keeper.HasOwner(ctx, msg.Name) {
		err := keeper.coinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	} else {
		_, err := keeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
		}
	}
	// 使用之前在Keeper上定义的 getter 和 setter，handler 将买方设置为新所有者，并将新价格设置为当前出价。
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return sdk.Result{}
}
