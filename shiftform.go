package main

var formHTML = `<!DOCTYPE html>
<html>
    <head>
        <title>DAE Shift tracker</title>
    </head>
    <body>
        <h1>Add Shift Task</h1>
        <form method="post" action="/">
            <label for="name">Name:</label>
            <input type="text" name="name" required><br>

            <label for="shift_date">Shift Date:</label>
			
            <input type="date" name="shift_date" value="{{ .Today }}" required><br>

            <label for="shift_type">Shift Type:</label>
            <select name="shift_type" required>
                <option value="first">first</option>
                <option value="second">second</option>
                <option value="night">night</option>
                <option value="onshore">onshore</option>
            </select><br>

            <label for="task_type">Task Type:</label>
            <select name="task_type" required>
                <option value="monitoring">monitoring</option>
                <option value="ticket triaging">ticket triaging</option>
                <option value="meeting">meeting</option>
                <option value="failure handling">failure handling</option>
                <option value="automation development">automation development</option>
                <option value="others">others</option>
            </select><br>

            <label for="Task">Enter your Task:</label><br>
<textarea name="task" id="task" rows="10" cols="80" required></textarea><br>


            <label for="hours">Hours:</label>
            <input type="number" name="hours" min="0" max="12" required><br>

            <label for="minutes">Minutes:</label>
            <input type="number" name="minutes" min="0" max="59" required><br>

            <input type="submit" value="Add Shift Task">
        </form>
    </body>
</html>`
