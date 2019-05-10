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

package response

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

// JSON renders a JSON output.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		log.Printf("[ERROR] Unable to encode JSON output: %s\n", err)
	}
}

// View compiles and renders view template.
func View(w io.Writer, template *template.Template, data map[string]string) {
	err := template.Execute(w, data)
	if err != nil {
		log.Printf("[ERROR] Unable to render template: %s\n", err)
	}
}
