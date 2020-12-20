class WasmProvider {
  constructor (wasmConnect, options = {}) {
    if (wasmConnect === undefined) {
      throw new Error('WasmProvider wasmConnect param is undefined')
    }
    this.wasmConnect = wasmConnect
    this.id = 0
    this.inflight = new Map()
    this.cancelled = new Map()
    this.subscriptions = new Map()
    if (typeof options.token === 'function') {
      this.tokenCallback = options.token
    } else {
      this.token = options.token
      if (this.token && this.token !== '') {
        this.url += `?token=${this.token}`
      }
    }
    this.options = options
  }

  connect () {
    if (!this.connectPromise) {
      const getConnectPromise = () => {
        return new Promise((resolve, reject) => {
          console.log('Jim wasm-provider.js wasmConnect', this.options && this.options.environment)
          this.sendToWasm = this.wasmConnect(
            this.receive.bind(this),
            this.options && this.options.environment
          )
          console.log('Jim wasm-provider.js got sendToWasm', this.sendToWasm)
          resolve()
        })
      }
      this.connectPromise = getConnectPromise()
    }
    return this.connectPromise
  }

  send (request, schemaMethod) {
    const jsonRpcRequest = {
      jsonrpc: '2.0',
      id: this.id++,
      ...request
    }
    console.log('Jim send to wasm', jsonRpcRequest)
    const promise = new Promise((resolve, reject) => {
      console.log('Jim1', jsonRpcRequest.id, this.sendToWasm)
      this.inflight.set(jsonRpcRequest.id, (err, result) => {
        if (err) {
          reject(err)
        } else {
          resolve(result)
        }
      })
      console.log('Jim2', this.inflight)
      this.sendToWasm(JSON.stringify(jsonRpcRequest))
      console.log('Jim3')
      // FIXME: Add timeout
    })
    return promise
  }

  /*
  sendWs (jsonRpcRequest) {
    const promise = new Promise((resolve, reject) => {
      this.ws.send(JSON.stringify(jsonRpcRequest))
      // FIXME: Add timeout
      this.inflight.set(jsonRpcRequest.id, (err, result) => {
        if (err) {
          reject(err)
        } else {
          resolve(result)
        }
      })
    })
    return promise
  }
  */

  sendSubscription (request, schemaMethod, subscriptionCb) {
    let chanId = null
    const json = {
      jsonrpc: '2.0',
      id: this.id++,
      ...request
    }
    console.log('Jim send subscription to wasm', jsonRpcRequest)
    /*
    if (this.transport !== 'ws') {
      return [
        () => {},
        Promise.reject(
          new Error('Subscriptions only supported for WebSocket transport')
        )
      ]
    }
    */
    /*
    const promise = this.connect().then(() => {
      this.ws.send(JSON.stringify(json))
      // FIXME: Add timeout
      return new Promise((resolve, reject) => {
        this.inflight.set(json.id, (err, result) => {
          chanId = result
          // console.info(`New subscription ${json.id} using channel ${chanId}`)
          this.subscriptions.set(chanId, subscriptionCb)
          if (err) {
            reject(err)
          } else {
            resolve()
          }
        })
      })
    })
    return [cancel.bind(this), promise]
    async function cancel () {
      await promise
      this.inflight.delete(json.id)
      if (chanId !== null) {
        this.subscriptions.delete(chanId)
        await new Promise(resolve => {
          // FIXME: Add timeout
          this.cancelled.set(chanId, {
            cancelledAt: Date.now(),
            closeCb: resolve
          })
          if (!this.destroyed) {
            this.sendWs({
              jsonrpc: '2.0',
              method: 'xrpc.cancel',
              params: [json.id]
            })
          }
        })
        // console.info(`Subscription ${json.id} cancelled, channel ${chanId} closed.`)
      }
    }
    */
  }

  receive (response) {
    try {
      const { id, error, result, method, params } = JSON.parse(response)
      // FIXME: Check return code, errors
      if (method === 'xrpc.ch.val') {
        // FIXME: Check return code, errors
        const [chanId, data] = params
        const subscriptionCb = this.subscriptions.get(chanId)
        if (subscriptionCb) {
          subscriptionCb(data)
        } else {
          const { cancelledAt } = this.cancelled.get(chanId)
          if (cancelledAt) {
            if (Date.now() - cancelledAt > 2000) {
              console.warn(
                'Received stale response for cancelled subscription on channel',
                chanId
              )
            }
          } else {
            console.warn('Could not find subscription for channel', chanId)
          }
        }
      } else if (method === 'xrpc.ch.close') {
        // FIXME: Check return code, errors
        const [chanId] = params
        const { closeCb } = this.cancelled.get(chanId)
        if (!closeCb) {
          console.warn(`Channel ${chanId} was closed before being cancelled`)
        } else {
          // console.info(`Channel ${chanId} was closed, calling callback`)
          closeCb()
        }
      } else {
        const cb = this.inflight.get(id)
        if (cb) {
          this.inflight.delete(id)
          if (error) {
            // FIXME: Return error class with error.code
            return cb(new Error(error.message))
          }
          cb(null, result)
        } else {
          console.warn(`Couldn't find request for ${id}`)
        }
      }
    } catch (e) {
      console.error('RPC receive error', e)
    }
  }

  async importFile (body) {
    throw new Error('not implemented')
  }

  async destroy () {}
}

module.exports = { WasmProvider }
