package nameservice

// 在这里定义应用程序用户可以对那些状态进行查询。
import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the nameservice Querier
const (
	// 传入一个域名返回nameservice给定的解析值。类似于DNS查询。
	QueryResolve = "resolve"
	// 传入一个域名返回价格，解析值和域名的所有者。用于确定你想要购买名称的成本。
	QueryWhois = "whois"
	QueryNames = "names"
)

// 该函数充当查询此模块的子路由器
// 因为querier没有类似于Msg的接口，所以需要手动定义switch语句（它们无法从query.Route()函数中删除）
// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryResolve:
			return queryResolve(ctx, path[1:], req, keeper)
		case QueryWhois:
			return queryWhois(ctx, path[1:], req, keeper)
		case QueryNames:
			return queryNames(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

// nolint: unparam
func queryResolve(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	value := keeper.ResolveName(ctx, path[0])

	if value == "" {
		return []byte{}, sdk.ErrUnknownRequest("could not resolve name")
	}
	// 按照惯例，每个输出类型都应该是 JSON marshallable 和 stringable（实现 Golang fmt.Stringer 接口）。
	// 返回的字节应该是输出结果的JSON编码。
	// 因此，对于输出类型的解析，我们将解析字符串包装在一个名为 QueryResResolve 的结构中，
	// 该结构既是JSON marshallable 的又有.String（）方法。
	// 在type/querier.go中
	res, err := codec.MarshalJSONIndent(keeper.cdc, QueryResResolve{value})
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

// nolint: unparam
func queryWhois(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	whois := keeper.GetWhois(ctx, path[0])
	// 对于 Whois 的输出，正常的 Whois 结构已经是 JSON marshallable 的，
	// 但我们需要在其上添加.String（）方法。 ??
	res, err := codec.MarshalJSONIndent(keeper.cdc, whois)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryNames(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var namesList QueryResNames

	iterator := keeper.GetNamesIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		namesList = append(namesList, string(iterator.Key()))
	}
	//名称查询的输出也一样，[]字符串本身已经可 marshallable ，但我们需要在其上添加.String（）方法。
	// 在type/querier.go中
	res, err := codec.MarshalJSONIndent(keeper.cdc, namesList)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
