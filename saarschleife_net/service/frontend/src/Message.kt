import M.toast

data class MyToastOptions(var html: String, var displayLength: Number = 5000, var classes: String = "")

object Message {

	fun success(message: String) {
		toast(MyToastOptions(message, classes = "light-green"))
	}

	fun error(message: String) {
		toast(MyToastOptions(message, classes = "red darken-2"))
	}
}