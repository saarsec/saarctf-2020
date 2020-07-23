package saarschleife

import saarschleife.models.StatusMessage
import saarschleife.models.VisitorAccount

object OutputFormat {

	fun formatVisitor(visitor: VisitorAccount, isFriend: Boolean): Map<String, Any> = mapOf(
			"username" to visitor.username,
			"points" to visitor.points,
			"statusMessage" to (visitor.statusMessage != null),
			"isFriend" to isFriend,
			"created" to visitor.created.time / 1000
	)

	fun formatStatusMessage(message: StatusMessage?) = if (message == null) mapOf() else mapOf(
			"title" to message.title,
			"text" to message.text,
			"isPublic" to message.isPublic,
			"writtenBy" to message.writtenBy.username
	)

}