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
                <th>Shift Type</th>
                <th>Task Type</th>
                <th>Task</th>
                <th>Hours</th>
                <th>Minutes</th>
                
            </tr>
        </thead>
        <tbody>
            {{range .ShiftTasks}}
            <tr>
                <td>{{.Name}}</td>
                <td>{{.ShiftDate}}</td>
                <td>{{.ShiftType}}</td>
                <td>{{.TaskType}}</td>
                <td>{{.Task}}</td>
                <td>{{.Hours}}</td>
                <td>{{.Minutes}}</td>
                
            </tr>
            {{end}}
        </tbody>
    </table>
    <h2>Total Hours: {{.TotalHours}}</h2>
	<a href="/search">Back to search</a><br>
    <a href="/">Back to Add Shift Task</a>
	</body>
	</html>
	`
