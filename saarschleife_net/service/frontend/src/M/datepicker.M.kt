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

external open class Datepicker : Component<DatepickerOptions>, Openable {
    override var isOpen: Boolean = definedExternally
    open var date: Date = definedExternally
    open var doneBtn: HTMLButtonElement = definedExternally
    open var clearBtn: HTMLButtonElement = definedExternally
    override fun open(): Unit = definedExternally
    override fun close(): Unit = definedExternally
    override fun toString(): String = definedExternally
    open fun setDate(date: String? = definedExternally /* null */, preventOnSelect: Boolean? = definedExternally /* null */): Unit = definedExternally
    open fun setDate(date: Date? = definedExternally /* null */, preventOnSelect: Boolean? = definedExternally /* null */): Unit = definedExternally
    open fun gotoDate(date: Date): Unit = definedExternally
    open fun setInputValue(): Unit = definedExternally
    companion object {
        fun getInstance(elem: Element): Datepicker = definedExternally
        fun init(els: Element, options: Any? = definedExternally /* null */): Datepicker = definedExternally
        fun init(els: MElements, options: Any? = definedExternally /* null */): Array<Datepicker> = definedExternally
    }
    open fun setDate(): Unit = definedExternally
}
external interface DatepickerOptions {
    var autoClose: Boolean
    var format: String
    var parse: (value: String, format: String) -> Date
    var defaultDate: Date
    var setDefaultDate: Boolean
    var disableWeekends: Boolean
    var disableDayFn: (day: Date) -> Boolean
    var firstDay: Number
    var minDate: Date
    var maxDate: Date
    var yearRange: dynamic /* Number | Array<Number> */
    var isRTL: Boolean
    var showMonthAfterYear: Boolean
    var showDaysInNextAndPreviousMonths: Boolean
    var container: Element
    var showClearBtn: Boolean
    var i18n: Any?
    var events: Array<String>
    var onSelect: (`this`: Datepicker, selectedDate: Date) -> Unit
    var onOpen: (`this`: Datepicker) -> Unit
    var onClose: (`this`: Datepicker) -> Unit
    var onDraw: (`this`: Datepicker) -> Unit
}
