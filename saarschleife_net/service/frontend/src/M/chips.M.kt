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

external open class Chips : Component<ChipsOptions> {
    open var chipsData: Array<ChipData> = definedExternally
    open var hasAutocomplete: Boolean = definedExternally
    open var autocomplete: Autocomplete = definedExternally
    open fun addChip(chip: ChipData): Unit = definedExternally
    open fun deleteChip(n: Number? = definedExternally /* null */): Unit = definedExternally
    open fun selectChip(n: Number): Unit = definedExternally
    companion object {
        fun getInstance(elem: Element): Chips = definedExternally
        fun init(els: Element, options: Any? = definedExternally /* null */): Chips = definedExternally
        fun init(els: MElements, options: Any? = definedExternally /* null */): Array<Chips> = definedExternally
    }
}
external interface ChipData {
    var tag: String
    var img: String? get() = definedExternally; set(value) = definedExternally
}
external interface ChipsOptions {
    var data: Array<ChipData>
    var placeholder: String
    var secondaryPlaceholder: String
    var autocompleteOptions: Any?
    var limit: Number
    var onChipAdd: (`this`: Chips, element: Element, chip: Element) -> Unit
    var onChipSelect: (`this`: Chips, element: Element, chip: Element) -> Unit
    var onChipDelete: (`this`: Chips, element: Element, chip: Element) -> Unit
}
