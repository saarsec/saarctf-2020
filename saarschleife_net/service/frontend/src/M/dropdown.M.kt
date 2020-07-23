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

external open class Dropdown : Component<DropdownOptions> {
    open var id: String = definedExternally
    open var dropdownEl: Element = definedExternally
    open var isOpen: Boolean = definedExternally
    open var isScrollable: Boolean = definedExternally
    open var focusedIndex: Number = definedExternally
    open fun open(): Unit = definedExternally
    open fun close(): Unit = definedExternally
    open fun recalculateDimensions(): Unit = definedExternally
    companion object {
        fun getInstance(elem: Element): Dropdown = definedExternally
        fun init(els: Element, options: Any? = definedExternally /* null */): Dropdown = definedExternally
        fun init(els: MElements, options: Any? = definedExternally /* null */): Array<Dropdown> = definedExternally
    }
}
external interface DropdownOptions {
    var alignment: dynamic /* String /* "left" */ | String /* "right" */ */
    var autoTrigger: Boolean
    var constrainWidth: Boolean
    var container: Element
    var coverTrigger: Boolean
    var closeOnClick: Boolean
    var hover: Boolean
    var inDuration: Number
    var outDuration: Number
    var onOpenStart: (`this`: Dropdown, el: Element) -> Unit
    var onOpenEnd: (`this`: Dropdown, el: Element) -> Unit
    var onCloseStart: (`this`: Dropdown, el: Element) -> Unit
    var onCloseEnd: (`this`: Dropdown, el: Element) -> Unit
}
