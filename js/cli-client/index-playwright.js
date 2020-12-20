const delay = require('delay')
const { LotusRPC } = require('@filecoin-shipyard/lotus-client-rpc')
const { mainnet } = require('@filecoin-shipyard/lotus-client-schema')
const { BrowserProvider } = require('./browser-provider')
const { WasmProvider } = require('./wasm-provider')

async function run () {
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
  status.innerText = 'All systems good! JS and Go loaded.'
  console.log('All systems go!')

  // console.log('Sleeping...')
  // await delay(3000)

  const wsUrl = 'wss://lotus.jimpick.com/spacerace_api/1/node/rpc/v0'
  const browserProvider = new BrowserProvider(wsUrl)
  await browserProvider.connect()
  const requestsForLotusHandler = async (req, responseHandler) => {
    const request = JSON.parse(req)
    console.log('JSON-RPC request => Lotus', request)
    async function waitForResult () {
      const result = await browserProvider.sendWs(request)
      console.log('Jim result', result)
      responseHandler(JSON.stringify(result))
    }
    waitForResult()
  }

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

  const order = {
    "Root": {
      "/": "bafykbzaced3v6jdz436uh2shndde7nwmjlmp6riomr6ps3fbapvaqb6dqpi2o"
    },
    "Piece": null,
    "Size": 8388608,
    "Total": "16777216",
    "UnsealPrice": "0",
    "PaymentInterval": 104857600,
    "PaymentIntervalIncrease": 104857600,
    "Client": "t3qkztmkfopk63qsel2xk3ek4w22epn3jnnlubnwjha2sl7rjhiuduwx24xivmhtdz7st3zmteuemeefply55q",
    "Miner": "f07281",
    "MinerPeer": {
      "Address": "f07283",
      "ID": "12D3KooWAU1x4P8XGCWyQBAapXXoGyom4Bx5QHnH4zQSeTqzMyQP",
      "PieceCID": null
    }
  }
  const fileref =  {
    "Path": "/tmp/wiki.zip.aa.aa-" + Math.random() * 100000000,
    "IsCAR": false
  } 
  console.log('Retrieve WSS')
  const result = await retrieveClient.clientRetrieve(
    order,
    fileref
  )
  console.log(`Retrieve WSS: ${JSON.stringify(result)}`)
}
run()

// "test": "node wasm_exec.js ../../wasm/bundlemain/main.wasm 12D3KooWEUS7VnaRrHF24GTWVGYtcEsmr3jsnNLcsEwPU7rDgjf5 f063655"
