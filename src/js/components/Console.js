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

import React from "react";
import keydown, {ALL_KEYS} from "react-keydown";
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

        this.refWebSocket.sendMessage({
            key: e.key,
            shift: e.shiftKey,
            ctrl: e.ctrlKey,
            alt: e.altKey,
            meta: e.metaKey
        });
    }

    render() {
        return (
            <div>
                <Websocket url='ws://localhost:8888/live' onMessage={ (m) => {} }
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
                    <div className="console__overlay">
                        <h2>Connecting to Remote Machine</h2>
                        <h2 className="console__loading_container">
                            <span className="console__loading">.</span>&nbsp;
                            <span className="console__loading">.</span>&nbsp;
                            <span className="console__loading">.</span>&nbsp;
                        </h2>
                    </div>
                </div>
            </div>
        )
    }
}

export default Console;
