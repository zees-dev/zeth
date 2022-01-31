import { BigNumber } from "ethers"

export function expandToNDecimals(n: any, decimals: number): BigNumber {
  return BigNumber.from(n).mul(BigNumber.from(10).pow(decimals))
}

export function floatToBigNumber(n: number): BigNumber {
  const [int, dec] = n.toString().split(".")
  if (dec) {
    return expandToNDecimals(int + dec, int.length + dec.length)
  }
  return expandToNDecimals(int, int.length)
}
