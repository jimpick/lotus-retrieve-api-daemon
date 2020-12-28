const delay = require('delay')
const { LotusRPC } = require('@filecoin-shipyard/lotus-client-rpc')
const { mainnet } = require('@filecoin-shipyard/lotus-client-schema')
const { BrowserProvider } = require('./browser-provider')
const { WasmProvider } = require('./wasm-provider')
const { Lotus } = require('./browser-retrieval/shared/lotus-client/Lotus')

declare const Go: any

async function run () {
  try {
    // Initialize Lotus client and filecoin-signing-tools from browser-retrieval
    console.log('Starting Lotus client')
    const lotus = await Lotus.create()

    console.log('Starting WASM...')
    const go = new Go()
    try {
      const wasmResult = await WebAssembly.instantiateStreaming(
        fetch('main.wasm'),
        go.importObject
      )
      go.run(wasmResult.instance)
    } catch (e) {
      console.error('Error', e)
    }
    await delay(500) // FIXME: Get rid of this
    console.log('All systems go!')

    // console.log('Sleeping...')
    // await delay(3000)

    const wsUrl = 'wss://lotus.jimpick.com/calibration_api/0/node/rpc/v0'
    const browserProvider = new BrowserProvider(wsUrl, {
      token: async () => {
        const response = await fetch('/token')
        return await response.text()
      }
    })
    await browserProvider.connect()
    const requestsForLotusHandler = async (req, responseHandler) => {
      const request = JSON.parse(req)
      console.log('JSON-RPC request => Lotus', JSON.stringify(request))
      async function waitForResult () {
        try {
          const result = await browserProvider.sendWs(request)
          console.log('Jim result', JSON.stringify(result))
          responseHandler(JSON.stringify(result))
        } catch (e) {
          console.error('JSON-RPC error', e.message)
        }
      }
      waitForResult()
    }

    /*
    const wasmRetrievalServiceProvider = new WasmProvider(
      global.connectRetrievalService,
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
    console.log(`window.collectFileDeposit`, window.collectFileDeposit)
    const fileData = window.collectFileDeposit(fileDepositId)
    console.log(`fileData`, fileData)
    const blob = new Blob([fileData], {'type': 'image/jpeg'})
    const url = URL.createObjectURL(blob)
    const imgEl = document.createElement('img')
    imgEl.src = url
    imgEl.width = 500
    document.body.appendChild(imgEl)
    console.log('Retrieve WSS Success')
    */
  } catch (e) {
    console.error('Error', e.message)
  }
}

run()

// "test": "node wasm_exec.js ../../wasm/bundlemain/main.wasm 12D3KooWEUS7VnaRrHF24GTWVGYtcEsmr3jsnNLcsEwPU7rDgjf5 f063655"
