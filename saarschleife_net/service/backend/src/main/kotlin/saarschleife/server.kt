package saarschleife

import io.ktor.application.Application
import io.ktor.application.call
import io.ktor.application.install
import io.ktor.features.*
import io.ktor.html.respondHtml
import io.ktor.http.CacheControl
import io.ktor.http.ContentType
import io.ktor.http.HttpStatusCode
import io.ktor.http.content.CachingOptions
import io.ktor.http.content.resources
import io.ktor.http.content.static
import io.ktor.jackson.jackson
import io.ktor.request.receive
import io.ktor.response.respond
import io.ktor.response.respondText
import io.ktor.routing.get
import io.ktor.routing.post
import io.ktor.routing.route
import io.ktor.routing.routing
import io.ktor.sessions.*
import io.ktor.util.date.GMTDate
import saarschleife.models.StatusMessage
import saarschleife.models.Storage
import saarschleife.models.VisitorAccount
import saarschleife.models.sha256
import saarschleife.template.layout
import saarschleife.template.pageHome
import saarschleife.template.pageSecrets
import saarschleife.template.pageVisitor


class PostDataLoginRegister(val username: String, val password: String)
class PostDataSecretFound(val name: String, val solution: String)
class PostDataStatusMessage(val title: String, val text: String, val isPublic: Boolean)

data class SessionData(val username: String)

open class HttpStatusException(val code: HttpStatusCode = HttpStatusCode.InternalServerError, message: String = "") : RuntimeException(message)
class NotFoundException(message: String = "Not found") : HttpStatusException(HttpStatusCode.NotFound, message)
class ForbiddenException(message: String = "Forbidden") : HttpStatusException(HttpStatusCode.Forbidden, message)

/**
 * Return current visitor, or throw exception if not logged in
 */
fun CurrentSession.currentVisitor() = Storage.visitors[this.get<SessionData>()?.username ?: ""]
		?: throw ForbiddenException("You have to be logged in")


fun Application.saarschleifeApp() {
	install(DefaultHeaders)
	install(AutoHeadResponse)
	install(ConditionalHeaders)
	install(PartialContent)
	install(Compression)
	install(CallLogging)

	// Caching assets
	install(CachingHeaders) {
		options { outgoingContent ->
			when (outgoingContent.contentType?.withoutParameters()) {
				ContentType.Text.CSS, ContentType.Application.JavaScript, ContentType.Image.JPEG ->
					CachingOptions(CacheControl.MaxAge(maxAgeSeconds = 7200, visibility = CacheControl.Visibility.Public), null as? GMTDate?)
				else -> null
			}
		}
	}


	// Authentication / Sessions
	install(Sessions) {
		cookie<SessionData>("saarschleife_session_id") {
			transform(SessionTransportTransformerMessageAuthentication(secretKey, "HmacSHA256"))
			cookie.path = "/"
		}
	}

	// JSON input / output
	install(ContentNegotiation) {
		jackson {}
	}

	// Exception handling
	install(StatusPages) {
		exception<HttpStatusException> { cause ->
			call.respondText(cause.message ?: cause.code.description, ContentType.Text.Plain, cause.code)
		}
	}


	routing {

		static("/static") {
			resources("")
		}

		get("/") {
			val currentVisitor = Storage.visitors[call.sessions.get<SessionData>()?.username ?: ""]
			call.respondHtml {
				layout("main", currentVisitor, onload = "HomePage.install()") { pageHome(currentVisitor) }
			}
		}

		get("/visitor") {
			val currentVisitor = call.sessions.currentVisitor()
			call.respondHtml {
				layout("visitor", currentVisitor, onload = "VisitorPage.install()") { pageVisitor() }
			}
		}

		get("/secrets") {
			val currentVisitor = call.sessions.currentVisitor()
			call.respondHtml {
				layout("secrets", currentVisitor, onload = "SecretPage.install()") { pageSecrets(Storage.secrets.values) }
			}
		}


		route("/api") {
			get("/me") {
				call.respond(OutputFormat.formatVisitor(call.sessions.currentVisitor(), true))
			}

			get("/visitor/{username}") {
				val visitor = Storage.visitors[call.parameters["username"]] ?: throw NotFoundException()
				call.respond(OutputFormat.formatVisitor(visitor, visitor.isFriendOf(call.sessions.currentVisitor())))
			}

			post("/visitor/friends/{username}") {
				val currentVisitor = call.sessions.currentVisitor()
				val otherVisitor = Storage.visitors[call.parameters["username"]] ?: throw NotFoundException()
				if (currentVisitor.points <= otherVisitor.points)
					throw ForbiddenException("${otherVisitor.username} does not want to be your friend, you loosy n00b!")
				currentVisitor.befriend(otherVisitor)
				call.respondText("OK")
			}

			route("/statusmessage") {
				post("set") {
					val visitor = call.sessions.currentVisitor()
					val post = call.receive<PostDataStatusMessage>()
					visitor.updateStatusMessage(StatusMessage(post.title, post.text, post.isPublic, visitor))
					call.respondText("OK")
				}

				post("remove") {
					call.sessions.currentVisitor().updateStatusMessage(null)
					call.respondText("OK")
				}

				get("{username}") {
					val visitor = Storage.visitors[call.parameters["username"]] ?: throw NotFoundException()
					if (visitor.statusMessage != null && !visitor.statusMessage!!.isPublic) {
						// check if this visitor is a friend
						val currentVisitor = call.sessions.currentVisitor()
						if (visitor.statusMessage?.writtenBy?.isFriendOf(currentVisitor) != true)
							throw ForbiddenException()
					}
					call.respond(OutputFormat.formatStatusMessage(visitor.statusMessage))
				}

				post("spread/{username}") {
					// Copy a visitor's status message, preserve the original author
					val currentVisitor = call.sessions.currentVisitor()
					val visitor = Storage.visitors[call.parameters["username"]] ?: throw NotFoundException()
					currentVisitor.updateStatusMessage(visitor.statusMessage)
					call.respondText("OK")
				}
			}



			post("/login-register") {
				val post = call.receive<PostDataLoginRegister>()
				if (post.username.isEmpty() || post.password.isEmpty())
					throw HttpStatusException(HttpStatusCode.InternalServerError, "Username / password required")
				val visitor = Storage.visitors.getOrElse(post.username) {
					Storage.addVisitor(VisitorAccount(post.username, sha256(post.password)))
				}
				if (visitor.password != sha256(post.password))
					throw ForbiddenException("Invalid credentials")
				call.sessions.set(SessionData(visitor.username))
				call.respondText("OK")
			}

			post("/logout") {
				call.sessions.clear<SessionData>()
				call.respondText("OK")
			}


			post("/secret/found") {
				val visitor = call.sessions.currentVisitor()
				val data = call.receive<PostDataSecretFound>()
				val challenge = Storage.secrets[data.name] ?: throw NotFoundException("Secret not found")
				if (challenge.hasBeenFoundBy(visitor))
					throw HttpStatusException(HttpStatusCode.BadRequest, "Already found")
				if (!challenge.checkSolution(data.solution))
					throw HttpStatusException(HttpStatusCode.BadRequest, "Wrong solution")

				challenge.setFoundBy(visitor)
				visitor.updatePoints(visitor.points + challenge.points)
				call.respondText("Secret found. You have ${visitor.points} points.")
			}
		}

	}
}
