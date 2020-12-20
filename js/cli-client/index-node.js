require('isomorphic-fetch')
const delay = require('delay')
const { LotusRPC } = require('@filecoin-shipyard/lotus-client-rpc')
const { mainnet } = require('@filecoin-shipyard/lotus-client-schema')
const { BrowserProvider } = require('./browser-provider')
const { WasmProvider } = require('./wasm-provider')

global.WebSocket = require('websocket').w3cwebsocket

const wasmExec = require('./wasm_exec_modified.js')

async function run () {
  console.log('Starting WASM...')
  wasmExec('../../wasm/bundlemain/main.wasm') // async
  console.log('Sleeping...')
  await delay(3000)

  const wsUrl = 'wss://lotus.jimpick.com/spacerace_api/0/node/rpc/v0'
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

  const wasmQueryAskServiceProvider = new WasmProvider(
    global.connectQueryAskService,
    {
      environment: {
        requestsForLotusHandler
      }
    }
  )

  const queryAskClient = new LotusRPC(wasmQueryAskServiceProvider, {
    schema: mainnet.fullNode
  })

  console.log('Query Ask WSS')
  const result = await queryAskClient.clientQueryAsk(
    '12D3KooWEUS7VnaRrHF24GTWVGYtcEsmr3jsnNLcsEwPU7rDgjf5',
    'f063655'
  )
  console.log(`Query Ask WSS: ${JSON.stringify(result)}`)
}
run()

// "test": "node wasm_exec.js ../../wasm/bundlemain/main.wasm 12D3KooWEUS7VnaRrHF24GTWVGYtcEsmr3jsnNLcsEwPU7rDgjf5 f063655"
