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

external open class Materialbox : Component<MaterialboxOptions> {
    open var overlayActive: Boolean = definedExternally
    open var doneAnimating: Boolean = definedExternally
    open var caption: String = definedExternally
    open var originalWidth: Number = definedExternally
    open var originalHeight: Number = definedExternally
    open fun open(): Unit = definedExternally
    open fun close(): Unit = definedExternally
    companion object {
        fun getInstance(elem: Element): Materialbox = definedExternally
        fun init(els: Element, options: Any? = definedExternally /* null */): Materialbox = definedExternally
        fun init(els: MElements, options: Any? = definedExternally /* null */): Array<Materialbox> = definedExternally
    }
}
external interface MaterialboxOptions {
    var inDuration: Number
    var outDuration: Number
    var onOpenStart: (`this`: Materialbox, el: Element) -> Unit
    var onOpenEnd: (`this`: Materialbox, el: Element) -> Unit
    var onCloseStart: (`this`: Materialbox, el: Element) -> Unit
    var onCloseEnd: (`this`: Materialbox, el: Element) -> Unit
}
