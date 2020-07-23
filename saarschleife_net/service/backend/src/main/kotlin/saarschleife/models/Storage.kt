package saarschleife.models

import com.mongodb.client.model.Updates
import kotlinx.coroutines.runBlocking
import org.litote.kmongo.async.KMongo
import org.litote.kmongo.async.getCollection
import org.litote.kmongo.coroutine.updateOneById
import org.litote.kmongo.coroutine.countDocuments
import org.litote.kmongo.coroutine.deleteMany
import org.litote.kmongo.coroutine.insertOne
import org.litote.kmongo.coroutine.toList
import java.util.*
import java.util.concurrent.ConcurrentHashMap
import java.util.concurrent.ConcurrentLinkedDeque


object Storage {
	private val client = KMongo.createClient()
	private val database = client.getDatabase("saarschleife")

	val visitors: MutableMap<String, VisitorAccount> = ConcurrentHashMap()
	val recentVisitors: ConcurrentLinkedDeque<VisitorAccount> = ConcurrentLinkedDeque()
	val visitorsCollection = database.getCollection<VisitorAccount>()

	// "recent" visitors are visitors that registered within the last 30 minutes
	fun visitorIsRecentAfterDate() = Date(System.currentTimeMillis() - 1000 * 60 * 30)

	suspend fun addVisitor(acc: VisitorAccount): VisitorAccount {
		if (visitors.putIfAbsent(acc.username, acc) != null)
			throw IllegalArgumentException("Visitor already exists")
		recentVisitors.addFirst(acc)

		visitorsCollection.insertOne(acc)

		println("- New visitor: ${acc.username}")
		return acc
	}


	val secrets: MutableMap<String, Secret> = ConcurrentHashMap()
	val secretsCollection = database.getCollection<Secret>()

	suspend fun addSecret(s: Secret) {
		if (secrets.putIfAbsent(s.name, s) != null)
			throw IllegalArgumentException("Secret already exists")

		secretsCollection.insertOne(s)
	}


	private suspend fun initStorage() {
		if (secretsCollection.countDocuments() == 0L) {
			println("[DB] Add default secrets...")
			addSecret(Secret("Cloef", points = 200))
			addSecret(Secret("Treetop Path", points = 300))
			addSecret(Secret("Observation Tower", points = 500))
		} else {
			println("[DB] Regenerate default secrets...")
			secretsCollection.find().toList().forEach { c: Secret -> secrets[c.name] = c }
		}
		secretsCollection.updateOneById(secrets["Cloef"]!!._id, Updates.set("solution", secretSolution("Cloef")))
		secretsCollection.updateOneById(secrets["Treetop Path"]!!._id, Updates.set("solution", secretSolution("Treetop Path")))
		secretsCollection.updateOneById(secrets["Observation Tower"]!!._id, Updates.set("solution", secretSolution("Observation Tower")))
	}

	private suspend fun loadStorage() {
		println("[DB] Loading visitors")
		val recentAfter = visitorIsRecentAfterDate()
		visitorsCollection.find().toList().sortedBy { it.created }.forEach { u: VisitorAccount ->
			visitors[u.username] = u
			if (u.created.after(recentAfter))
				recentVisitors.addFirst(u)
		}
		secretsCollection.find().toList().forEach { c: Secret -> secrets[c.name] = c }
		println("[DB] Server ready.   #visitors: ${visitors.size}   #secrets: ${secrets.size}")
	}

	suspend fun cleanStorage() {
		visitorsCollection.deleteMany()
		secretsCollection.deleteMany()
	}

	fun load() {
		runBlocking {
			initStorage()
			loadStorage()
		}
	}

}
