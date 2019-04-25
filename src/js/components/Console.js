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

import React from "react";
import keydown, { ALL_KEYS } from "react-keydown";
import Websocket from "react-websocket";

class Console extends React.Component {
    prohibitedKeys = [
        "Copy",
        "Cut",
        "Paste",
        "Undo",
        "Redo"
    ];

    constructor(props) {
        super(props);

        this.state = {
            connected: false,
            ws: false
        };
        this.handleOpen = this.handleOpen.bind(this);
        this.handleClose = this.handleClose.bind(this);
    }

    handleOpen() {
        this.setState({
            ws: true
        })
    }

    handleClose() {
        this.setState({
            ws: false
        })
    }

    @keydown(ALL_KEYS)
    handleKeys(e) {
        if (this.prohibitedKeys.contains(e.key)) {
            alert("Special keys are not currently supported!");
            return;
        }

        if (!this.state.ws) {
            alert("Connecting to remote Websocket server... Please try again.");
            return;
        }

        const message = {
            key: e.key,
            shift: e.shiftKey,
            ctrl: e.ctrlKey,
            alt: e.altKey,
            meta: e.metaKey
        };
        this.refWebSocket.sendMessage(JSON.stringify(message));
    }

    render() {
        let url = "ws://" + window.domain + "/api/keystrokes";

        return (
            <div>
                <Websocket url={ url } onMessage={ (m) => {} }
                           onOpen={ this.handleOpen } onClose={ this.handleClose }
                           reconnect={ true } debug={ false }
                           ref={ Websocket => {
                               this.refWebSocket = Websocket;
                           } }/>
                <p>
                    <strong>Status:</strong> Connecting...
                </p>
                <div className="console__container">
                    <video className="console" />
                    { this.state.connected ? "" : <div className="console__overlay">
                        <h2>Connecting to Remote Machine</h2>
                        <h2 className="console__loading_container">
                            <span className="console__loading">.</span>&nbsp;
                            <span className="console__loading">.</span>&nbsp;
                            <span className="console__loading">.</span>&nbsp;
                        </h2>
                    </div> }
                </div>
            </div>
        )
    }
}

export default Console;
