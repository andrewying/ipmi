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

import "@babel/polyfill";

import React from "react";
import ReactDOM from "react-dom";
import { HashRouter as Router, Route } from "react-router-dom";
import { createStore } from "redux";
import { Provider } from "react-redux";
import { useCookies } from "react-cookie";

import Header from "./partials/Header";
import Sidebar from "./partials/Sidebar";
import About from "./components/About";
import Accounts from "./components/Accounts";
import Console from "./components/Console";
import CreateAccount from "./components/CreateAccout";
import Images from "./components/Images";
import app from "./reducers";

function App() {
    const [ cookie, setCookie, removeCookie ] = useCookies([window.cookieName]);
    // if (!cookie) {
    //    window.location.replace("/auth/login");
    //    return "";
    // }

    return (
        <div>
            <Header />
            <Router>
                <div className="container">
                    <Sidebar />
                    <main className="content">
                        <Route path="/" exact component={ Console } />
                        <Route path="/about" component={ About } />
                        <Route path="/accounts" exact component={ Accounts } />
                        <Route path="/accounts/create" component={ CreateAccount } />
                        <Route path="/console" component={ Console } />
                        <Route path="/images" component={ Images } />
                    </main>
                </div>
            </Router>
        </div>
    )
}

const store = createStore(app);
const MOUNT_NODE = document.getElementById("app");

ReactDOM.render(
    <Provider store={store}>
        <App />
    </Provider>,
    MOUNT_NODE,
);
