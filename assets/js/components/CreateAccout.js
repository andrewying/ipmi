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
import { Link, Redirect } from 'react-router-dom';

export default function CreateAccount(props) {
  const newAccount = typeof(props.identity) === 'undefined';

  const [ redirect, setRedirect ] = useState('');
  const [ error, setError ] = useState('');
  const [ email, setEmail ] = useState('');
  const [ publicKey, setPublicKey ] = useState('');
  const [ accessLevel, setAccessLevel ] = useState(0);

  const loadAccount = () => {
    setEmail(props.identity);

    fetch(`/api/keys?identity=${ props.identity }`, {
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'same-origin',
    })
      .then(res => res.json())
      .then(res => {
        if (res.code !== 200) {
          setRedirect('/accounts');
          return;
        }

        setPublicKey(res.key);
        setAccessLevel(res.accessLevel);
      });
  };

  const createAccount = (e) => {
    e.preventDefault();

    fetch('/api/keys', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'same-origin',
      body: JSON.stringify({
        email: email,
        key: publicKey,
        accessLevel: accessLevel,
      }),
    })
      .then(res => res.json())
      .then(res => {
        if (res.code !== 201) {
          setError(res.message);
          return;
        }

        setRedirect('/accounts');
      });
  };

  if (redirect !== '') {
    return <Redirect to={ redirect }/>;
  }

  return (
    <div>
      <h2>{ newAccount ? 'Create Account' : `Account for ${ email }` }</h2>
      { error !== '' ? <div className="alert alert-danger">
        <p><strong>{ error }</strong></p>
      </div> : '' }
      <form className="horizontal-form">
        <div className="form-row">
          <label>Email</label>
          <input type="email" id="email" name="email" value={ email }
                 onChange={ (e) => setEmail(e.target.value) }/>
        </div>
        <div className="form-row">
          <label>Public Key</label>
          <textarea id="publicKey" name={ publicKey } value={ publicKey }
                    onChange={ (e) => setPublicKey(e.target.value) }/>
        </div>
        <div className="form-row">
          <label>Access Level</label>
          <input type="number" min="0" step="1" id="accessLevel"
                 name="accessLevel" value={ accessLevel }
                 onChange={ (e) => setAccessLevel(e.target.value) }/>
        </div>
        <div className="form-row">
          <span></span>
          <div>
            <Link to="/accounts" className="btn btn-outline mt-1 mr-2 mb-1">Back</Link>
            <input type="submit" className="btn btn-primary mt-1 mb-1"
                   value="Create" onClick={ createAccount }/>
          </div>
        </div>
      </form>
    </div>
  );
}
