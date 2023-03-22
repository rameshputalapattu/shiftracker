package main

var searchresults = `<!DOCTYPE html>
<html>
    <head>
        <title>DAE Shift tracker search Results</title>
    </head>
	<style>
        table {
            border-collapse: collapse;
            width: 100%;
        }
        th, td {
            text-align: left;
            padding: 8px;
            border: 1px solid black;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
    <body>

<table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Shift Date</th>
                <th>Task</th>
                <th>Hours</th>
                <th>Shift Type</th>
            </tr>
        </thead>
        <tbody>
            {{range .}}
            <tr>
                <td>{{.Name}}</td>
                <td>{{.ShiftDate}}</td>
                <td>{{.Task}}</td>
                <td>{{.Hours}}</td>
                <td>{{.ShiftType}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
	<a href="/search">Back to search</a>
	</body>
	</html>
	`
