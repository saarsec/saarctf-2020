@file:Suppress("INTERFACE_WITH_SUPERCLASS", "OVERRIDING_FINAL_MEMBER", "RETURN_TYPE_MISMATCH_ON_OVERRIDE", "CONFLICTING_OVERLOADS", "EXTERNAL_DELEGATION", "NESTED_CLASS_IN_EXTERNAL_INTERFACE")
@file:JsQualifier("M")
package M

import kotlin.js.*
import kotlin.js.Json
import org.khronos.webgl.*
import org.w3c.dom.*
import org.w3c.dom.events.*
import org.w3c.dom.parsing.*
import org.w3c.dom.svg.*
import org.w3c.dom.url.*
import org.w3c.fetch.*
import org.w3c.files.*
import org.w3c.notifications.*
import org.w3c.performance.*
import org.w3c.workers.*
import org.w3c.xhr.*

external open class MElements : NodeList

external open class Component<TOptions>(elem: Element, options: Any? = definedExternally /* null */) : ComponentBase<TOptions> {
    open fun destroy(): Unit = definedExternally
}
external open class ComponentBase<TOptions>(options: Any? = definedExternally /* null */) {
    open var el: Element = definedExternally
    open var options: TOptions = definedExternally
}
external interface Openable {
    var isOpen: Boolean
    fun open()
    fun close()
}
external interface InternationalizationOptions {
    var cancel: String
    var clear: String
    var done: String
    var previousMonth: String
    var nextMonth: String
    var months: Array<String>
    var monthsShort: Array<String>
    var weekdays: Array<String>
    var weekdaysShort: Array<String>
    var weekdaysAbbrev: Array<String>
}
