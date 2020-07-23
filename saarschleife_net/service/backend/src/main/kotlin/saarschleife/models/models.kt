package saarschleife.models

import com.fasterxml.jackson.annotation.JsonIgnore
import com.mongodb.client.model.Filters
import com.mongodb.client.model.Updates
import org.litote.kmongo.Id
import org.litote.kmongo.coroutine.*
import org.litote.kmongo.newId
import org.litote.kmongo.set
import java.net.NetworkInterface
import java.security.MessageDigest
import java.util.*


fun sha256(pwd: String): String {
	return MessageDigest.getInstance("SHA-256").digest(pwd.toByteArray()).fold("") { str, it -> str + "%02x".format(it) }
}

fun secretSolution(name: String): String {
	val xyz = NetworkInterface.getNetworkInterfaces().toList().flatMap { ni ->
		ni.inetAddresses.toList().filter { it.address.size == 4 }.filter { !it.isLoopbackAddress }.map { it.hostAddress }
	}.first().split(".")[2]
	return "geocache{" + sha256("$name|$xyz").substring(0..11) + "}"
}


data class VisitorAccount(val username: String, val password: String,
						  var points: Int = 0, var statusMessage: StatusMessage? = null,
						  val created: Date = Date(), val _id: Id<VisitorAccount> = newId()) {

	suspend fun isFriendOf(visitor: VisitorAccount): Boolean {
		return Storage.visitorsCollection.check(_id, Filters.eq("friends", visitor.username)) || visitor.username == this.username
	}

	suspend fun befriend(visitor: VisitorAccount) {
		Storage.visitorsCollection.updateOneById(_id, Updates.addToSet("friends", visitor.username))
		Storage.visitorsCollection.updateOneById(visitor._id, Updates.addToSet("friends", username))
	}

	suspend fun updatePoints(points: Int) {
		this.points = points
		Storage.visitorsCollection.updateOneById(_id, set(VisitorAccount::points, points))
	}

	suspend fun updateStatusMessage(statusMessage: StatusMessage?) {
		this.statusMessage = statusMessage
		Storage.visitorsCollection.updateOneById(_id, set(VisitorAccount::statusMessage, statusMessage))
	}

}

data class StatusMessage(val title: String, val text: String, val isPublic: Boolean, var writtenByUsername: String) {

	// shorthand: StatusMessage.writtenBy == visitors[StatusMessage.writtenByUsername]

	constructor(title: String, text: String, isPublic: Boolean, writtenBy: VisitorAccount) : this(title, text, isPublic, writtenBy.username)

	var writtenBy: VisitorAccount
		@JsonIgnore get() = Storage.visitors[writtenByUsername]
				?: throw RuntimeException("User $writtenByUsername not found")
		@JsonIgnore set(u) {
			writtenByUsername = u.username
		}

}

data class Secret(val name: String, val points: Int, val _id: Id<Secret> = newId()) {

	suspend fun checkSolution(solution: String): Boolean {
		return Storage.secretsCollection.check(_id, Filters.eq("solution", solution))
	}

	suspend fun hasBeenFoundBy(visitor: VisitorAccount): Boolean {
		return Storage.secretsCollection.check(_id, Filters.eq("foundBy", visitor.username))
	}

	suspend fun setFoundBy(visitor: VisitorAccount) {
		Storage.secretsCollection.updateOneById(_id, Updates.addToSet("foundBy", visitor.username))
	}
}
