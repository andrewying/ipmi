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

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faGithub } from '@fortawesome/free-brands-svg-icons';
import config from '../../../package.json';

export default function About() {
  return (
    <div>
      <h2>Adsisto v{ config.version }</h2>
      <a href="https://github.com/adsisto/adsisto" className="btn btn-outline" target="_blank">
        <FontAwesomeIcon icon={ faGithub } size="lg"/>&nbsp;&nbsp;GitHub
      </a>
      <p>Copyright &copy; { config.author.name } 2019.</p>
      <p>
        Adsisto is free software: you can redistribute it and/or modify
        it under the terms of version 3 of the <a
          href="https://github.com/adsisto/adsisto/blob/master/LICENSE.md"
          target="_blank">GNU General Public License</a> as published by the
        Free Software Foundation. In addition, this program is also subject
        to certain additional terms available <a
          href="https://github.com/adsisto/adsisto/blob/master/SUPPLEMENT.md" target="_blank">here</a>.
      </p>
      <p>
        This program is distributed in the hope that it will be useful,
        but <strong>WITHOUT ANY WARRANTY</strong>; without even the implied
        warranty of <strong>MERCHANTABILITY</strong> or <strong>FITNESS
        FOR A PARTICULAR PURPOSE</strong>. See the GNU General Public License
        for more details.
      </p>
    </div>
  );
}
