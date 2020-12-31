// based on https://webpack.js.org/guides/development/

const fs = require('fs')
const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const AddAssetHtmlPlugin = require('add-asset-html-webpack-plugin')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')
const Dotenv = require('dotenv-webpack')

const { size: wasmSize } = fs.statSync('./tmp/main.wasm.gz')

module.exports = {
  mode: 'development',
  target: 'web',
  entry: {
    index: './index.ts'
  },
  devtool: 'inline-source-map',
  devServer: {
    contentBase: './dist'
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/,
      },
    ],
  },
  resolve: {
    extensions: [ '.tsx', '.ts', '.js' ],
    fallback: { "util": require.resolve("util/") }
  },
  plugins: [
    new webpack.DefinePlugin({ 'process.env.WASM_SIZE': wasmSize }),
    new CleanWebpackPlugin({ cleanStaleWebpackAssets: false }),
    new HtmlWebpackPlugin({
      title: 'Filecoin WASM Retrieval Demo'
    }),
    new AddAssetHtmlPlugin({ filepath: require.resolve('./wasm_exec.js') }),
    new Dotenv(),
  ],
  output: {
    filename: '[name].bundle.js',
    path: path.resolve(__dirname, 'dist'),
    publicPath: '/'
  },
  experiments: {
    syncWebAssembly: true,
  }
}
