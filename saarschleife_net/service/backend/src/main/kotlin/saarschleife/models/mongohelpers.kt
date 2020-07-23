package saarschleife.models

import com.mongodb.async.client.MongoCollection
import com.mongodb.client.model.CountOptions
import com.mongodb.client.model.Filters
import org.bson.conversions.Bson
import org.litote.kmongo.Id
import org.litote.kmongo.and
import org.litote.kmongo.coroutine.*

/**
 * Count the number of entries that fulfill a condition
 */
suspend fun <T> MongoCollection<T>.countDocuments(filter: Bson, options: CountOptions = CountOptions()): Long {
	return singleResult { countDocuments(filter, options, it) } ?: 0L
}

/**
 * Check if a condition holds for an entry identified by _id
 */
suspend fun <T> MongoCollection<T>.check(id: Id<T>, condition: Bson): Boolean {
	val filter = and(Filters.eq("_id", id), condition)
	return countDocuments(filter = filter) > 0
}