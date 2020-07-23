if ((!file.startsWith("/home/saarlendar/config/") && !file.startsWith("/home/saarlendar/module/")) || file.indexOf("..") !== -1) {
    print("Access denied!");
}
else {

print(readfile(file));

}