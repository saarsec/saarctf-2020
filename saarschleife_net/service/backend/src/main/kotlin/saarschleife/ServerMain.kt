package saarschleife

import saarschleife.models.Storage
import io.ktor.application.Application
import io.ktor.server.engine.embeddedServer
import io.ktor.server.netty.Netty
import kotlinx.coroutines.runBlocking


object ServerMain {
	@JvmStatic
	fun main(args: Array<String>) {
		var host = "0.0.0.0"
		var port = 8080
		if (args.size == 2) {
			host = args[0]
			port = args[1].toInt()
		} else if (args.isNotEmpty()) {
			when (args[0]) {
				"clean" -> {
					runBlocking {
						Storage.load()
						Storage.cleanStorage()
					}
					println("Database wiped, please restart server.")
				}
				"init" -> {
					Storage.load()
					println("Default entries are present now")
				}
				else -> {
					println("Parameters: [<host> <port>]     // start server")
					println("            clean               // wipe database")
					println("            init                // add default entries, without wiping database")
				}
			}
			return
		}
		println("Binding to $host:$port ...")

		// Initialize database
		Storage.load()
		// Start server
		embeddedServer(Netty, module = Application::saarschleifeApp, host = host, port = port, configure = {
			// limit number of used cores (default config for 2 cores)
			connectionGroupSize = 2
			workerGroupSize = 2
			callGroupSize = 2
			shareWorkGroup = false
		}).start(wait = true)
	}
}