
code-gen:
	abigen --abi=abis/IERC20/IERC20.json --pkg=IERC20 --out=abis/IERC20/IERC20.go
	abigen --abi=abis/dex_trader/dex_trader.json --pkg=dex_trader --out=abis/dex_trader/dex_trader.go
	abigen --abi=abis/uniswapv3_pool/uniswapv3_pool.json --pkg=uniswapv3_pool --out=abis/uniswapv3_pool/uniswapv3_pool.go
	abigen --abi=abis/uniswapv3_factory/uniswapv3_factory.json --pkg=uniswapv3_factory --out=abis/uniswapv3_factory/uniswapv3_factory.go
	abigen --abi=abis/uniswapv3_quoter/uniswapv3_quoter.json --pkg=uniswapv3_quoter --out=abis/uniswapv3_quoter/uniswapv3_quoter.go

clean:
	rm abis/IERC20/IERC20.go
	rm abis/dex_trader/dex_trader.go
	rm abis/uniswapv3_pool/uniswapv3_pool.go
	rm abis/uniswapv3_factory/uniswapv3_factory.go
	rm abis/uniswapv3_quoter/uniswapv3_quoter.go
