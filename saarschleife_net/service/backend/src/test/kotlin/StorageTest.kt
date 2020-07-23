import kotlinx.coroutines.*
import kotlinx.coroutines.time.delay
import saarschleife.models.StatusMessage
import saarschleife.models.Storage
import saarschleife.models.VisitorAccount
import org.junit.Assert.*
import org.junit.Test
import org.litote.kmongo.coroutine.findOneById
import java.time.Duration
import java.util.*
import kotlin.coroutines.CoroutineContext
import kotlin.coroutines.resume
import kotlin.streams.asSequence


@InternalCoroutinesApi
class TestDirectContext : CoroutineDispatcher(), Delay {
	override fun scheduleResumeAfterDelay(timeMillis: Long, continuation: CancellableContinuation<Unit>) {
		Thread.sleep(timeMillis, 0)
		continuation.resume(Unit)
	}

	override fun dispatch(context: CoroutineContext, block: Runnable) {
		block.run()
	}
}

@InternalCoroutinesApi
class StorageTest {

	private val alphanumChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	private fun randomString(len: Long) =
			Random().ints(len, 0, alphanumChars.length).asSequence().map(alphanumChars::get).joinToString("")

	private val context = TestDirectContext()

	private suspend fun createUser(prefix: String = ""): VisitorAccount {
		val username = prefix + randomString(10)
		val user = VisitorAccount(username, randomString(12))
		Storage.addVisitor(user)
		return user
	}

	@Test
	fun testUserFriends() {
		runBlocking(context) {
			Storage.load()
			val user1 = createUser("test")
			val user2 = createUser("test")

			assertTrue(user1.isFriendOf(user1))
			assertFalse(user1.isFriendOf(user2))
			assertFalse(user2.isFriendOf(user1))
			user1.befriend(user2)
			assertTrue("user1 is friend of user2", user1.isFriendOf(user2))
			assertTrue("user2 is friend of user1", user2.isFriendOf(user1))

			user1.befriend(user2)
			assertTrue(user1.isFriendOf(user2))
			assertTrue(user2.isFriendOf(user1))
		}
	}


	@Test
	fun testSolveOnlyOnce() {
		runBlocking(context) {
			Storage.load()
			val user = createUser("test")
			val challenge = Storage.secrets.values.first()

			assertFalse(challenge.hasBeenFoundBy(user))
			challenge.setFoundBy(user)
			assertTrue(challenge.hasBeenFoundBy(user))
		}
	}

	@Test
	fun testUserPropertyChange() {
		runBlocking(context) {
			Storage.load()
			val user = createUser("test")
			assertEquals(0, user.points)

			user.updatePoints(350)
			assertEquals(350, user.points)

			val user2 = Storage.visitorsCollection.findOneById(user._id) ?: throw AssertionError("Not found")
			assertEquals(350, user2.points)
		}
	}

	@Test
	fun testStatusMessage() {
		runBlocking(context) {
			Storage.load()
			val user = createUser("test")
			assertNull(user.statusMessage)

			user.updateStatusMessage(StatusMessage("TestTitle", "TestText", true, user))
			assertEquals("TestTitle", user.statusMessage?.title)

			val user2 = Storage.visitorsCollection.findOneById(user._id) ?: throw AssertionError("Not found")
			assertEquals("TestTitle", user2.statusMessage?.title)
			assertEquals(user.username, user2.statusMessage?.writtenBy?.username)
		}
	}

	@Test
	fun testTiming() {
		runBlocking(context) {
			Storage.load()
			val user = createUser("test")
			val challenge = Storage.secrets.values.first()

			println("-------------------------")
			for (i in 1..10) {
				val t = System.nanoTime()
				challenge.hasBeenFoundBy(user)
				val dt = System.nanoTime() - t
				println("  Time: ${dt / 1000 / 1000.0} ms")
				delay(150)
			}
			println("-------------------------")
		}
	}
}