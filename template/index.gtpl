<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">
<html>
<head>
    <title>Test Module</title>
</head>
<body>
<div>
    <a href="#">Home</a>
</div>
<form action="/telnet" method="post">
    Host(URL):<input type="host" name="host">
    Port:<input type="port" name="port">
    <input type="submit" value="Telnet">
</form>
<form action="/ping" method="post">
    Host(URL):<input type="host" name="host">
    <input type="submit" value="Ping">
</form>
<form action="/nslookup" method="post">
    Host(URL):<input type="host" name="host">
    <input type="submit" value="LookupHost">
</form>
</body>
</html>
