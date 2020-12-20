// Based on:
// https://webpack.js.org/guides/development/

const path = require('path')
const express = require('express')
const webpack = require('webpack')
const webpackDevMiddleware = require('webpack-dev-middleware')
const { chromium } = require('playwright')

const app = express()
const config = require('./webpack.config.js')
const compiler = webpack(config)

const bundleFile = path.resolve('../../wasm/bundlemain/main.wasm')
app.get('/main.wasm', (req, res) => {
  res.sendFile(bundleFile)
})

// Tell express to use the webpack-dev-middleware and use the webpack.config.js
// configuration file as a base.
app.use(
  webpackDevMiddleware(compiler, {
    publicPath: config.output.publicPath
  })
)

// Serve the files on port 3000.
app.listen(3000, function () {
  console.log('Example app listening on port 3000!\n')
  run()
})

async function run () {
  // const browser = await chromium.launch({ headless: false })
  const browser = await chromium.launch()
  const page = await browser.newPage()
  page.on('console', msg => {
    for (let i = 0; i < msg.args().length; ++i) {
      const line = `${msg.args()[i]}`
      console.log(`${i}: ${line}`)
      if (line.match(/Query Ask WSS:/)) {
        process.exit(0)
      }
    }
  })
  await page.goto('http://localhost:3000/')
}
