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

external open class Autocomplete : Component<AutocompleteOptions> {
    open fun selectOption(el: Element): Unit = definedExternally
    open fun updateData(data: AutocompleteData): Unit = definedExternally
    open var isOpen: Boolean = definedExternally
    open var count: Number = definedExternally
    open var activeIndex: Number = definedExternally
    companion object {
        fun getInstance(elem: Element): Autocomplete = definedExternally
        fun init(els: Element, options: Any? = definedExternally /* null */): Autocomplete = definedExternally
        fun init(els: MElements, options: Any? = definedExternally /* null */): Array<Autocomplete> = definedExternally
    }
}
external interface AutocompleteData {
    //@nativeGetter
    operator fun get(key: String): String?
    //@nativeSetter
    operator fun set(key: String, value: String?)
}
external interface AutocompleteOptions {
    var data: AutocompleteData
    var limit: Number
    var onAutocomplete: (`this`: Autocomplete, text: String) -> Unit
    var minLength: Number
    var sortFunction: (a: String, b: String, inputText: String) -> Number
}
