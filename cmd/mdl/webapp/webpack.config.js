/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-env node */
const HtmlWebPackPlugin = require('html-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const path = require('path');

const base = require('./webpack.config.base')

base.entry = {
	main: './src/index.tsx',
	headless: './src/headless.ts',
}
base.plugins.push(
	new HtmlWebPackPlugin({
		template: './src/index.html',
		filename: './index.html',
		chunks: ['main'],
	}),
	new HtmlWebPackPlugin({
		template: './src/index.html',
		filename: './headless.html',
		chunks: ['headless'],
	}),
	new CleanWebpackPlugin({
		protectWebpackAssets: false,
		cleanAfterEveryBuildPatterns: ['*.LICENSE.txt'],
	})
)
base.output = {
	path: path.resolve(__dirname, 'dist/'),
	publicPath: '/',
}
base.devtool = 'source-map'

module.exports = base
