package saarschleife.template

import kotlinx.html.*
import saarschleife.models.VisitorAccount

fun HtmlBlockTag.parallax(src: String, alt: String? = null, content: (DIV.() -> Unit)? = null) {
	div("parallax-container valign-wrapper") {
		if (content != null) {
			div("container") {
				div("row center") {
					div("col s12 light") {
						content()
					}
				}
			}
		}
		div("parallax") {
			img(src = src, alt = alt)
		}
	}
}

fun UL.mainmenu(page: String, visitorAccount: VisitorAccount?) {
	if (visitorAccount != null) {
		li(classes = if (page == "visitor") "active" else null) {
			a("/visitor") {
				i("material-icons left") { +"people" }
				+"Visitors"
			}
		}
		li(classes = if (page == "secrets") "active" else null) {
			a("/secrets") {
				i("material-icons left") { +"work" }
				+"Secrets"
			}
		}

		li {
			a(classes = "link-logout") {
				i("material-icons left") { +"account_circle" }
				+"Logout"
			}
		}
	} else {
		li {
			a("/#login") {
				i("material-icons left") { +"account_circle" }
				+"Login"
			}
		}
	}
}

fun HTML.layout(page: String, visitor: VisitorAccount?, onload: String = "Endpoints.install()", layoutbody: DIV.() -> Unit) {
	head {
		title { +"saarschleife.net" }
		link("https://fonts.googleapis.com/icon?family=Material+Icons", "stylesheet")
		link("/static/css/materialize.css", "stylesheet", "text/css")
		unsafe { +"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"/>" }
		style {
			unsafe {
				+"#main-content { margin-bottom: 120px }"
				+"dt { font-weight: bold }"
				+".parallax{ margin-left: -1px }" // fix bug in css library
				+".nofloat{ float: none !important; }"
				+".btn-large.top-fab{ top: -28px; bottom: auto; }"
			}
		}
	}
	body {
		attributes["onload"] = "frontend.$onload"
		attributes["data-my-username"] = visitor?.username ?: ""

		nav(classes = "light-blue") {
			div(classes = "nav-wrapper container") {
				a("/", classes = "brand-logo") { +"saarschleife.net" }
				ul(classes = "right hide-on-med-and-down") {
					mainmenu(page, visitor)
				}
				unsafe { +"<a href=\"#\" data-target=\"nav-mobile\" class=\"sidenav-trigger\"><i class=\"material-icons\">menu</i></a>" }
			}
		}
		ul("sidenav") {
			id = "nav-mobile"
			mainmenu(page, visitor)
		}

		// Parallax image
		if (page == "main") {
			parallax("/static/img/saarschleife.jpg") {
				br()
				br()
				br()
				br()
				h5("header center white-text") { +"Explore the beauty of the Saarland" }
				br()
				br()
				a(if (visitor != null) "/visitor" else "#login", classes = "btn-large") {
					+"Get started"
				}
			}
		}

		div(classes = "container") {
			id = "main-content"
			br()
			br()

			layoutbody()

			br()
			br()
		}

		parallax("/static/img/saarschleife-treetop.jpg")

		footer("page-footer orange") {
			div("container") {
				h5("white-text") { +"We are saarschleife.net" }
				p("grey-text text-lighten-4") {
					+"The only and one Social Network for all Saarschleife Visitors"
				}
				p("grey-text text-lighten-4") {
					+"Share your Saarschleife Discoveries with all your friends!"
				}
			}
			div("footer-copyright") {
				div("container") {
					+"Proudly made by [saarsec] mkb, 2020"
				}
			}
		}

		script(src = "/static/js/materialize.js") {}
		script(src = "/static/kotlin.js") {}
		script(src = "/static/kotlinx-html-js.js") {}
		script(src = "/static/frontend.js") {}
	}
}
