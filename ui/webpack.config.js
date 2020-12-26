const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");

// Helpers
function relativepath(filepath) {
	return path.resolve(__dirname, filepath)
}

module.exports = {
	entry: relativepath("src/main.js"),
	output: {
		filename: "bundle.js",
		path: relativepath("dist")
	},
	resolve: {
		alias: {
			"$base": relativepath("src"),
			"$app": relativepath("src/app"),
		}
	},
	module: {
		rules: [
			{
				test: /\.m?js$/,
				exclude: /node_modules/,
				use: {
					loader: "babel-loader",
					options: {
						presets: ["@babel/preset-env"]
					}
				}
			},
			{
				test: /\.css$/,
				use: [
				  'style-loader',
				  'css-loader'
				]
			}
		]
	},
	// development options here
	devtool: "source-map",
	devServer: {
		contentBase: relativepath("dist"),
		compress: false,
		port: 9000,
		historyApiFallback: true,
		proxy: {
			"/event": {
				target: process.env.PROXY_SERVER || "http://localhost:8080",
				pathRewrite: {
					"^/event" : ""
				}
			}
    	}
	},
	// plugins here
	plugins: [
		new HtmlWebpackPlugin({
			template: relativepath("src/template.html")
		})
	]
};
