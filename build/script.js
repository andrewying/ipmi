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

const path = require('path');
const program = require('commander');
require('colors');
const emoji = require('node-emoji');
const git = require('nodegit');
const assetCompiler = require('./assets');
const Server = require('./server');

const version = require('../package.json').version;
let commit;

let server;

process.on('SIGTERM', function () {
    console.log('Gracefully shutting down...'.green);

    if (server) server.kill();
    process.exit(0);
});

program
    .version(version)
    .option('-e, --env [env]', 'set build environment', 'development')
    .parse(process.argv);

console.log(` ${emoji.get('package')} Adsisto Build Script `.bgMagenta.green);

git.Repository.open(path.resolve(__dirname, '../'))
        .then(function (repo) {
            repo.getHeadCommit()
                .then(function (res) {
                    commit = res.id().tostrS();
                    console.log(`Commit: ${ commit }\n`.green);

                    assetCompiler(program.env, function () {
                        if (program.env !== 'production') {
                            server = Server.run();
                        } else {
                            Server.build({
                                version: version,
                                commit: commit
                            });

                            console.log(`${emoji.get('white_check_mark')} Successfully built binaries.`.green);
                            process.exit(0);
                        }
                    });
                })
                .catch(function (error) {
                    console.error(
                        `${emoji.get('exclamation')} Unable to get current commit`.bold.red,
                        '\n',
                        error.toString()
                    );
                    process.exit(1);
                });
        })
        .catch(function (error) {
            console.error(
                `${emoji.get('exclamation')} Local directory is not a valid Git repository`.bold.red,
                '\n',
                error.toString()
            );
            process.exit(1);
        });
