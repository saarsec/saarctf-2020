#!/usr/bin/python2
# -*- coding: iso-8859-15 -*-
#
# Hey pwners,
# 
# I guess you hate this service already.
#
# If you think it's to easy, be fair and run this script to make your life even harder. 
#
# Happy pwning,
# alfink
#
# 
#
import sys
import os

v = []
files = set()

def loadfile(fn):
    global v
    global files
    files.add(fn)
    result = ""
    with open(fn) as f:
        for l in f.readlines():
            l2 = l.strip()
            if l2.startswith("include "):
                result += loadfile(l2[8:-1].strip())
            elif l2.startswith("set "):
                l3 = l2.split()
                name = l3[1].replace("{","").replace("}","")
                v.append((name, "$var"+str(len(v))+"x"))
                result += l2 + "\n"
                v.append(("${"+name[1:]+"}", "${var"+str(len(v)-1)+"x}"))
            elif l2 == "":
                continue
            else:
                result += l2 + "\n"
    return result

minified = loadfile("./nginx.conf")

obfuscated = ""
if len(sys.argv) == 2 and sys.argv[1] == "insane":
    for n, o in sorted(v,key=lambda x: -len(x[0].replace("{","").replace("}",""))):
        print "#",n, o
    for l in minified.split("\n"):
        for n, o in sorted(v,key=lambda x: -len(x[0].replace("{","").replace("}",""))):
            l = l.replace(n, o)
        obfuscated += l + "\n"
    out = obfuscated
elif len(sys.argv) == 2 and sys.argv[1] == "hard":
    out = minified
else:
    print("usage: python ./next_level.py [difficulty]")
    print("    difficulty in ['hard', 'insane']")
    exit(1)

for fn in files:
    os.remove(fn)

with open("./nginx.conf", "w") as f:
    f.write(out)


#
# The obfuscation for insane mode is still a bit broken, so you also have to fix things to keep the SLA up  ¯\_(ツ)_/¯
#
