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

import React, { useState } from 'react';
import ReactDOM from 'react-dom';
import { useCookies } from 'react-cookie';
import * as log from 'loglevel';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faLock } from '@fortawesome/free-solid-svg-icons';
import Logo from '../images/logo.png';

function Login() {
  const [ setCookie ] = useCookies([ window.cookieName ]);

  const [ clavis, setClavis ] = useState(true);
  const [ clavisError, setClavisError ] = useState('');
  const [ token, setToken ] = useState('');
  const [ help, setHelp ] = useState(false);
  const [ loginText, setLoginText ] = useState('Login');
  const [ loginDisabled, setLoginDisabled ] = useState(false);
  const [ success, setSuccess ] = useState(undefined);
  const [ redirect, setRedirect ] = useState(undefined);

  const launchHelp = e => {
    e.preventDefault();
    setHelp(true);
  };

  const tokenOnChange = e => {
    setToken(e.target.value);
  };

  const authenticate = e => {
    e.preventDefault();
    setLoginText('Logging in...');
    setLoginDisabled(true);

    if (clavis) {
      fetch('clavis://auth', {
        headers: {
          'Content-Type': 'application/json',
        },
      })
        .then(res => res.json())
        .then(data => {
          if (data.code !== 200) {
            setClavisError(data.message);
            return;
          }

          setToken(data.token);
        })
        .catch(error => {
          log.error('Error while obtaining authentication token from Clavis instance: '
            + error.toString());
          setClavisError(error.toString());
          setClavis(false);
        });
    }

    log.debug('Authenticating with server.');
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
        log.debug('Received authentication response from server.');

        if (data.code !== 200) {
          log.warn('Failed to authenticate with server: ' + data.message);
          setSuccess(false);
          setLoginText('Login');
          setLoginDisabled(false);
          return;
        }

        log.info('Sucessfully authenticated with server.');
        setCookie(window.cookieName, data.token, {
          secure: true,
        });
        setRedirect('/');
      });
  };

  if (redirect) {
    window.location.replace(redirect);
    return;
  }

  fetch('clavis://ping', {
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'same-origin',
  })
    .then(res => res.json())
    .then(data => {
      if (data.code !== 200) {
        log.warn('Unable to communicate with local Clavis application');
        setClavis(false);
        return;
      }

      setClavis(true);
    })
    .catch(error => {
      log.error('Error while communicating with local Clavis application: ' + error.toString());
      setClavis(false);
    });

  return (
    <div className="container">
      <div className="bg-primary-100 flex flex-grow flex-wrap justify-center min-h-screen px-4 fullwidth-content z-50">
        <div className="login-container">
          <h1 className="text-2xl my-2 font-bold">Login to IPMI</h1>
          { success === false ? <div className="alert alert-danger">
            <p><strong>Identity token invalid</strong></p>
            <p>The identity token you have provided was invalid. Please check that
            you have followed the instructions and that the token you enter was
            not expired.</p>
            <p>If you continue to have issues logging in, please contact your system
            administrator.</p>
          </div> : <div className="alert alert-warning">
            <p><strong>Access to this system is restricted to authorised users only.
            </strong> Use by others will be in contravention of the <i>Computer
            Misuse Act 1990</i>.</p>
            <p>All users are informed, in accordance to the <i>Investigatory Powers
            (Interception by Businesses etc. for Monitoring and Record-keeping Purposes)
            Regulations 2018</i>, that their communication may be intercepted as
            permitted by the <i>Investigatory Powers Act 2016</i>.</p>
          </div> }
          <form>
            { clavis ? '' : <div>
              <div className="input-group">
                <span><FontAwesomeIcon icon={ faLock }/></span>
                <input type="password" id="token" name="token"
                  placeholder="Security Token" value={ token }
                  onChange={ tokenOnChange } maxLength={ 500 }/>
              </div>
              <a className="help" onClick={ launchHelp }>
                Need help generating the token? Click here.
              </a>
            </div> }
            <input type="submit" className="btn btn-primary mt-4 mb-1 px-3 py-2 text-lg"
              onClick={ authenticate } disabled={ loginDisabled } value={ loginText }/>
          </form>
        </div>
        <footer className="flex justify-center mt-auto w-full">
          <img className="login-logo" src={ Logo } alt="Adsisto"/>
        </footer>
      </div>
    </div>
  );
}

const MOUNT_NODE = document.getElementById('app');
ReactDOM.render(<Login />, MOUNT_NODE);
