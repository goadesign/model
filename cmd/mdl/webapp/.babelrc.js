module.exports = api => {
	// This caches the Babel config by environment.
	api.cache.using(() => process.env.NODE_ENV);

	return {
		presets: [
			'@babel/preset-env',
			'@babel/preset-typescript',
			'@babel/preset-react'
		],
		plugins: [
			// Applies the react-refresh Babel plugin on development modes only
			api.env('development') && 'react-refresh/babel'
		].filter(Boolean)
	};
};
