/*
 * Adsisto
 * Copyright (c) 2019 Andrew Ying
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of version 3 of the GNU General Public License as published by the
 * Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
 * A PARTICULAR PURPOSE.  See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import "@babel/polyfill";

import React from "react";
import ReactDOM from "react-dom";
import { HashRouter as Router, Route, Link } from "react-router-dom";
import { createStore } from "redux";
import { Provider } from "react-redux";
import { CookiesProvider, withCookies } from "react-cookie";

import Console from "./components/Console";
import Header from "./components/Header";
import Images from "./components/Images";
import app from "./reducers";

class App extends React.Component {
    render() {
        const cookie = this.props.cookies.get(window.cookieName);
        if (!cookie) {
            window.location.replace("/auth/login");
            return "";
        }

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
App = withCookies(App);

ReactDOM.render(
    <CookiesProvider>
        <Provider store={store}>
            <App />
        </Provider>
    </CookiesProvider>,
    MOUNT_NODE,
);
