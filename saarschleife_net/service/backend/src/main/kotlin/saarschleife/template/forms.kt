package saarschleife.template

import kotlinx.html.*

fun FORM.inputGroup(labeltext: String? = null, classes: String = "", inp: FORM.() -> Unit) {
	div(classes = "input-field $classes") {
		this@inputGroup.inp()
		if (labeltext != null)
			label {
				+labeltext
			}
	}
}

fun FORM.submit(text: String = "Submit") {
	button(classes = "btn waves-effect waves-light", type = ButtonType.submit) {
		i("material-icons left") {+"send"}
		+text
	}
}