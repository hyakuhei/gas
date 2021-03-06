// (c) Copyright 2016 Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package output

import (
	"encoding/json"
	"html/template"
	"io"

	gas "github.com/HewlettPackard/gas/core"
)

// The output format for reported issues
type ReportFormat int

const (
	ReportText ReportFormat = iota // Plain text format
	ReportJSON                     // Json format
	ReportCSV                      // CSV format
)

var text = `Results:
{{ range $index, $issue := .Issues }}
[{{ $issue.File }}:{{ $issue.Line }}] - {{ $issue.What }} (Confidence: {{ $issue.Confidence}}, Severity: {{ $issue.Severity }})
  > {{ $issue.Code }}

{{ end }}
Summary:
   Files: {{.Stats.NumFiles}}
   Lines: {{.Stats.NumLines}}
   Nosec: {{.Stats.NumNosec}}
  Issues: {{.Stats.NumFound}}

`

var csv = `{{ range $index, $issue := .Issues -}}
{{- $issue.File -}},
{{- $issue.Line -}},
{{- $issue.What -}},
{{- $issue.Severity -}},
{{- $issue.Confidence -}},
{{- printf "%q" $issue.Code }}
{{ end }}`

func CreateReport(w io.Writer, format string, data *gas.Analyzer) error {
	var err error
	switch format {
	case "json":
		err = reportJSON(w, data)
	case "csv":
		err = reportFromTemplate(w, csv, data)
	case "text":
		err = reportFromTemplate(w, text, data)
	default:
		err = reportFromTemplate(w, text, data)
	}
	return err
}

func reportJSON(w io.Writer, data *gas.Analyzer) error {
	raw, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		panic(err)
	}

	_, err = w.Write(raw)
	if err != nil {
		panic(err)
	}
	return err
}

func reportFromTemplate(w io.Writer, reportTemplate string, data *gas.Analyzer) error {
	t, e := template.New("gas").Parse(reportTemplate)
	if e != nil {
		return e
	}

	return t.Execute(w, data)
}
