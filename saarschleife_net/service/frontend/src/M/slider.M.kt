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

external open class Slider : Component<SliderOptions> {
    override var el: Element = definedExternally
    override var options: SliderOptions = definedExternally
    open var activeIndex: Number = definedExternally
    open fun pause(): Unit = definedExternally
    open fun start(): Unit = definedExternally
    open fun next(): Unit = definedExternally
    open fun prev(): Unit = definedExternally
    companion object {
        fun getInstance(elem: Element): Slider = definedExternally
        fun init(els: Element, options: Any? = definedExternally /* null */): Slider = definedExternally
        fun init(els: MElements, options: Any? = definedExternally /* null */): Array<Slider> = definedExternally
    }
}
external interface SliderOptions {
    var indicators: Boolean
    var height: Number
    var duration: Number
    var interval: Number
}
