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

import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import ReactTable from 'react-table';

function Accounts() {
  const [ loading, setLoading ] = useState(true);
  const [ accounts, setAccounts ] = useState([]);
  const [ error, setError ] = useState('');

  const columns = [
    {
      Header: 'Email',
      accessor: 'identity',
    }, {
      Header: 'Access Level',
      accessor: 'accessLevel',
    },
  ];

  const loadAccounts = () => {
    fetch('/api/keys', {
      credentials: 'same-origin',
    })
      .then(res => res.json())
      .then(res => {
        if (res.code === 200) {
          setAccounts(res.keys);
          setLoading(false);
          return;
        }

        setError(res.message);
        setLoading(false);
      });
  };

  useEffect(() => {
    loadAccounts();
  });

  return (
    <div>
      <h1 className="font-bold my-2 text-2xl">Accounts</h1>
      { error !== '' ? <div className="alert alert-danger">
        <p><strong>{ error }</strong></p>
      </div> : '' }
      <Link to="/accounts/new" className="btn btn-primary mb-4">Create</Link>
      { loading ? <p className="text-gray-800">Loading...</p> : <ReactTable columns={ columns } data={ accounts }/> }
    </div>
  );
}

const mapStateToProps = state => {
  return {
    email: state.email,
    accessLevel: state.accessLevel,
  };
};

export default connect(mapStateToProps)(Accounts);
