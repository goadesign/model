/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-env node */
const HtmlWebPackPlugin = require('html-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const path = require('path');

const base = require('./webpack.config.base')

base.entry = './src/static/index.tsx'
base.plugins.push(
	new HtmlWebPackPlugin({
		template: './src/static/index.html',
		filename: './index.html',
	}),
	new CleanWebpackPlugin({
		protectWebpackAssets: false,
		cleanAfterEveryBuildPatterns: ['*.LICENSE.txt'],
	})
)
base.output = {
	path: path.resolve(__dirname, 'dist-static/'),
	publicPath: '.',
}

module.exports = base
