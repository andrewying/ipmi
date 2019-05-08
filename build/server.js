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

const path = require('path');
const { exec } = require('child_process');
const opn = require('opn');

function run() {
    let server = exec('go run ./cmd/adsisto -dev',{
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
    setTimeout(() => opn('http://localhost:8080'), 10000);

    server.on('close', (code) => {
        console.log('[SERVER]'.bgGreen + ` Exited with code ${code}`.bold);
        process.exit(1);
    });

    return server;
}

function build(variables) {
    exec(`go build ./cmd/adsisto ${getArgs(variables)}`,{
        cwd: path.resolve(__dirname, '../')
    }, (error, stdout, stderr) => {
        if (error) {
            console.error('[SERVER]'.bgGreen + ` ${error}`.red);
            process.exit(1);
        }

        console.log('[SERVER]'.bgGreen + ` ${stdout}`);
        console.log('[SERVER]'.bgGreen + ` ${stderr}`.red);
    });

    console.log('[SERVER]'.bgGreen + ' Successfully built server binary');
}

function getArgs(arguments) {
    let array = [];
    for (let [ key, value ] of arguments) {
        array.push(`-X github.com/adsisto/adsisto/cmd/adsisto.${key}=${value}`);
    }

    return array.concat(' ');
}

module.exports = {
    run: run,
    build: build
};
