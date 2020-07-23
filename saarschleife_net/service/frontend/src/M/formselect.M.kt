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

external open class FormSelect : Component<FormSelectOptions> {
    open var isMultiple: Boolean = definedExternally
    open var wrapper: Element = definedExternally
    open var dropdownOptions: HTMLUListElement = definedExternally
    open var input: HTMLInputElement = definedExternally
    open var dropdown: Dropdown = definedExternally
    open fun getSelectedValues(): Array<String> = definedExternally
    companion object {
        fun getInstance(elem: Element): FormSelect = definedExternally
        fun init(els: Element, options: Any? = definedExternally /* null */): FormSelect = definedExternally
        fun init(els: MElements, options: Any? = definedExternally /* null */): Array<FormSelect> = definedExternally
    }
}
external interface FormSelectOptions {
    var classes: String
    var dropdownOptions: Any?
}
