import LayoutScripts.defaultErrorHandler
import LayoutScripts.getFormData
import kotlinx.html.*
import kotlinx.html.dom.create
import kotlinx.html.js.div
import kotlinx.html.js.onClickFunction
import org.w3c.dom.*
import org.w3c.dom.events.Event
import kotlin.browser.document

@Suppress("unused")
object VisitorPage {
	private var myUsername: String = ""
	private var selectedVisitor: ClientVisitor? = null
	private var selectedStatusMessage: ClientStatusMessage? = null
	private lateinit var statusForm: HTMLFormElement

	/**
	 * Create HTML to show a visitor
	 */
	fun displayVisitor(visitor: ClientVisitor?) {
		val div = document.create.div(classes = "card-content") {
			if (visitor != null) {
				if (visitor.username == myUsername) {
					span("new badge light-blue") {
						//i("material-icons left") { +"account_circle" }
						attributes["data-badge-caption"] = "that's you!"
					}
				} else if (visitor.isFriend) {
					span("new badge green") {
						i("material-icons left") { +"supervisor_account" }
						attributes["data-badge-caption"] = "You are friends!"
					}
				} else if (!visitor.isFriend) {
					// "Add as friend" button
					a(classes = "btn-floating btn-large halfway-fab top-fab waves-effect waves-light light-blue") {
						title = "Add as friend"
						onClickFunction = VisitorPage::addFriend
						i("material-icons") { +"supervisor_account" }
					}
				}

				h4 { +"Visitor: ${visitor.username}" }
				p {
					strong { +"Points: " }
					+"${visitor.points}"
				}
				p {
					strong { +"Status: " }
				}
				div {
					attributes["id"] = "status-message"
				}
			} else {
				h5 { +"No visitor found" }
			}
		}
		val container = document.getElementById("visitor-content")
		container?.innerHTML = ""
		container?.append(div)

		selectedVisitor = visitor
		statusForm.classList.add("hide")
	}

	/**
	 * Create HTML to show a status message
	 */
	fun displayStatusMessage(statusMessage: ClientStatusMessage?) {
		val tag = document.create.div {
			if (statusMessage != null) {
				blockQuote {
					h6 { +statusMessage.title }
					+statusMessage.text
					p {
						small { +"by ${statusMessage.writtenBy}" }
						if (statusMessage.isPublic) {
							span(classes = "new badge light-blue nofloat") { attributes["data-badge-caption"] = "public" }
						} else {
							span(classes = "new badge red nofloat") { attributes["data-badge-caption"] = "private" }
						}
					}
				}
				if (selectedVisitor?.username == myUsername) {
					a(classes = "btn-floating right red") {
						i("material-icons") { +"delete" }
						onClickFunction = VisitorPage::removeStatusMessage
					}
					a(classes = "btn-floating right waves-effect waves-light light-blue") {
						i("material-icons") { +"edit" }
						onClickFunction = VisitorPage::editStatusMessage
					}
				} else {
					a(classes = "btn-floating right waves-effect waves-light light-blue") {
						i("material-icons") { +"share" }
						onClickFunction = VisitorPage::spreadStatusMessage
					}
				}
			} else {
				span { +"no status message" }
				if (selectedVisitor?.username == myUsername)
					a(classes = "btn-floating right waves-effect waves-light red") {
						i("material-icons") { +"add" }
						onClickFunction = VisitorPage::editStatusMessage
					}
			}
		}
		val container = document.getElementById("status-message")
		container?.innerHTML = ""
		container?.append(tag)

		selectedStatusMessage = statusMessage
	}

	/**
	 * Create HTML to show that a status message is not available
	 */
	fun displayInaccessibleStatusMessage() {
		val tag = document.create.div {
			if (selectedVisitor?.username == myUsername) {
				a(classes = "btn-floating right red") {
					i("material-icons") { +"delete" }
					onClickFunction = VisitorPage::removeStatusMessage
				}
				a(classes = "btn-floating right waves-effect waves-light light-blue") {
					i("material-icons") { +"edit" }
					onClickFunction = VisitorPage::editStatusMessage
				}
			}
			i { +"status message inaccessible (you're not a friend of the author)" }
		}

		val container = document.getElementById("status-message")
		container?.innerHTML = ""
		container?.append(tag)

		selectedStatusMessage = null
	}


	fun searchVisitor(e: Event) {
		e.preventDefault()
		val username = e.currentTarget.unsafeCast<HTMLFormElement>().elements[0]?.unsafeCast<HTMLInputElement>()?.value
		searchVisitor(username)
	}

	private fun searchVisitor(username: String?) {
		getVisitor(username) { visitor ->
			displayVisitor(visitor)
			if (visitor?.statusMessage == true) {
				LayoutScripts.get<ClientStatusMessage>("/api/statusmessage/${visitor.username}", success = this::displayStatusMessage, error = { xhr ->
					if (xhr.status == 403.toShort()) {
						displayInaccessibleStatusMessage()
					} else {
						defaultErrorHandler(xhr)
					}
				})
			} else {
				displayStatusMessage(null)
			}
		}
	}

	private fun clickVisitor(e: Event) {
		e.preventDefault()
		val username = e.currentTarget.unsafeCast<HTMLElement>().attributes["data-username"]?.value
		val searchField = document.getElementById("visitor-search-input").unsafeCast<HTMLInputElement>()
		searchVisitor(username)
		searchField.value = username ?: ""
		searchField.scrollIntoView()
		M.updateTextFields()
	}

	private fun getVisitor(username: String?, cb: (ClientVisitor?) -> Unit) {
		LayoutScripts.get<ClientVisitor>(if (username != null) "/api/visitor/$username" else "/api/me", error = { xhr ->
			if (xhr.status == 404.toShort())
				cb(null)
			else
				defaultErrorHandler(xhr)
		}, success = cb)
	}

	fun addFriend(e: Event) {
		e.preventDefault()
		LayoutScripts.post("/api/visitor/friends/" + selectedVisitor?.username, null) {
			Message.success("You are now a friend of ${selectedVisitor?.username}")
			searchVisitor(selectedVisitor?.username)
		}
	}


	fun editStatusMessage(e: Event) {
		e.preventDefault()
		statusForm.elements["title"]?.unsafeCast<HTMLInputElement>()?.value = selectedStatusMessage?.title ?: ""
		statusForm.elements["text"]?.unsafeCast<HTMLTextAreaElement>()?.value = selectedStatusMessage?.text ?: ""
		statusForm.elements["isPublic"]?.unsafeCast<HTMLInputElement>()?.checked = selectedStatusMessage?.isPublic ?: false
		statusForm.classList.remove("hide")
		M.updateTextFields()
		statusForm.elements["title"]?.unsafeCast<HTMLInputElement>()?.focus()
	}

	fun submitStatusMessage(e: Event) {
		e.preventDefault()
		val form = e.currentTarget.unsafeCast<HTMLFormElement>()
		val endpoint = form.attributes["action"]?.value ?: throw RuntimeException("No action")
		LayoutScripts.post(endpoint, getFormData(form)) {
			Message.success("Status message saved")
			searchVisitor(null)
		}
	}

	fun removeStatusMessage(e: Event) {
		e.preventDefault()
		LayoutScripts.post("/api/statusmessage/remove", null) {
			Message.success("Status message removed")
			searchVisitor(null)
		}
	}

	fun spreadStatusMessage(e: Event) {
		e.preventDefault()
		LayoutScripts.post("/api/statusmessage/spread/${selectedVisitor?.username}", null) {
			Message.success("You shared this status message")
			searchVisitor(null)
		}
	}

	fun install() {
		LayoutScripts.install()
		myUsername = document.body?.attributes?.get("data-my-username")?.value ?: ""
		statusForm = document.getElementsByClassName("form-statusmessage-set")[0] as HTMLFormElement
		statusForm.addEventListener("submit", VisitorPage::submitStatusMessage)
		document.getElementsByClassName("form-visitor-search").asList().forEach { it.addEventListener("submit", VisitorPage::searchVisitor) }
		document.getElementsByClassName("link-visitor").asList().forEach { it.addEventListener("click", VisitorPage::clickVisitor) }
		searchVisitor(null)
	}
}
