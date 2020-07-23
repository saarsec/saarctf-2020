user = {"name": undefined, "password": undefined};
messages = undefined;
function get_messages(){
    if (messages === undefined && readfile !== undefined)
        messages = readfile("/home/saarlendar/users/"+user["name"]+"_messages");
    return messages;
}
events = undefined;
function get_events(){
    if (events === undefined && readfile !== undefined)
        events = readfile("/home/saarlendar/events/"+user["name"]);
    return events;
}
function login(){
    fetch("http://localhost:1337/login?username="+user["name"]+"&password="+user["password"])
}