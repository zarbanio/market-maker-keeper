// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dex_trader

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// DexTraderMetaData contains all meta data concerning the DexTrader contract.
var DexTraderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_executor\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"name\":\"Trade\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"executor\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swapRouter\",\"outputs\":[{\"internalType\":\"contractISwapRouter\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"poolFee\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMinimum\",\"type\":\"uint256\"}],\"name\":\"trade\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_executor\",\"type\":\"address\"}],\"name\":\"updateExecutor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_router\",\"type\":\"address\"}],\"name\":\"updateRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DexTraderABI is the input ABI used to generate the binding from.
// Deprecated: Use DexTraderMetaData.ABI instead.
var DexTraderABI = DexTraderMetaData.ABI

// DexTrader is an auto generated Go binding around an Ethereum contract.
type DexTrader struct {
	DexTraderCaller     // Read-only binding to the contract
	DexTraderTransactor // Write-only binding to the contract
	DexTraderFilterer   // Log filterer for contract events
}

// DexTraderCaller is an auto generated read-only Go binding around an Ethereum contract.
type DexTraderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DexTraderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DexTraderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DexTraderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DexTraderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DexTraderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DexTraderSession struct {
	Contract     *DexTrader        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DexTraderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DexTraderCallerSession struct {
	Contract *DexTraderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// DexTraderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DexTraderTransactorSession struct {
	Contract     *DexTraderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// DexTraderRaw is an auto generated low-level Go binding around an Ethereum contract.
type DexTraderRaw struct {
	Contract *DexTrader // Generic contract binding to access the raw methods on
}

// DexTraderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DexTraderCallerRaw struct {
	Contract *DexTraderCaller // Generic read-only contract binding to access the raw methods on
}

// DexTraderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DexTraderTransactorRaw struct {
	Contract *DexTraderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDexTrader creates a new instance of DexTrader, bound to a specific deployed contract.
func NewDexTrader(address common.Address, backend bind.ContractBackend) (*DexTrader, error) {
	contract, err := bindDexTrader(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DexTrader{DexTraderCaller: DexTraderCaller{contract: contract}, DexTraderTransactor: DexTraderTransactor{contract: contract}, DexTraderFilterer: DexTraderFilterer{contract: contract}}, nil
}

// NewDexTraderCaller creates a new read-only instance of DexTrader, bound to a specific deployed contract.
func NewDexTraderCaller(address common.Address, caller bind.ContractCaller) (*DexTraderCaller, error) {
	contract, err := bindDexTrader(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DexTraderCaller{contract: contract}, nil
}

// NewDexTraderTransactor creates a new write-only instance of DexTrader, bound to a specific deployed contract.
func NewDexTraderTransactor(address common.Address, transactor bind.ContractTransactor) (*DexTraderTransactor, error) {
	contract, err := bindDexTrader(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DexTraderTransactor{contract: contract}, nil
}

// NewDexTraderFilterer creates a new log filterer instance of DexTrader, bound to a specific deployed contract.
func NewDexTraderFilterer(address common.Address, filterer bind.ContractFilterer) (*DexTraderFilterer, error) {
	contract, err := bindDexTrader(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DexTraderFilterer{contract: contract}, nil
}

// bindDexTrader binds a generic wrapper to an already deployed contract.
func bindDexTrader(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DexTraderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DexTrader *DexTraderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DexTrader.Contract.DexTraderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DexTrader *DexTraderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DexTrader.Contract.DexTraderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DexTrader *DexTraderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DexTrader.Contract.DexTraderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DexTrader *DexTraderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DexTrader.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DexTrader *DexTraderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DexTrader.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DexTrader *DexTraderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DexTrader.Contract.contract.Transact(opts, method, params...)
}

// Executor is a free data retrieval call binding the contract method 0xc34c08e5.
//
// Solidity: function executor() view returns(address)
func (_DexTrader *DexTraderCaller) Executor(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DexTrader.contract.Call(opts, &out, "executor")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Executor is a free data retrieval call binding the contract method 0xc34c08e5.
//
// Solidity: function executor() view returns(address)
func (_DexTrader *DexTraderSession) Executor() (common.Address, error) {
	return _DexTrader.Contract.Executor(&_DexTrader.CallOpts)
}

// Executor is a free data retrieval call binding the contract method 0xc34c08e5.
//
// Solidity: function executor() view returns(address)
func (_DexTrader *DexTraderCallerSession) Executor() (common.Address, error) {
	return _DexTrader.Contract.Executor(&_DexTrader.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DexTrader *DexTraderCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DexTrader.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DexTrader *DexTraderSession) Owner() (common.Address, error) {
	return _DexTrader.Contract.Owner(&_DexTrader.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DexTrader *DexTraderCallerSession) Owner() (common.Address, error) {
	return _DexTrader.Contract.Owner(&_DexTrader.CallOpts)
}

// SwapRouter is a free data retrieval call binding the contract method 0xc31c9c07.
//
// Solidity: function swapRouter() view returns(address)
func (_DexTrader *DexTraderCaller) SwapRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DexTrader.contract.Call(opts, &out, "swapRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SwapRouter is a free data retrieval call binding the contract method 0xc31c9c07.
//
// Solidity: function swapRouter() view returns(address)
func (_DexTrader *DexTraderSession) SwapRouter() (common.Address, error) {
	return _DexTrader.Contract.SwapRouter(&_DexTrader.CallOpts)
}

// SwapRouter is a free data retrieval call binding the contract method 0xc31c9c07.
//
// Solidity: function swapRouter() view returns(address)
func (_DexTrader *DexTraderCallerSession) SwapRouter() (common.Address, error) {
	return _DexTrader.Contract.SwapRouter(&_DexTrader.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0xe1f21c67.
//
// Solidity: function approve(address _token, address spender, uint256 amount) returns()
func (_DexTrader *DexTraderTransactor) Approve(opts *bind.TransactOpts, _token common.Address, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _DexTrader.contract.Transact(opts, "approve", _token, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0xe1f21c67.
//
// Solidity: function approve(address _token, address spender, uint256 amount) returns()
func (_DexTrader *DexTraderSession) Approve(_token common.Address, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _DexTrader.Contract.Approve(&_DexTrader.TransactOpts, _token, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0xe1f21c67.
//
// Solidity: function approve(address _token, address spender, uint256 amount) returns()
func (_DexTrader *DexTraderTransactorSession) Approve(_token common.Address, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _DexTrader.Contract.Approve(&_DexTrader.TransactOpts, _token, spender, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DexTrader *DexTraderTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DexTrader.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DexTrader *DexTraderSession) RenounceOwnership() (*types.Transaction, error) {
	return _DexTrader.Contract.RenounceOwnership(&_DexTrader.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DexTrader *DexTraderTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DexTrader.Contract.RenounceOwnership(&_DexTrader.TransactOpts)
}

// Trade is a paid mutator transaction binding the contract method 0x917d61cd.
//
// Solidity: function trade(address token0, address token1, uint24 poolFee, uint256 amountIn, uint256 amountOutMinimum) returns(uint256 amountOut)
func (_DexTrader *DexTraderTransactor) Trade(opts *bind.TransactOpts, token0 common.Address, token1 common.Address, poolFee *big.Int, amountIn *big.Int, amountOutMinimum *big.Int) (*types.Transaction, error) {
	return _DexTrader.contract.Transact(opts, "trade", token0, token1, poolFee, amountIn, amountOutMinimum)
}

// Trade is a paid mutator transaction binding the contract method 0x917d61cd.
//
// Solidity: function trade(address token0, address token1, uint24 poolFee, uint256 amountIn, uint256 amountOutMinimum) returns(uint256 amountOut)
func (_DexTrader *DexTraderSession) Trade(token0 common.Address, token1 common.Address, poolFee *big.Int, amountIn *big.Int, amountOutMinimum *big.Int) (*types.Transaction, error) {
	return _DexTrader.Contract.Trade(&_DexTrader.TransactOpts, token0, token1, poolFee, amountIn, amountOutMinimum)
}

// Trade is a paid mutator transaction binding the contract method 0x917d61cd.
//
// Solidity: function trade(address token0, address token1, uint24 poolFee, uint256 amountIn, uint256 amountOutMinimum) returns(uint256 amountOut)
func (_DexTrader *DexTraderTransactorSession) Trade(token0 common.Address, token1 common.Address, poolFee *big.Int, amountIn *big.Int, amountOutMinimum *big.Int) (*types.Transaction, error) {
	return _DexTrader.Contract.Trade(&_DexTrader.TransactOpts, token0, token1, poolFee, amountIn, amountOutMinimum)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DexTrader *DexTraderTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DexTrader.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DexTrader *DexTraderSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DexTrader.Contract.TransferOwnership(&_DexTrader.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DexTrader *DexTraderTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DexTrader.Contract.TransferOwnership(&_DexTrader.TransactOpts, newOwner)
}

// UpdateExecutor is a paid mutator transaction binding the contract method 0x74936c16.
//
// Solidity: function updateExecutor(address _executor) returns()
func (_DexTrader *DexTraderTransactor) UpdateExecutor(opts *bind.TransactOpts, _executor common.Address) (*types.Transaction, error) {
	return _DexTrader.contract.Transact(opts, "updateExecutor", _executor)
}

// UpdateExecutor is a paid mutator transaction binding the contract method 0x74936c16.
//
// Solidity: function updateExecutor(address _executor) returns()
func (_DexTrader *DexTraderSession) UpdateExecutor(_executor common.Address) (*types.Transaction, error) {
	return _DexTrader.Contract.UpdateExecutor(&_DexTrader.TransactOpts, _executor)
}

// UpdateExecutor is a paid mutator transaction binding the contract method 0x74936c16.
//
// Solidity: function updateExecutor(address _executor) returns()
func (_DexTrader *DexTraderTransactorSession) UpdateExecutor(_executor common.Address) (*types.Transaction, error) {
	return _DexTrader.Contract.UpdateExecutor(&_DexTrader.TransactOpts, _executor)
}

// UpdateRouter is a paid mutator transaction binding the contract method 0xc851cc32.
//
// Solidity: function updateRouter(address _router) returns()
func (_DexTrader *DexTraderTransactor) UpdateRouter(opts *bind.TransactOpts, _router common.Address) (*types.Transaction, error) {
	return _DexTrader.contract.Transact(opts, "updateRouter", _router)
}

// UpdateRouter is a paid mutator transaction binding the contract method 0xc851cc32.
//
// Solidity: function updateRouter(address _router) returns()
func (_DexTrader *DexTraderSession) UpdateRouter(_router common.Address) (*types.Transaction, error) {
	return _DexTrader.Contract.UpdateRouter(&_DexTrader.TransactOpts, _router)
}

// UpdateRouter is a paid mutator transaction binding the contract method 0xc851cc32.
//
// Solidity: function updateRouter(address _router) returns()
func (_DexTrader *DexTraderTransactorSession) UpdateRouter(_router common.Address) (*types.Transaction, error) {
	return _DexTrader.Contract.UpdateRouter(&_DexTrader.TransactOpts, _router)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address _token, uint256 amount, address recipient) returns()
func (_DexTrader *DexTraderTransactor) WithdrawToken(opts *bind.TransactOpts, _token common.Address, amount *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _DexTrader.contract.Transact(opts, "withdrawToken", _token, amount, recipient)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address _token, uint256 amount, address recipient) returns()
func (_DexTrader *DexTraderSession) WithdrawToken(_token common.Address, amount *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _DexTrader.Contract.WithdrawToken(&_DexTrader.TransactOpts, _token, amount, recipient)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x3ccdbb28.
//
// Solidity: function withdrawToken(address _token, uint256 amount, address recipient) returns()
func (_DexTrader *DexTraderTransactorSession) WithdrawToken(_token common.Address, amount *big.Int, recipient common.Address) (*types.Transaction, error) {
	return _DexTrader.Contract.WithdrawToken(&_DexTrader.TransactOpts, _token, amount, recipient)
}

// DexTraderOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DexTrader contract.
type DexTraderOwnershipTransferredIterator struct {
	Event *DexTraderOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DexTraderOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DexTraderOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DexTraderOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DexTraderOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DexTraderOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DexTraderOwnershipTransferred represents a OwnershipTransferred event raised by the DexTrader contract.
type DexTraderOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DexTrader *DexTraderFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DexTraderOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DexTrader.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DexTraderOwnershipTransferredIterator{contract: _DexTrader.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DexTrader *DexTraderFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DexTraderOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DexTrader.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DexTraderOwnershipTransferred)
				if err := _DexTrader.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DexTrader *DexTraderFilterer) ParseOwnershipTransferred(log types.Log) (*DexTraderOwnershipTransferred, error) {
	event := new(DexTraderOwnershipTransferred)
	if err := _DexTrader.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DexTraderTradeIterator is returned from FilterTrade and is used to iterate over the raw logs and unpacked data for Trade events raised by the DexTrader contract.
type DexTraderTradeIterator struct {
	Event *DexTraderTrade // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DexTraderTradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DexTraderTrade)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DexTraderTrade)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DexTraderTradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DexTraderTradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DexTraderTrade represents a Trade event raised by the DexTrader contract.
type DexTraderTrade struct {
	Token0    common.Address
	Token1    common.Address
	Fee       *big.Int
	AmountIn  *big.Int
	AmountOut *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTrade is a free log retrieval operation binding the contract event 0x2fea533d210d2ab12cebb50fa91293400f42c2121ae7b66b10c589514bd416c0.
//
// Solidity: event Trade(address indexed token0, address indexed token1, uint24 fee, uint256 amountIn, uint256 amountOut)
func (_DexTrader *DexTraderFilterer) FilterTrade(opts *bind.FilterOpts, token0 []common.Address, token1 []common.Address) (*DexTraderTradeIterator, error) {

	var token0Rule []interface{}
	for _, token0Item := range token0 {
		token0Rule = append(token0Rule, token0Item)
	}
	var token1Rule []interface{}
	for _, token1Item := range token1 {
		token1Rule = append(token1Rule, token1Item)
	}

	logs, sub, err := _DexTrader.contract.FilterLogs(opts, "Trade", token0Rule, token1Rule)
	if err != nil {
		return nil, err
	}
	return &DexTraderTradeIterator{contract: _DexTrader.contract, event: "Trade", logs: logs, sub: sub}, nil
}

// WatchTrade is a free log subscription operation binding the contract event 0x2fea533d210d2ab12cebb50fa91293400f42c2121ae7b66b10c589514bd416c0.
//
// Solidity: event Trade(address indexed token0, address indexed token1, uint24 fee, uint256 amountIn, uint256 amountOut)
func (_DexTrader *DexTraderFilterer) WatchTrade(opts *bind.WatchOpts, sink chan<- *DexTraderTrade, token0 []common.Address, token1 []common.Address) (event.Subscription, error) {

	var token0Rule []interface{}
	for _, token0Item := range token0 {
		token0Rule = append(token0Rule, token0Item)
	}
	var token1Rule []interface{}
	for _, token1Item := range token1 {
		token1Rule = append(token1Rule, token1Item)
	}

	logs, sub, err := _DexTrader.contract.WatchLogs(opts, "Trade", token0Rule, token1Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DexTraderTrade)
				if err := _DexTrader.contract.UnpackLog(event, "Trade", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTrade is a log parse operation binding the contract event 0x2fea533d210d2ab12cebb50fa91293400f42c2121ae7b66b10c589514bd416c0.
//
// Solidity: event Trade(address indexed token0, address indexed token1, uint24 fee, uint256 amountIn, uint256 amountOut)
func (_DexTrader *DexTraderFilterer) ParseTrade(log types.Log) (*DexTraderTrade, error) {
	event := new(DexTraderTrade)
	if err := _DexTrader.contract.UnpackLog(event, "Trade", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
