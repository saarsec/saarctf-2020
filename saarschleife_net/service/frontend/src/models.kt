data class ClientVisitor(val username: String, val points: Int, val statusMessage: Boolean, val isFriend: Boolean)
data class ClientStatusMessage(val title: String, val text: String, val isPublic: Boolean, val writtenBy: String)
