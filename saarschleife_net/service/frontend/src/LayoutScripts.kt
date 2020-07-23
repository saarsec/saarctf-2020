import org.w3c.dom.*
import org.w3c.dom.events.Event
import org.w3c.xhr.XMLHttpRequest
import kotlin.browser.document
import kotlin.browser.window


object LayoutScripts {

	fun defaultErrorHandler(xhr: XMLHttpRequest) {
		console.error(xhr)
		if (!xhr.responseText.isEmpty())
			Message.error(xhr.responseText)
	}

	private fun get(url: String, error: (XMLHttpRequest) -> Unit = LayoutScripts::defaultErrorHandler, success: (XMLHttpRequest) -> Unit) {
		val xhr = XMLHttpRequest()
		xhr.open("GET", url, true)
		xhr.onreadystatechange = {
			if (xhr.readyState == 4.toShort()) {
				if (xhr.status in 200..299) {
					success(xhr)
				} else {
					error(xhr)
				}
			}
		}
		xhr.send()
	}

	fun <T> get(url: String, error: (XMLHttpRequest) -> Unit = LayoutScripts::defaultErrorHandler, success: (T) -> Unit) {
		get(url, error) { xhr ->
			success(JSON.parse(xhr.responseText))
		}
	}

	fun post(url: String, body: Any?, error: (XMLHttpRequest) -> Unit = LayoutScripts::defaultErrorHandler, success: (XMLHttpRequest) -> Unit) {
		val xhr = XMLHttpRequest()
		xhr.open("POST", url, true)
		xhr.setRequestHeader("Content-Type", "application/json")
		xhr.onreadystatechange = {
			if (xhr.readyState == 4.toShort()) {
				if (xhr.status in 200..299) {
					success(xhr)
				} else {
					error(xhr)
				}
			}
		}
		xhr.send(JSON.stringify(body))
	}

	fun getFormData(form: HTMLFormElement): Any {
		val data: dynamic = object {}
		form.elements.asList().forEach { input ->
			when (input) {
				is HTMLInputElement -> {
					if (input.type == "checkbox") {
						data[input.name] = input.checked
					} else {
						data[input.name] = input.value
					}
				}
				is HTMLTextAreaElement -> {
					data[input.name] = input.value
				}
			}
		}
		return data
	}

	private fun logout(e: Event) {
		e.preventDefault()
		post("/api/logout", null) {
			window.location.href = "/"
		}
	}

	fun install() {
		document.getElementsByClassName("link-logout").asList().forEach { it.addEventListener("click", LayoutScripts::logout) }

		val elements = document.querySelector(".sidenav")
		if (elements != null) M.Sidenav.init(elements)

		//val parallax = document.querySelector(".parallax")
		//if (parallax != null) M.Parallax.init(parallax)
		document.getElementsByClassName("parallax").asList().forEach {
			M.Parallax.init(it)
		}
	}
}
