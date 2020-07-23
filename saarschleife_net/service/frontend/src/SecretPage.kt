import org.w3c.dom.*
import org.w3c.dom.events.Event
import kotlin.browser.document
import kotlin.dom.addClass
import kotlin.dom.removeClass

@Suppress("unused")
object SecretPage {

	var secretName: String? = null

	private fun selectSecret(e: Event) {
		e.preventDefault()
		secretName = e.currentTarget?.unsafeCast<HTMLLinkElement>()?.attributes?.get("data-name")?.value
		val form = document.getElementById("form-secret-find")!!.unsafeCast<HTMLFormElement>()
		form.removeClass("hide")
		val input = form.elements["solution"].unsafeCast<HTMLInputElement>()
		input.value = ""
		input.focus()
	}

	private fun findSecret(e: Event) {
		e.preventDefault()
		val form = document.getElementById("form-secret-find")!!.unsafeCast<HTMLFormElement>()
		form.addClass("hide")
		val data: dynamic = object {}
		data["name"] = secretName
		data["solution"] = form.elements["solution"].unsafeCast<HTMLInputElement>().value
		LayoutScripts.post("/api/secret/found", data) {
			Message.success(it.responseText)
		}
	}

	fun install() {
		LayoutScripts.install()
		document.getElementsByClassName("link-secret").asList().forEach { it.addEventListener("click", SecretPage::selectSecret) }
		document.getElementById("form-secret-find")?.addEventListener("submit", SecretPage::findSecret)
	}
}