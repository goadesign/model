/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-env node */
const ReactRefreshWebpackPlugin = require('@pmmmwh/react-refresh-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');

const isDevelopment = process.env.NODE_ENV !== 'production';

module.exports = {
	mode: isDevelopment ? 'development' : 'production',

	devtool: isDevelopment ? 'eval-source-map' : 'source-map',
	resolve: {
		extensions: ['.ts', '.tsx', '.js', '.jsx']
	},
	module: {
		rules: [
			{
				test: /\.(j|t)sx?$/,
				exclude: /node_modules/,
				use: {
					loader: 'babel-loader'
				}
			},
			{
				test: /\.css$/i,
				use: ['style-loader', 'css-loader']
			},
			{
				test: /\.html$/,
				use: [
					{
						loader: 'html-loader'
					}
				]
			},
			{ test: /\.svg$/, loader: 'svg-react-loader' },
			{
				test: /\.(png|jp(e*)g)$/,
				use: [{
					loader: 'url-loader',
					options: {
						limit: 8000, // Convert images < 8kb to base64 strings
						name: 'images/[hash]-[name].[ext]'
					}
				}]
			}
		]
	},
	plugins: [
		isDevelopment && new ReactRefreshWebpackPlugin()
	].filter(Boolean),
	optimization: {
		runtimeChunk: 'single',
		splitChunks: {
			chunks: 'all',
			cacheGroups: {
				vendor: {
					test: /[\/\\]node_modules[\/\\]/,
					name: 'vendors',
					chunks: 'all',
					priority: 10,
				},
				fontawesome: {
					test: /[\/\\]node_modules[\/\\]@fortawesome[\/\\]/,
					name: 'fontawesome',
					chunks: 'all',
					priority: 20,
				},
				elkjs: {
					test: /[\/\\]node_modules[\/\\]elkjs[\/\\]/,
					name: 'elkjs',
					chunks: 'all',
					priority: 20,
				},
				svgtools: {
					test: /[\/\\]node_modules[\/\\]svg-editor-tools[\/\\]/,
					name: 'svg-editor-tools',
					chunks: 'all',
					priority: 20,
				},
			},
		},
	},
	devServer: {
		historyApiFallback: true,
		hot: true,
	}
};
