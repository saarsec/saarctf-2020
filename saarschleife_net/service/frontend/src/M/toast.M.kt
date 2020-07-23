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

external open class Toast : ComponentBase<ToastOptions> {
    open var panning: Boolean = definedExternally
    open var timeRemaining: Number = definedExternally
    open fun dismiss(): Unit = definedExternally
    companion object {
        fun getInstance(elem: Element): Toast = definedExternally
        fun dismissAll(): Unit = definedExternally
    }
}
external interface ToastOptions {
    var html: String
    var displayLength: Number
    var inDuration: Number
    var outDuration: Number
    var classes: String
    var completeCallback: () -> Unit
    var activationPercent: Number
}
external fun toast(options: Any?): Toast = definedExternally
