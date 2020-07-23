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

external open class Carousel : Component<CarouselOptions> {
    open var pressed: Boolean = definedExternally
    open var dragged: Number = definedExternally
    open var center: Number = definedExternally
    open fun next(n: Number? = definedExternally /* null */): Unit = definedExternally
    open fun prev(n: Number? = definedExternally /* null */): Unit = definedExternally
    open fun set(n: Number? = definedExternally /* null */): Unit = definedExternally
    companion object {
        fun getInstance(elem: Element): Carousel = definedExternally
        fun init(els: Element, options: Any? = definedExternally /* null */): Carousel = definedExternally
        fun init(els: MElements, options: Any? = definedExternally /* null */): Array<Carousel> = definedExternally
    }
}
external interface CarouselOptions {
    var duration: Number
    var dist: Number
    var shift: Number
    var padding: Number
    var numVisible: Number
    var fullWidth: Boolean
    var indicators: Boolean
    var noWrap: Boolean
    var onCycleTo: (`this`: Carousel, current: Element, dragged: Boolean) -> Unit
}
