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

external open class Pushpin : Component<PushpinOptions> {
    open var originalOffset: Number = definedExternally
    companion object {
        fun getInstance(elem: Element): Pushpin = definedExternally
        fun init(els: Element, options: Any? = definedExternally /* null */): Pushpin = definedExternally
        fun init(els: MElements, options: Any? = definedExternally /* null */): Array<Pushpin> = definedExternally
    }
}
external interface PushpinOptions {
    var top: Number
    var bottom: Number
    var offset: Number
    var onPositionChange: (`this`: Pushpin, position: dynamic /* String /* "pinned" */ | String /* "pin-top" */ | String /* "pin-bottom" */ */) -> Unit
}
