/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-env node */
const HtmlWebPackPlugin = require('html-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const path = require('path');

const base = require('./webpack.config.base')

base.entry = './src/index.tsx'
base.plugins.push(
	new HtmlWebPackPlugin({
		template: './src/index.html',
		filename: './index.html',
	}),
	new CopyWebpackPlugin([{
		from: '*',
		context: 'static/'
	}])
)
base.output = {
	path: path.resolve(__dirname, 'dist/'),
		publicPath: '/',
}

module.exports = base
