package saarschleife.template

import saarschleife.models.Secret
import saarschleife.models.Storage
import saarschleife.models.VisitorAccount
import kotlinx.html.*
import java.text.SimpleDateFormat

val DateFormatter = SimpleDateFormat("dd.MM.yyyy HH:mm")
val ShortDateFormatter = SimpleDateFormat("HH:mm:ss")

fun DIV.pageHome(visitor: VisitorAccount?) {
	h1(classes = "header center orange-text") { +"Welcome to saarschleife.net!" }
	div("row") {
		div("col m8 offset-m2 l6 offset-l3") {

			if (visitor != null) {
				h5(classes = "header light center") { +"Your Account" }

				br()

				dl("row") {
					dt("col s4") { +"Username:" }
					dd("col s8") { +visitor.username }
					dt("col s4") { +"Points:" }
					dd("col s8") { +"${visitor.points}" }
					dt("col s4") { +"Member since:" }
					dd("col s8") { +DateFormatter.format(visitor.created) }
				}

				br()

				div("row center") {
					a(href = "/visitor", classes = "btn-large waves-effect waves-light orange") {
						+"Discover other members!"
					}
				}

			} else {

				// Login form
				h5(classes = "header light") {
					id = "login"
					+"Please log in or register to continue"
				}
				br()
				form("/api/login-register", classes = "form-login") {
					inputGroup("Username") {
						textInput(name = "username") { required = true }
					}
					inputGroup("Password") {
						passwordInput(name = "password") { required = true }
					}
					submit("Login / Register")
				}
			}
		}
	}
}

fun DIV.pageVisitor() {
	style {
		unsafe {
			+".link-visitor small{ margin-left: 10px }"
		}
	}

	h1(classes = "header center orange-text") { +"Visitors" }
	h5("header center light") { +"Socialize with other visitors!" }
	br()

	form(classes = "form-visitor-search") {
		inputGroup("Search visitor", classes = "inline") { input(type = InputType.text, name = "username") { id = "visitor-search-input" } }
		button(classes = "btn waves-effect waves-light", type = ButtonType.submit) {
			i("material-icons") { +"search" }
		}
	}
	div("card") {
		id = "visitor-content"
	}

	br()

	// Form to create/edit status message
	form("/api/statusmessage/set", classes = "hide form-statusmessage-set card") {
		div(classes = "card-content") {
			span("card-title") { +"Write a status message" }
			this@form.inputGroup("Title") {
				input(type = InputType.text, name = "title") { required = true }
			}
			this@form.inputGroup("Text") {
				textArea(classes = "materialize-textarea") { name = "text"; required = true }
			}
			p {
				label {
					input(type = InputType.checkBox, name = "isPublic") { }
					span { +"public" }
				}
			}
		}
		div("card-action") {
			this@form.submit("Save")
		}
	}

	br()
	br()

	// List of recent visitors
	if (Storage.recentVisitors.isNotEmpty()) {
		h4(classes = "header center orange-text") { +"Other visitors" }
		div("collection") {
			val recentAfter = Storage.visitorIsRecentAfterDate()
			for (member in Storage.recentVisitors) {
				if (!member.created.after(recentAfter))
					break

				a("#", classes = "collection-item link-visitor") {
					attributes["data-username"] = member.username
					span("badge") { +"${member.points} points" }
					+member.username
					+" "
					small {
						+"(since "
						+ShortDateFormatter.format(member.created)
						+")"
					}
				}
			}
		}
	}
}

fun DIV.pageSecrets(secrets: Iterable<Secret>) {
	h1(classes = "header center orange-text") { +"Secrets" }
	h5("header center light") { +"Explore and find secrets to get respected by other visitors!" }
	br()

	div(classes = "collection") {
		secrets.sortedBy { c -> c.points }.forEach { c ->
			a("#", classes = "collection-item link-secret") {
				attributes["data-name"] = c.name
				+c.name
				span("badge") { +"${c.points} points" }
			}
		}
	}

	br()

	form(classes = "hide card") {
		id = "form-secret-find"
		div("card-content") {
			span("card-title") { +"Found a secret? Tell us!" }
			this@form.inputGroup("Secret Text") { input(InputType.text, name = "solution") { required = true } }
		}
		div("card-action") {
			this@form.submit("Submit")
		}
	}
}
