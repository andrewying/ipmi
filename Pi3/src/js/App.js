import "@babel/polyfill";

import React from "react";
import ReactDOM from "react-dom";
import { HashRouter as Router, Route, Link } from "react-router-dom";
import { createStore } from "redux";
import { Provider } from "react-redux";
import Console from "./components/Console";
import Header from "./components/Header";

class App extends React.Component {
    render() {
        return (
            <div className="container">
                <Header />
                <Router>
                    <div className="sidebar">
                        <ul className="sidebar__menu">
                            <li>
                                <a href="#">Home</a>
                            </li>
                        </ul>
                    </div>
                    <div className="content">
                        <Route path="/" component={ Console } />
                    </div>
                </Router>
            </div>
        )
    }
}

const store = createStore();
const MOUNT_NODE = document.getElementById("app");

ReactDOM.render(
    <Provider store={store}>
        <App />
    </Provider>,
    MOUNT_NODE,
);