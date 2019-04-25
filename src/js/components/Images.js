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

import React, { useState } from "react";
import { withRouter } from "react-router-dom";
import { useDropzone } from "react-dropzone";

import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner, faExclamationTriangle, faCheck } from "@fortawesome/free-solid-svg-icons";

function Images() {
    const [ errors, setErrors ] = useState([]);
    const [ files, setFiles ] = useState({});

    const statusIcon = status => {
        switch(status) {
            case 0:
                return <FontAwesomeIcon icon={ faSpinner } spin />;
            case 5:
                return <FontAwesomeIcon className="text-warning" icon={ faExclamationTriangle } />;
            case 10:
                return <FontAwesomeIcon className="text-success" icon={ faCheck } />;
        }
    };

    const filesList = () => {
        let array = Object.values(files);

        return <div className="images__files">
            <p><strong>Files</strong></p>
            { array.map(entry => <div className="images__file" key={ entry.file }>
                <span>{ entry.file }</span>
                <progress max="100" value={ entry.progress } />
                <div>{ statusIcon(entry.status) }</div>
            </div>) }
        </div>;
    };

    const uploadFiles = acceptedFiles => {
        acceptedFiles.map(file => {
            for (let existingFile in files) {
                if (existingFile.file === file.name) {
                    return 2;
                }
            }

            let clone = Object.assign({}, files);
            clone[file.name] = {
                status: 0,
                file: file.name,
                uploadedFile: "",
                progress: 0,
            };
            setFiles(clone);

            let data = new FormData();
            data.append("file", file);

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
                clone[file.name].progress = e.loaded/e.total;
                setFiles(clone);
            };

            request.open("POST", "/api/images");
            request.responseType = "json";
            request.send(data);
        });
    };

    const {
        getRootProps,
        getInputProps,
        isDragActive,
    } = useDropzone({
        onDropAccepted: uploadFiles,
        multiple: false,
        accept: ".iso",
    });

    return (
        <div>
            { errors.length !== 0 ? <div className="alert alert-danger">
                <p><strong>The following errors occurred while uploading the image files:</strong></p>
                <ul>
                    { errors.map(error => <li>{ error }</li>) }
                </ul>
            </div> : "" }
            <div className={ isDragActive ? "images images__active" : "images" }>
                <div { ...getRootProps({ className: "images__drop" }) }>
                    <input { ...getInputProps({ name: "images" }) } />
                    { Object.keys(files).length !== 0 ? filesList(files) : <h3>Select Files</h3> }
                </div>
            </div>
            <button className="btn btn-primary">Upload</button>
        </div>
    )
}

export default withRouter(Images);
