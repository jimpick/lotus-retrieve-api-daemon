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

  /*
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
  //const result = await queryAskClient.clientQueryAsk(
  //  '12D3KooWHeqPF4yLunXpxaZyf8z4WbgK12YcYYeYLB5HEeQGaxAk',
  //  'f0105208'
  //)
  console.log(`Query Ask WSS: ${JSON.stringify(result)}`)
  */
}
run()

// "test": "node wasm_exec.js ../../wasm/bundlemain/main.wasm 12D3KooWEUS7VnaRrHF24GTWVGYtcEsmr3jsnNLcsEwPU7rDgjf5 f063655"
