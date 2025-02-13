package call

import (
	"math/big"
	"testing"

	"github.com/depocket/multicall-go/utils"
	"github.com/stretchr/testify/assert"
)

var TestAddresses = map[Chain]string{
	Arbitrum:  "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
	Aurora:    "0x8BEc47865aDe3B172A928df8f990Bc7f2A3b9f79",
	Avalanche: "0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7",
	Bsc:       "0x7d99eda556388Ad7743A1B658b9C4FC67D7A9d74", // TotalSupply: 21000000
	Ethereum:  "0xB8c77482e45F1F44dE1745F52C74426C631bDD52",
	Fantom:    "0xe1146b9AC456fCbB60644c36Fd3F868A9072fc6E",
	Moonbeam:  "0xeFAeeE334F0Fd1712f9a8cc375f427D9Cdd40d73",
	Moonriver: "0xE3F5a90F9cb311505cd691a46596599aA1A0AD7D",
}

var TestAddressesV2 = map[Chain]string{
	Ethereum: "0xB8c77482e45F1F44dE1745F52C74426C631bDD52",
}

func TestContractBuilder_Default(t *testing.T) {
	caller := NewContractBuilder().AddMethod("function totalSupply()(uint256)")
	_, result, err := caller.
		AddCall("ts", TestAddresses[Ethereum], "totalSupply").
		Call(nil)
	if err != nil {
		assert.Failf(t, "Error calling %s contract", string(Ethereum))
	} else {
		assert.Equal(t, result["ts"][0].(*big.Int).Cmp(big.NewInt(0)), 1)
	}
}

func TestContractBuilder_BSC_ETH(t *testing.T) {
	caller := NewContractBuilder().
		WithChainConfig(DefaultChainConfigs[Bsc]).
		AddMethod("function totalSupply()(uint256)")
	_, result, err := caller.
		AddCall("ts", TestAddresses[Bsc], "totalSupply").
		Call(nil)
	if err != nil {
		assert.Failf(t, "Error calling %s contract", string(Bsc))
	} else {
		assert.Equal(t, utils.WeiToEther(result["ts"][0].(*big.Int)).String(), "21000000")
	}
}

func TestContractBuilder_Call(t *testing.T) {
	for chain, address := range TestAddresses {
		caller := NewContractBuilder().
			WithChainConfig(DefaultChainConfigs[chain]).
			AddMethod("function totalSupply()(uint256)")
		_, result, err := caller.AddCall("ts", address, "totalSupply").Call(nil)
		if err != nil {
			assert.Failf(t, "Error calling %s contract", string(chain))
		} else {
			assert.Equal(t, result["ts"][0].(*big.Int).Cmp(big.NewInt(0)), 1)
		}
	}
}

func TestContractBuilder_FlexibleCall(t *testing.T) {
	for chain, address := range TestAddressesV2 {
		caller := NewContractBuilder().
			WithChainConfig(DefaultChainConfigs[chain]).
			AddMethod("totalSupply()(uint256)")
		result, err := caller.AddCall("ts", address, "totalSupply").FlexibleCall(false)
		if err != nil {
			assert.Failf(t, "Error calling %s contract", string(chain))
		} else {
			assert.Equal(t, result["ts"].ReturnData[0].(*big.Int).Cmp(big.NewInt(0)), 1)
		}
	}
}

func TestContractBuilder_AddMethodTuple(t *testing.T) {
	caller := NewContractBuilder().
		WithChainConfig(DefaultChainConfigs[Fantom]).
		AddMethod("function claimableRewards(address)((address,uint256)[])").
		AddMethod("earnedBalances(address)(uint256,(uint256,uint256)[])").
		AddMethod("claimableRewardsTwo(address) ( (address,uint256)[], (address,address)[] )").
		AddMethod("claimableRewardsThree(address)((address,uint256)[],(address,address)[], uint256,(address,address,address,uint256)[]")

	assert.Equal(
		t,
		caller.contractAbi.Methods["claimableRewards"].String(),
		"function claimableRewards(address input0) view returns((address,uint256)[] output0)",
	)
	assert.Equal(
		t,
		caller.contractAbi.Methods["earnedBalances"].String(),
		"function earnedBalances(address input0) view returns(uint256 output0, (uint256,uint256)[] output1)",
	)
	assert.Equal(
		t,
		caller.contractAbi.Methods["claimableRewardsTwo"].String(),
		"function claimableRewardsTwo(address input0) view returns((address,uint256)[] output0, (address,address)[] output1)",
	)
	assert.Equal(
		t,
		caller.contractAbi.Methods["claimableRewardsThree"].String(),
		"function claimableRewardsThree(address input0) view returns((address,uint256)[] output0, (address,address)[] output1, uint256 output2, (address,address,address,uint256)[] output3)",
	)
}
