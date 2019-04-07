import React from "react";
import { withRouter } from "react-router-dom";

import Logo from "../../images/logo.png";

class Header extends React.Component {
    render() {
        return (
            <header>
                <img className="header__logo" src={ Logo } alt="Cam FM"/>
                <h1 className="header__title">IPMI <small>v1.0.0</small></h1>
                <div className="header__status">
                    Hostname: <br />
                    ISO file: <br />
                    Time:
                </div>
                <ul className="header__actions">
                    <li><a href="" className="btn" role="button">Logout</a></li>
                </ul>
            </header>
        )
    }
}

export default withRouter(Header);