package saarschleife

fun String.hexStringToByteArray() = ByteArray(this.length / 2) { this.substring(it * 2, it * 2 + 2).toInt(16).toByte() }

// Auto-generated at build time. Rebuild project to change the key.
val secretKey = "21E27C652F7BE883FFC4FC0A81BF3830437F2052D2EBC2FF2B5C31E17B64327A".hexStringToByteArray()
