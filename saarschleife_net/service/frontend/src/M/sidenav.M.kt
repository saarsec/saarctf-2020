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

external open class Sidenav : Component<SidenavOptions>, Openable {
    override fun open(): Unit = definedExternally
    override fun close(): Unit = definedExternally
    override var isOpen: Boolean = definedExternally
    open var isFixed: Boolean = definedExternally
    open var isDragged: Boolean = definedExternally
    companion object {
        fun getInstance(elem: Element): Sidenav = definedExternally
        fun init(els: Element, options: Any? = definedExternally /* null */): Sidenav = definedExternally
        fun init(els: Any, options: Any? = definedExternally /* null */): Array<Sidenav> = definedExternally
    }
}
external interface SidenavOptions {
    var edge: dynamic /* String /* "left" */ | String /* "right" */ */
    var draggable: Boolean
    var inDuration: Number
    var outDuration: Number
    var onOpenStart: (`this`: Sidenav, elem: Element) -> Unit
    var onOpenEnd: (`this`: Sidenav, elem: Element) -> Unit
    var onCloseStart: (`this`: Sidenav, elem: Element) -> Unit
    var onCloseEnd: (`this`: Sidenav, elem: Element) -> Unit
}
