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
import { connect } from "react-redux";

import Logo from "../../images/logo.png";
import Clock from "react-live-clock";

class Header extends React.Component {
    render() {
        return (
            <header className="header">
                <img className="header__logo" src={ Logo } alt="Cam FM"/>
                <h1 className="header__title">IPMI <small>v1.0.0</small></h1>
                <div className="header__status">
                    <strong>Hostname:</strong><br />
                    <strong>ISO file:</strong>&nbsp;
                    { this.props.iso !== undefined ? this.props.iso : 'Not attached' }<br />
                    <strong><Clock format={'HH:mm:ss D MMM YYYY'} ticking={ true } /></strong>
                </div>
                <ul className="header__actions">
                    <li><a href="auth/logout" className="btn" role="button">Logout</a></li>
                </ul>
            </header>
        )
    }
}

const mapStateToProps = state => {
    return {
        iso: state.iso
    }
};

export default connect(mapStateToProps)(Header);
