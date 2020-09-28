/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-env node */
const HtmlWebPackPlugin = require('html-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');

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

module.exports = base
