#!/usr/bin/env node

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

const program = require('commander');
const { exec } = require('child_process');
const colors = require('colors');

program
    .version(require('../package.json').version)
    .option('-e, --env [env]', 'set build environment', 'development')
    .parse(process.argv);

process.env.NODE_ENV = program.env;
const production = program.env === 'production';

const Webpack = require('webpack');
const config = require('../webpack.config');

let lastHash = null;
let compiler;
let server;
try {
    compiler = Webpack(config);
} catch (err) {
    if (err.name === "WebpackOptionsValidationError") {
        console.error(err.message);
        process.exit(1);
    }

    throw err;
}

process.on('SIGTERM', function () {
    console.log('Gracefully shutting down...'.green);

    if (server) server.kill();
    process.exit(0);
});

function compilerCallback(err, stats) {
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

function runWatcher() {
    const path = require("path");
    const opn = require('opn');

    compiler.watch(config.watch, compilerCallback);
    console.log('[WEBPACK]'.bgBlue.white + ' Asset builder started...');

    server = exec('go run .',{
        cwd: path.resolve(__dirname, '../')
    }, (error, stdout, stderr) => {
        if (error) {
            console.error('[SERVER]'.bgGreen + ` ${error}`.red);
            process.exit(1);
        }

        console.log('[SERVER]'.bgGreen + ` ${stdout}`);
        console.log('[SERVER]'.bgGreen + ` ${stderr}`.red);
    });
    console.log('[SERVER]'.bgGreen + ' Web server started');
    setTimeout(() => opn('http://localhost:8080'), 5000);

    server.on('close', (code) => {
        console.log('[SERVER]'.bgGreen + ` Exited with code ${code}`.bold);
        process.exit(1);
    });
}

compiler.run((err, stats) => {
    if (compiler.close) {
        compiler.close(err2 => {
            compilerCallback(err || err2, stats);
        });
    } else {
        compilerCallback(err, stats);
    }

    if (!production) runWatcher();
});