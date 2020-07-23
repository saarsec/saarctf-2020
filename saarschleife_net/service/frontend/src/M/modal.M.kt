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

external open class Modal : Component<ModalOptions>, Openable {
    override fun open(): Unit = definedExternally
    override fun close(): Unit = definedExternally
    override var isOpen: Boolean = definedExternally
    open var id: String = definedExternally
    companion object {
        fun getInstance(elem: Element): Modal = definedExternally
        fun init(els: Element, options: Any? = definedExternally /* null */): Modal = definedExternally
        fun init(els: MElements, options: Any? = definedExternally /* null */): Array<Modal> = definedExternally
    }
}
external interface ModalOptions {
    var opacity: Number
    var inDuration: Number
    var outDuration: Number
    var preventScrolling: Boolean
    var onOpenStart: (`this`: Modal, el: Element) -> Unit
    var onOpenEnd: (`this`: Modal, el: Element) -> Unit
    var onCloseStart: (`this`: Modal, el: Element) -> Unit
    var onCloseEnd: (`this`: Modal, el: Element) -> Unit
    var dismissible: Boolean
    var startingTop: String
    var endingTop: String
}
