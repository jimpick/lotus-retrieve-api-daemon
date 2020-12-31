const delay = require('delay')
const pako = require('pako')
const { LotusRPC } = require('@filecoin-shipyard/lotus-client-rpc')
const { mainnet } = require('@filecoin-shipyard/lotus-client-schema')
const { BrowserProvider } = require('./browser-provider')
const { WasmProvider } = require('./wasm-provider')
const { Lotus } = require('./browser-retrieval/shared/lotus-client/Lotus')
import { appStore } from './browser-retrieval/shared/store/appStore'
const { toByteArray, fromByteArray } = require('base64-js')
const cbor = require('ipld-dag-cbor').util

// Global functions
declare const Go: any
declare const connectRetrievalService: any
declare const collectFileDeposit: any

async function download (url, defaultLength, status) {
  // https://dev.to/samthor/progress-indicator-with-fetch-1loo
  const response = await fetch(url)
  let length = response.headers.get('Content-Length')
  if (!length) {
    length = defaultLength
    // something was wrong with response, just give up
    // return await response.arrayBuffer()
  }
  const array = new Uint8Array(Number(length))
  let at = 0 // to index into the array
  const reader = response.body.getReader()
  for (;;) {
    const { done, value } = await reader.read()
    if (done) {
      break
    }
    status.textContent = `Fetching WASM bundle... ${at} of ${length} bytes`
    array.set(value, at)
    at += value.length
  }
  return array
}

async function run () {
  try {
    const h1 = document.createElement('h1')
    h1.textContent = 'Lotus WASM Retrieval Demo'
    document.body.appendChild(h1)

    const network = document.createElement('p')
    network.textContent = `Network: calibration`
    document.body.appendChild(network)

    const wallet = document.createElement('p')
    wallet.textContent = `Wallet: ${process.env.WALLET_1}`
    document.body.appendChild(wallet)

    const status = document.createElement('p')
    document.body.appendChild(status)

    // Initialize Lotus client and filecoin-signing-tools from browser-retrieval
    console.log('Starting Lotus client')
    const lotus = await Lotus.create()

    appStore.optionsStore.wallet = process.env.WALLET_1
    appStore.optionsStore.privateKey = process.env.WALLET_1_SECRET

    console.log('Starting WASM...', process.env.WASM_SIZE)
    const go = new Go()
    try {
      const url = 'main.wasm.gz' // the gzip-compressed wasm file

      const compressed = await download(url, process.env.WASM_SIZE, status)
      status.textContent = `Uncompressing WASM... (compressed size: ${compressed.byteLength} bytes)`
      await delay(100)
      const wasm = pako.ungzip(compressed)
      const size = +wasm.buffer.byteLength
    
      status.textContent = `Instantiating WASM... (uncompressed size: ${size} bytes)`
      const result = await WebAssembly.instantiate(wasm, go.importObject)
      go.run(result.instance)
    } catch (e) {
      console.error('Error', e)
    }
    await delay(500) // FIXME: Get rid of this
    status.innerText = 'All systems good! JS and Go loaded.'

    const lotusUrl = process.env.REACT_APP_LOTUS_ENDPOINT
    const browserProvider = new BrowserProvider(lotusUrl, {
      token: process.env.REACT_APP_LOTUS_TOKEN
    })
    await browserProvider.connect()
    const requestsForLotusHandler = makeRequestsForLotusHandler(
      browserProvider,
      lotus
    )

    const wasmRetrievalServiceProvider = new WasmProvider(
      connectRetrievalService,
      {
        environment: {
          requestsForLotusHandler
        }
      }
    )

    const retrieveClient = new LotusRPC(wasmRetrievalServiceProvider, {
      schema: mainnet.fullNode
    })

    const cid =
      'bafyaa4asfyfcmakvudsaeiaxi25hifvrnuwwqh5kl3hn7zysxdimchuxnllbrvzbbhn2my7boijaageaqbabelqkeyavliheaiqoizfeekow64ww6ibex2q62au7dstwjrtjpjwih4yx357q3s5pvqasaamj3zjdbihaqaqytxswgieaqbacbhpfem' // seal.jpg
    // const cid = 'bafykbzaced3v6jdz436uh2shndde7nwmjlmp6riomr6ps3fbapvaqb6dqpi2o' // wikipedia chunk
    const order = {
      Root: {
        '/': cid
      },
      Piece: null,
      Size: 8388608,
      Total: '16777216',
      UnsealPrice: '0',
      PaymentInterval: 104857600,
      PaymentIntervalIncrease: 104857600,
      Client:
        't3qkztmkfopk63qsel2xk3ek4w22epn3jnnlubnwjha2sl7rjhiuduwx24xivmhtdz7st3zmteuemeefply55q',
      Miner: 'f07281',
      MinerPeer: {
        Address: 'f07283',
        ID: '12D3KooWAU1x4P8XGCWyQBAapXXoGyom4Bx5QHnH4zQSeTqzMyQP',
        PieceCID: null
      }
    }
    const fileref = {
      Path: '/tmp/wiki.zip.aa.aa-' + Math.random() * 100000000,
      IsCAR: false
    }
    console.log('Retrieve WSS')
    const fileDepositId = await retrieveClient.clientRetrieve(order, fileref)
    console.log(`Retrieve WSS FileDepositID: ${JSON.stringify(fileDepositId)}`)
    console.log(`collectFileDeposit`, collectFileDeposit)
    const fileData = collectFileDeposit(fileDepositId)
    console.log(`fileData`, fileData)
    const blob = new Blob([fileData], { type: 'image/jpeg' })
    const url = URL.createObjectURL(blob)
    const imgEl = document.createElement('img')
    imgEl.src = url
    imgEl.width = 500
    document.body.appendChild(imgEl)
    console.log('Retrieve WSS Success')
  } catch (e) {
    console.error('Error', e.message)
  }
}

run()

// "test": "node wasm_exec.js ../../wasm/bundlemain/main.wasm 12D3KooWEUS7VnaRrHF24GTWVGYtcEsmr3jsnNLcsEwPU7rDgjf5 f063655"

function makeRequestsForLotusHandler (browserProvider, lotus) {
  const voucherNonces = {}
  const requestsForLotusHandler = async (req, responseHandler) => {
    const request = JSON.parse(req)
    console.log('JSON-RPC request => Lotus', JSON.stringify(request))
    if (request.method === 'Filecoin.PaychGet') {
      // Request: {"jsonrpc":"2.0","id":4,"method":"Filecoin.PaychGet",
      // "params": [
      //  "f3qkztmkfopk63qsel2xk3ek4w22epn3jnnlubnwjha2sl7rjhiuduwx24xivmhtdz7st3zmteuemeefply55q",
      //  "f07281",
      //  "16777216" ] }
      // Response: {"jsonrpc":"2.0","result":{
      //  "Channel":"t25dlrlbhotbryscbp5vgcijc4atlrigg6iqabu5a",
      //  "WaitSentinel":{
      //    "/":"bafy2bzacedl7gm6b3kaf7cd6c7e6l7xvbwfp73sz7y43lwq3kop2oqwkdrtcw"
      //   }},"id":4}
      const toAddr = request.params[0]
      const pchAmount = request.params[2]
      const { paymentChannel, msgCid } = await lotus.createPaymentChannel({
        toAddr,
        pchAmount
      })
      console.log('Zondax Payment channel', paymentChannel, msgCid)
      voucherNonces[paymentChannel.slice(1)] = 0
      responseHandler(
        JSON.stringify({
          jsonrpc: '2.0',
          result: {
            Channel: paymentChannel,
            WaitSentinel: msgCid
          },
          id: request.id
        })
      )
    } else if (request.method === 'Filecoin.PaychAllocateLane') {
      // Request: {"jsonrpc":"2.0","id":5,"method":"Filecoin.PaychAllocateLane",
      // "params": [
      //   "f25dlrlbhotbryscbp5vgcijc4atlrigg6iqabu5a"
      // ] }
      // Response: {"jsonrpc":"2.0","result":104,"id":5}
      await delay(0)
      responseHandler(
        JSON.stringify({
          jsonrpc: '2.0',
          result: 0, // Hard-coded
          id: request.id
        })
      )
    } else if (request.method === 'Filecoin.PaychVoucherCreate') {
      // Request: {"jsonrpc":"2.0","id":7,"method":"Filecoin.PaychVoucherCreate",
      // "params":[ "f25dlrlbhotbryscbp5vgcijc4atlrigg6iqabu5a", "3270170", 104] }
      // Response: {"jsonrpc":"2.0","result":{
      //  "Voucher":{
      //    "ChannelAddr": "t25dlrlbhotbryscbp5vgcijc4atlrigg6iqabu5a",
      //    "TimeLockMin": 0,
      //    "TimeLockMax": 0,
      //    "SecretPreimage": null,
      //    "Extra": null,
      //    "Lane": 104,
      //    "Nonce": 1,
      //    "Amount": "3270170",
      //    "MinSettleHeight": 0,
      //    "Merges": null,
      //    "Signature": {
      //      "Type": 2,
      //      "Data": "jXNHAkAVVghzKmheTA+DVGzA89ggzozL+3mhEaNm6iwV2uilx2HSBCCj04XNYg12Du23D+5vMX3ZYp0bgGb9erUYo96Sd6SiGGaCe4yZ4qBI/cwMn8En7HDZ7vO7liiv"
      //    }
      //  },
      // "Shortfall": "0"
      // }, "id":7}
      const paymentChannel = request.params[0]
      const amount = request.params[1]
      const nonce = voucherNonces[paymentChannel.slice(1)]++
      const signedVoucher = await lotus.createSignedVoucher(
        paymentChannel,
        amount,
        nonce
      )
      console.log('Jim signedVoucher', JSON.stringify(signedVoucher))
      let sigBytes = cbor.deserialize(toByteArray(signedVoucher))[10]
      console.log('Jim voucher sigBytes', JSON.stringify(sigBytes))
      responseHandler(
        JSON.stringify({
          jsonrpc: '2.0',
          result: {
            Voucher: {
              ChannelAddr: paymentChannel,
              TimeLockMin: 0,
              TimeLockMax: 0,
              SecretPreimage: null,
              Extra: null,
              Lane: 0, // hardcoded
              Nonce: nonce,
              Amount: amount,
              MinSettleHeight: 0,
              Merges: null,
              Signature: {
                Type: sigBytes[0],
                Data: fromByteArray(sigBytes.slice(1))
              }
            },
            Shortfall: '0'
          },
          id: request.id
        })
      )
    } else {
      async function callLotus () {
        try {
          const result = await browserProvider.sendHttp(request)
          console.log('Jim result', JSON.stringify(result))
          responseHandler(JSON.stringify(result))
        } catch (e) {
          console.error('JSON-RPC error', e.message)
        }
      }
      callLotus()
    }
  }
  return requestsForLotusHandler
}
