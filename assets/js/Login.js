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

import React from 'react';
import ReactDOM from 'react-dom';
import { CookiesProvider, withCookies } from 'react-cookie';
import keydown from 'react-keydown';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faLock } from '@fortawesome/free-solid-svg-icons';
import Logo from '../images/logo.png';

class Login extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      token: '',
      help: false,
      loginText: 'Login',
      loginDisabled: false,
      success: undefined,
      redirect: undefined,
    };

    this.launchHelp = this.launchHelp.bind(this);
    this.fieldOnChange = this.fieldOnChange.bind(this);
    this.authenticate = this.authenticate.bind(this);
  }

  launchHelp(e) {
    e.preventDefault();
    this.setState({
      help: true,
    });
  }

  fieldOnChange(e) {
    let field = e.target.id;
    let update = {};
    update[field] = e.target.value;
    this.setState(update);
  }

  @keydown('enter')
  authenticate(e) {
    e.preventDefault();
    this.setState({
      loginText: 'Logging in...',
    });

    console.debug('Authenticating with IPMI server...');

    let parent = this;

    fetch('/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'same-origin',
      body: JSON.stringify({
        token: this.state.token,
      }),
    })
      .then(res => res.json())
      .then(data => {
        if (data.code !== 200) {
          parent.setState({
            success: false,
            loginText: 'Login',
          });
          return;
        }

        parent.props.cookies.set(window.cookieName, data.token, {
          secure: true,
        });
        parent.setState({
          redirect: '/',
        });
      });
  }

  render() {
    if (this.state.redirect) {
      window.location.replace(this.state.redirect);
      return '';
    }

    return (
      <div className="container">
        <div className="fullwidth-content">
          <div className="login__container">
            <h2>Login to IPMI</h2>
            { this.state.success === false ? <div className="alert alert-danger">
              <p><strong>Identity token invalid</strong></p>
              <p>The identity token you have provided was invalid.
                Please check that you have followed the instructions
                and that the token you enter was not expired.</p>
              <p>If you continue to have issues logging in, please
                contact your system administrator.</p>
            </div> : <div className="alert alert-warning">
              <p><strong>Access to this system is restricted to authorised
                users only.</strong> Use by others will be in contravention
                of the <i>Computer Misuse Act 1990</i>.</p>
              <p>All users are informed, in accordance to the <i>
                Investigatory Powers (Interception by Businesses etc.
                for Monitoring and Record-keeping Purposes) Regulations
                2018</i>, that their communication may be intercepted
                as permitted by the <i>Investigatory Powers Act 2016
                </i>.</p>
            </div> }
            <form className="login__form">
              <div className="input-group">
                <span><FontAwesomeIcon icon={ faLock }/></span>
                <input type="password" id="token" name="token"
                       placeholder="Security Token" value={ this.state.token }
                       onChange={ this.fieldOnChange } maxLength={ 500 }/>
              </div>
              <a className="help" onClick={ this.launchHelp }>
                Need help generating the token? Click here.
              </a>
              <input type="submit" className="btn btn-primary mt-1 mb-1"
                     onClick={ this.authenticate }
                     disabled={ this.state.loginDisabled }
                     value={ this.state.loginText }/>
            </form>
          </div>
          <footer className="login__footer">
            <img className="login__logo" src={ Logo } alt="Cam FM"/>
          </footer>
        </div>
      </div>
    );
  }
}

const MOUNT_NODE = document.getElementById('app');
Login = withCookies(Login);

ReactDOM.render(
  <CookiesProvider>
    <Login/>
  </CookiesProvider>,
  MOUNT_NODE,
);
