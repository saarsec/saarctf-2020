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

external open class Timepicker : Component<TimepickerOptions> {
    open var isOpen: Boolean = definedExternally
    open var time: String = definedExternally
    open fun open(): Unit = definedExternally
    open fun close(): Unit = definedExternally
    open fun showView(view: String /* "hours" */): Unit = definedExternally
    open fun showView(view: String /* "minutes" */): Unit = definedExternally
    companion object {
        fun getInstance(elem: Element): Timepicker = definedExternally
        fun init(els: Element, options: Any? = definedExternally /* null */): Timepicker = definedExternally
        fun init(els: MElements, options: Any? = definedExternally /* null */): Array<Timepicker> = definedExternally
    }
}
external interface TimepickerOptions {
    var duration: Number
    var container: String
    var showClearBtn: Boolean
    var defaultTime: String
    var fromNow: Number
    var i18n: Any?
    var autoClose: Boolean
    var twelveHour: Boolean
    var vibrate: Boolean
    var onOpenStart: (`this`: Modal, el: Element) -> Unit
    var onOpenEnd: (`this`: Modal, el: Element) -> Unit
    var onCloseStart: (`this`: Modal, el: Element) -> Unit
    var onCloseEnd: (`this`: Modal, el: Element) -> Unit
    var onSelect: (`this`: Modal, hour: Number, minute: Number) -> Unit
}
