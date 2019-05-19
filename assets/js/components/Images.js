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
import { withRouter } from 'react-router-dom';
import { useDropzone } from 'react-dropzone';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSpinner, faExclamationTriangle, faCheck } from '@fortawesome/free-solid-svg-icons';

function Images() {
  const [ errors, setErrors ] = useState([]);
  const [ images, setImages ] = useState([
    { name: 'test.iso', file: '123456.iso' },
  ]);
  const [ selectedImage, selectImage ] = useState('');
  const [ files, setFiles ] = useState({});

  const statusIcon = status => {
    switch (status) {
      case 0:
        return <FontAwesomeIcon icon={ faSpinner } spin/>;
      case 5:
        return <FontAwesomeIcon className="text-warning" icon={ faExclamationTriangle }/>;
      case 10:
        return <FontAwesomeIcon className="text-success" icon={ faCheck }/>;
    }
  };

  const filesList = () => {
    let array = Object.values(files);

    return <div className="py-2 text-left w-full">
      <p><strong>Files</strong></p>
      { array.map(entry => <div className="flex mb-1 w-full" key={ entry.file }>
        <span className="w-1/4">{ entry.file }</span>
        <progress className="mr-4 w-1/2" max="100" value={ entry.progress }/>
        <div>{ statusIcon(entry.status) }</div>
      </div>) }
    </div>;
  };

  const uploadFiles = acceptedFiles => {
    acceptedFiles.map(file => {
      for (let existingFile in files) {
        if (existingFile.file === file.name
          && existingFile.status >= 5) {
          return 2;
        }
      }

      let clone = Object.assign({}, files);
      clone[file.name] = {
        status: 0,
        file: file.name,
        uploadedFile: '',
        progress: 0,
      };
      setFiles(clone);

      let data = new FormData();
      data.append('file', file);

      let request = new XMLHttpRequest();
      request.onreadystatechange = () => {
        if (request.readyState === 4) {
          let response = request.response;

          let clone = Object.assign({}, files);
          if (response.code === 200) {
            clone[file.name].status = 10;
          } else {
            clone[file.name].status = 5;

            let oldErrors = errors.slice(0);
            oldErrors.push(response.error);
            setErrors(oldErrors);
          }
          setFiles(clone);
        }
      };
      request.onprogress = e => {
        let clone = Object.assign({}, files);
        clone[file.name].progress = e.loaded / e.total;
        setFiles(clone);
      };

      request.open('POST', '/api/images');
      request.responseType = 'json';
      request.send(data);
    });
  };

  const commitFiles = e => {
    e.preventDefault();

    let data = [];
    for (let file in files) {
      data.push({
        file: file.uploadedFile,
        name: file.file,
        commit: file.status === 10,
      });
    }

    fetch('/api/images', {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'same-origin',
      body: JSON.stringify({ files: data }),
    })
      .then(res => res.json())
      .then(res => {
        if (res.code === 200) {
          setErrors([]);
          setFiles({});
          setImages(res.images);
          return;
        }

        setErrors([ res.error ]);
      });
  };

  const loadImage = () => {
    const pattern = /^[0-9A-F]{64}.iso$/;
    if (!pattern.test(selectedImage)) {
      alert('Invalid image file name supplied.');
      return;
    }

    fetch('/api/images/' + selectedImage + '/load', {
      method: 'PATCH',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'same-origin',
    });
  };

  const {
    getRootProps,
    getInputProps,
    isDragActive,
  } = useDropzone({
    onDropAccepted: uploadFiles,
    multiple: false,
    accept: '.iso',
  });

  return (
    <div className="flex">
      { images.length !== 0 ? <div className="border border-solid flex-grow m-4 ml-0 px-4 py-5 rounded-lg images__list">
        <strong>Uploaded Images</strong>
        <div>
          { images.map(image =>
            <div className="mt-3 p-2 images__item" data-selected={ image.file === selectedImage
              ? '1' : '0' } data-value={ image.file } onClick={ () =>
              selectImage(image.file) }>
              { image.name }
            </div>,
          ) }
        </div>
        <button className="btn btn-primary mt-2">Load</button>
      </div> : '' }
      <div className="relative images__drop_container">
        { errors.length !== 0 ? <div className="alert alert-danger">
          <p><strong>The following errors occurred while uploading the image files:</strong>
          </p>
          <ul>
            { errors.map(error => <li>{ error }</li>) }
          </ul>
        </div> : '' }
        <div className={ isDragActive ? 'my-4 rounded-lg images images__active' : 'my-4 rounded-lg images' }>
          <div { ...getRootProps({ className: 'images__drop' }) }>
            <input { ...getInputProps({ name: 'images' }) } />
            { Object.keys(files).length !== 0 ? filesList(files) : <h3>Select Files</h3> }
          </div>
        </div>
        <button className="btn btn-primary" onClick={ commitFiles }>Upload</button>
      </div>
    </div>
  );
}

export default withRouter(Images);
