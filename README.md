# blog_echo

## Data Base 

The data base is created in the blog.go file using sqlite. The database has 3 columns. 
The id, the title of the blog entry, and the body of the entry. We are able to add and delete things from the database with two functions called add_entry() and delete_entry()

## Web Server

A local webserver is created using the echo framework, and it has 4 requests; GET, POST, DELETE, and PUT. GET should print all of the entries as JSON objects. POST takes the title and body from the request and put it into the database. DELETE will delete a post based on the specified id of the blog in the URL (for example: localhost:8080/delete/3 to delete a blog with an id of 3). PUT will take an existing post and update it. PUT will also take the specified id from the URL.