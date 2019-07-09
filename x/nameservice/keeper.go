package nameservice

// 处理同存储的交互，引用其他的keeper进行跨模块的交互，
// 并包含模块的大部分核心功能。
import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"	//types包含了整个SDK常用的类型。
)

// Keeper maintains the link to data storage and exposes getter/setter methods
// for the various parts of the state machine
type Keeper struct {
	// bank模块控制账户和转账
	// 这是bank模块的Keeper引用。
	// 包括它来允许该模块中的代码调用bank模块的函数。SDK使用对象能力来访问应用程序状态的各个部分。
	// 这是为了允许开发人员采用小权限准入原则，限制错误或恶意模块的去影响其不需要访问的状态的能力。
	coinKeeper bank.Keeper
	// 通过它来访问一个持久化保存你的应用程序状态 sdk.KVStore
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	// 用于二进制编码/解码的线编解码器,提供负责Cosmos编码格式的工具 -- Amino
	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// Keeper的构造函数
// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Sets the entire Whois metadata struct for a name
// 添加一个函数来为指定域名设置解析字符串值
func (k Keeper) SetWhois(ctx sdk.Context, name string, whois Whois) {
	if whois.Owner.Empty() {
		return
	}
	//这个函数使用sdk.Context。该对象持有访问像blockHeight和chainID这样重要部分状态的函数。
	store := ctx.KVStore(k.storeKey)
	//.Set([]byte,[]byte)向存储中插入<name, value>键值对。
	// 由于存储只接受[]byte,想要把string转化成[]byte再把它们作为参数传给Set方法。
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(whois))
}


// Gets the entire Whois metadata struct for a name
// 添加一个函数来解析域名（即查找域名对应的解析值）
func (k Keeper) GetWhois(ctx sdk.Context, name string) Whois {
	//首先使用StoreKey访问存储
	store := ctx.KVStore(k.storeKey)
	//如果一个域名尚未在存储中，它返回一个新的 Whois 信息，包含最低价格 MinPrice。
	if !store.Has([]byte(name)) {
		return NewWhois()
	}
	bz := store.Get([]byte(name))
	var whois Whois
	k.cdc.MustUnmarshalBinaryBare(bz, &whois)
	return whois
}


// ResolveName - returns the string that the name resolves to
//根据名称返回域名解析出的值
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	return k.GetWhois(ctx, name).Value
}

// SetName - sets the value string that a name resolves to
//设置已有name的whois值为新name
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	whois := k.GetWhois(ctx, name)
	whois.Value = value
	k.SetWhois(ctx, name, whois)
}

//判断地址是否已被使用
// HasOwner - returns whether or not the name already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	return !k.GetWhois(ctx, name).Owner.Empty()
}

//获取地址所有者
// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	return k.GetWhois(ctx, name).Owner
}

//设置地址所有者
// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois := k.GetWhois(ctx, name)
	whois.Owner = owner
	k.SetWhois(ctx, name, whois)
}

//获取价格
// GetPrice - gets the current price of a name
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	return k.GetWhois(ctx, name).Price
}

//设置价格
// SetPrice - sets the current price of a name
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois := k.GetWhois(ctx, name)
	whois.Price = price
	k.SetWhois(ctx, name, whois)
}

// 获得迭代器，用于遍历指定 store 中的所有 <Key, Value> 对。
// Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
