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
const git = require('nodegit');
require('colors');
const program = require('commander');
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

git.Repository.open(path.resolve(__dirname, '../'))
        .then(function (repo) {
            repo.getHeadCommit()
                .then(function (res) {
                    commit = res.id().tostrS();

                    assetCompiler(program.env, function () {
                        server = Server.run();
                    });
                })
                .catch(function (error) {
                    console.error('Unable to get current commit. '.bold.red);
                    console.error(error.toString());
                    process.exit(1);
                });
        })
        .catch(function (error) {
            console.error('Local directory is not a valid Git repository. '.bold.red);
            console.error(error.toString());
            process.exit(1);
        });
