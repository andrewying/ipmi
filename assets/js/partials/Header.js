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
import { connect } from "react-redux";

import Logo from "../../images/logo.png";
import Clock from "react-live-clock";

class Header extends React.Component {
    version = require('../../../package.json').version;

    render() {
        return (
            <header className="header">
                <img className="header__logo" src={ Logo } alt="Adsisto"/>
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
