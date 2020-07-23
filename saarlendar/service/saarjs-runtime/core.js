include("/home/saarlendar/saarjs-runtime/base64.js");
include("/home/saarlendar/saarjs-runtime/helper.js");
include("/home/saarlendar/saarjs-runtime/api.js");

function drop_privileges(){
    include("/home/saarlendar/saarjs-runtime/drop_privileges.js");
    drop_privileges = undefined;
}