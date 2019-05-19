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
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';

import config from '../../../package.json';

function Sidebar(props) {
  const menuItems = [
    {
      text: 'Console',
      link: '/console',
      accessLevel: 0,
    }, {
      text: 'Virtual Images',
      link: '/images',
      accessLevel: 0,
    }, {
      text: 'Accounts',
      link: '/accounts',
      accessLevel: 100,
    }, {
      text: 'About',
      link: '/about',
      accessLevel: 0,
    },
  ];

  return (
    <nav className="sidebar">
      <ul className="sidebar-menu">
        { menuItems.map(item => props.accessLevel >= item.accessLevel ? <li>
          <Link to={ item.link }>{ item.text }</Link>
        </li> : '') }
      </ul>
      <div className="px-2 py-4">
        <span className="text-description">
          Adsisto v{ config.version }
        </span>
      </div>
    </nav>
  );
}

const mapStateToProps = state => {
  return {
    accessLevel: state.accessLevel,
  };
};

export default connect(mapStateToProps)(Sidebar);
