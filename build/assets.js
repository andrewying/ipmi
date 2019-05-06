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

const Webpack = require('webpack');
require('colors');

const config = require('../webpack.config');

let compiler;
let server;
let lastHash = null;

function run(env, cb) {
    process.env.NODE_ENV = env;
    const production = env === 'production';

    try {
        compiler = Webpack(config);
    } catch (err) {
        console.error(err.message);
        process.exit(1);
        return;
    }

    compiler.run(function (err, stats) {
        if (compiler.close) {
            compiler.close(err2 => {
                callback(err || err2, stats);
            });
        } else {
            callback(err, stats);
        }

        if (!production) {
            compiler.watch(config.watch, callback);
            console.log('[WEBPACK]'.bgBlue.white + ' Asset builder started...');
        }

        cb();
    });
}

function callback(err, stats) {
    if (!config.watch || err) {
        compiler.purgeInputFileSystem();
    }

    if (err) {
        lastHash = null;
        console.error('[WEBPACK]'.bgBlue.white + ` ${ err.stack || err }`);

        if (server) server.kill();
        process.exit(1);
    }

    if (stats.hash !== lastHash) {
        lastHash = stats.hash;
        if (stats.compilation && stats.compilation.errors.length !== 0) {
            const errors = stats.compilation.errors;
            if (errors[0].name === "EntryModuleNotFoundError") {
                console.error("\n" + "[WEBPACK]".bgBlue.white +
                    " Insufficient number of arguments or no entry found.".red);
            }
        }

        const statsString = stats.toString({
            _env: process.env.NODE_ENV,
            cached: false,
            cachedAssets: false,
            children: false,
            chunks: false,
            colors: true,
            exclude: ["node_modules", "bower_components", "components"]
        });
        if (statsString) console.log('[WEBPACK]'.bgBlue.white + ` ${ statsString }\n`);
    }
}

module.exports = run;
