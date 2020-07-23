import org.w3c.dom.HTMLFormElement
import org.w3c.dom.asList
import org.w3c.dom.events.Event
import org.w3c.dom.get
import kotlin.browser.document
import kotlin.browser.window

@Suppress("unused")
object HomePage {
	private fun login(e: Event) {
		e.preventDefault()
		val form = e.currentTarget.unsafeCast<HTMLFormElement>()
		val endpoint = form.attributes["action"]?.value ?: throw RuntimeException("No action")
		LayoutScripts.post(endpoint, LayoutScripts.getFormData(form)) {
			window.location.href = "/visitor"
		}
	}

	fun install() {
		LayoutScripts.install()
		document.getElementsByClassName("form-login").asList().forEach { it.addEventListener("submit", HomePage::login) }
	}
}