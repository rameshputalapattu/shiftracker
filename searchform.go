package main

var searchForm = `<!DOCTYPE html>
<html>
    <head>
        <title>DAE Shift tracker search</title>
    </head>
    <body>
        <h1>Search</h1>
        <form method="post" action="/searchresults">
            <label for="name">Name:</label>
            <input type="text" name="name" required><br>

            <label for="shift_date">Shift Date:</label>
			
            <input type="date" name="shift_date" value="{{ .Today }}" required><br>

            <input type="submit" value="search">
			
        </form>
    </body>
</html>`
