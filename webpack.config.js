/*
 * Adsisto
 * Copyright (c) 2019 Andrew Ying
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of version 3 of the GNU General Public License as published by the
 * Free Software Foundation. In addition, this program is also subject to certain
 * additional terms available at <SUPPLEMENT.md>.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
 * A PARTICULAR PURPOSE.  See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program.  If not, see <https://www.gnu.org/licenses/>.
 */

const path = require("path");
const webpack = require("webpack");
const ManifestPlugin = require("webpack-manifest-plugin");

const ExtractCssChunks = require("extract-css-chunks-webpack-plugin");
const UglifyJsPlugin = require("uglifyjs-webpack-plugin");

const production = process.env.NODE_ENV === "production";

let config = {
    mode: production ? "production" : "development",
    context: path.resolve(__dirname, "src"),
    entry: {
        app: "./index.js",
        login: "./login.js",
    },
    output: {
        path: path.join(__dirname, "public"),
        filename: production ? "js/[name]-[chunkhash].js" : "js/[name].js",
        chunkFilename: production ? "js/[name]-[chunkhash].js" : "js/[name].js",
    },
    module: {
        rules: [
            {
                test: /\.jsx?$/,
                use: "babel-loader",
                exclude: /node_modules/
            },
            {
                test: /\.(jpe?g|png|gif|ttf|svg)$/,
                use: {
                    loader: "file-loader",
                    options: {
                        name: '[path][name].[ext]',
                    },
                }
            },
            {
                test: /\.(sa|sc|c)ss$/,
                use: [
                    {
                        loader: ExtractCssChunks.loader,
                        options: {
                            hot: !production
                        }
                    },
                    "css-loader",
                    "postcss-loader",
                    "sass-loader"
                ]
            },
        ]
    },
    resolve: {
        extensions: ['.js', '.jsx', '.css', '.scss'],
        alias: {
            "~": path.resolve(__dirname, "node_modules")
        }
    },
    plugins: [
        new webpack.ProvidePlugin({
            fetch: "exports-loader?self.fetch!whatwg-fetch/dist/fetch.umd"
        }),
        new ExtractCssChunks(
            {
                filename: production ? "css/[name]-[chunkhash].css" : "css/[name].css",
                chunkfilename: production ? "css/[name]-[id].css" : "css/[name].css"
            }
        ),
        new ManifestPlugin({
            writeToFileEmit: true,
        }),
    ],
    optimization: {
        minimize: production,
        splitChunks: {
            cacheGroups: {
                default: false,
                vendors: {
                    test: /[\\/]node_modules[\\/]/,
                    name: "vendor",
                    chunks: "initial",
                    enforce: true
                }
            }
        }
    }
};

if (production) {
    config.plugins.push(
        new webpack.DefinePlugin({
            "process.env": { NODE_ENV: JSON.stringify("production") }
        })
    );
    config.minimizer.push(new UglifyJsPlugin({
        cache: true,
        parallel: true,
        sourceMap: true
    }));
} else {
    config.plugins.push(
        new webpack.NamedModulesPlugin(),
    );
    config.devtool = "source-map";
}

module.exports = config;