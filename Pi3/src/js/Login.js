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

import React from "react"
import ReactDOM from "react-dom";
import keydown from "react-keydown";

import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import Logo from "../images/logo.png";

class Login extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            token: "",
            help: false,
            loginText: "Login"
        };

        this.launchHelp = this.launchHelp.bind(this);
        this.fieldOnChange = this.fieldOnChange.bind(this);
        this.authenticate = this.authenticate.bind(this);
    }

    launchHelp(e) {
        e.preventDefault();
        this.setState({
            help: true
        });
    }

    fieldOnChange(e) {
        let field = e.target.id;
        let update = {};
        update[field] = e.target.value;
        this.setState(update);
    }

    @keydown("enter")
    authenticate(e) {
        e.preventDefault();
        console.debug("Authenticating with IPMI server...");
    }

    render() {
        return (
            <div className="container">
                <div className="fullwidth-content">
                    <div className="login__container">
                        <h2>Login to IPMI</h2>
                        <div className="alert-warning">
                            <p><strong>Access to this system is restricted to authorised
                            users only.</strong> Use by others will be in contravention
                            of the <i>Computer Misuse Act 1990</i>.</p>
                            <p>All users are informed, in accordance to the <i>
                            Investigatory Powers (Interception by Businesses etc.
                            for Monitoring and Record-keeping Purposes) Regulations
                            2018</i>, that their communication may be intercepted
                            as permitted by the <i>Investigatory Powers Act 2016
                            </i>.</p>
                        </div>
                        <form className="login__form">
                            <div className="input-group">
                                <span><FontAwesomeIcon icon={[ "fas", "lock" ]} /></span>
                                <input type="password" id="token" name="token"
                                       placeholder="Security Token"
                                       value={ this.state.token }
                                       onChange={ this.fieldOnChange }/>
                            </div>
                            <a className="help" onClick={ this.launchHelp }>
                                Need help generating the token? Click here.
                            </a>
                            <input type="submit" className="btn btn-primary"
                                   onClick={ this.authenticate } readOnly={ true }
                                   value={ this.state.loginText } />
                        </form>
                    </div>
                    <footer className="login__footer">
                        <img className="login__logo" src={ Logo } alt="Cam FM"/>
                    </footer>
                </div>
            </div>
        )
    }
}

const MOUNT_NODE = document.getElementById("app");

ReactDOM.render(
    <Login />,
    MOUNT_NODE,
);