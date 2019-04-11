/*
 * Copyright (c) Andrew Ying 2019.
 *
 * This file is part of the Intelligent Platform Management Interface (IPMI) software.
 * IPMI is licensed under the API Copyleft License. A copy of the license is available
 * at LICENSE.md.
 *
 * As far as the law allows, this software comes as is, without any warranty or
 * condition, and no contributor will be liable to anyone for any damages related
 * to this software or this license, under any kind of legal claim.
 */

const path = require("path");
const webpack = require("webpack");
const ManifestPlugin = require("webpack-manifest-plugin");

const ExtractCssChunks = require("extract-css-chunks-webpack-plugin");
const UglifyJsPlugin = require("uglifyjs-webpack-plugin");

const production = process.env.NODE_ENV === "production";

let config = {
    mode: production ? "production" : "development",
    entry: {
        app: path.join(__dirname, "src", "index.js"),
        login: path.join(__dirname, "src", "login.js"),
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
                test: /\.(jpe?g|png|gif)$/,
                use: {
                    loader: "file-loader",
                    options: {
                        name: 'images/[name].[ext]',
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
            fetch: "imports-loader?this=>global!exports-loader?global.fetch!whatwg-fetch"
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
    config.watch = true;
    config.devtool = "source-map";
}

module.exports = config;