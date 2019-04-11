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

import "@babel/polyfill";

import React from "react";
import ReactDOM from "react-dom";
import { HashRouter as Router, Route, Link } from "react-router-dom";
import { createStore } from "redux";
import { Provider } from "react-redux";
import Console from "./components/Console";
import Header from "./components/Header";
import Images from "./components/Images";
import app from "./reducers";

class App extends React.Component {
    render() {
        return (
            <div className="container">
                <Header />
                <Router>
                    <div className="sidebar">
                        <ul className="sidebar__menu">
                            <li>
                                <Link to="/console">Console</Link>
                            </li>
                            <li>
                                <Link to="/images">Virtual Images</Link>
                            </li>
                        </ul>
                    </div>
                    <div className="content">
                        <Route path="/" exact component={ Console } />
                        <Route path="/console" component={ Console } />
                        <Route path="/images" component={ Images } />
                    </div>
                </Router>
            </div>
        )
    }
}

const store = createStore(app);
const MOUNT_NODE = document.getElementById("app");

ReactDOM.render(
    <Provider store={store}>
        <App />
    </Provider>,
    MOUNT_NODE,
);