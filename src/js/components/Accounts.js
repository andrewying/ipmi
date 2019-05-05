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

import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { connect } from "react-redux";
import ReactTable from "react-table";

function Accounts(props) {
    const [ loading, setLoading ] = useState(true);
    const [ accounts, setAccounts ] = useState([]);
    const [ errors, setErrors ] = useState([]);

    const columns = [
        {
            Header: 'Email',
            accessor: 'accessor'
        }
    ];

    const loadAccounts = () => {
        fetch('/api/accounts', {
            credentials: "same-origin",
        })
            .then(res => res.json())
            .then(res => {
                if (res.code === 200) {
                    setAccounts(res.accounts);
                    setLoading(false);
                    return;
                }

                setLoading(false);
            });
    };

    useEffect(() => {
        loadAccounts();
    });

    return (
        <div>
            <h2>Accounts</h2>
            { errors.length !== 0 ? <div className="alert alert-danger">
                <p><strong>The following errors occurred while retrieving the list of accounts:</strong></p>
                <ul>
                    { errors.map(error => <li>{ error }</li>) }
                </ul>
            </div> : "" }
            <Link to="/accounts/create" className="btn btn-primary">Create</Link>
            { loading ? <p>Loading...</p> : <ReactTable columns={ columns } data={ accounts } /> }
        </div>
    )
}

const mapStateToProps = state => {
    return {
        email: state.email,
        accounts: state.accounts
    }
};

export default connect(mapStateToProps)(Accounts);
