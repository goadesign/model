// CSS text is embedded in SVGs so exported files use the measured local font.
declare module "*.css?inline" {
	const css: string;
	export default css;
}

// The font's original copyright and license accompany every embedded copy.
declare module "*LICENSE?raw" {
	const license: string;
	export default license;
}
