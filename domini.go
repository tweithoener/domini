// +build js,wasm

// Package domini is minimalistiv package to access the web pages DOM from a go web app when using WebAssembly (wasm)
package domini

import "syscall/js"

type dom struct {
	js.Value
}

// GetWindow returns the root handle into JS
func GetWindow() Window {
	return &dom{js.Global()}
}

// Underlying is implemented by all DOM elements. It is used to
// access the underlying js.Value.
type Underlyer interface {
	Underlying() js.Value
	IsUndefined() bool
	IsNull() bool
}

// Queryer implements document querying functions
type Queryer interface {
	QuerySelector(query string) HTMLElement
	QuerySelectorAll(query string) []HTMLElement
}

// Class represents the class list of an HTML element
type Class interface {
	Underlyer
	Add(classes ...string)
	Remove(classes ...string)
	Contains(class string) bool
}

// Style represents the CSS style declaration of an HTML element
type Style interface {
	Underlyer
	SetProperty(prop, value, prio string)
	Property(prop string) string
	RemoveProperty(prop string)
}

// HTMLElement represents an HTML element in a DOM
type HTMLElement interface {
	Underlyer
	EventTarget
	EventSource
	Queryer
	ID() string
	SetID(id string)
	TagName() string
	Class() Class
	Style() Style
	SetInnerHTML(html string)
	InnerHTML() string
	SetAttribute(attr, value string)
	AppendChild(el HTMLElement)
	RemoveChild(el HTMLElement)
	ChildElements() []HTMLElement
	ParentElement() HTMLElement
	SetData(name, value string)
	Data(name string) js.Value
}

// Event represents a DOM event
type Event interface {
	Underlyer
	PreventDefault()
	StopPropagation()
}

// Window represents the window
type Window interface {
	Underlyer
	EventTarget
	EventSource
	Document() Document
}

// Document represents the DOM root element
type Document interface {
	Underlyer
	EventTarget
	EventSource
	Queryer
	GetElementByID(id string) HTMLElement
	Body() HTMLElement
	DocumentElement() HTMLElement
	CreateElement(tag string) HTMLElement
	CreateElementNS(ns, tag string) HTMLElement
}

func (d *dom) IsNull() bool {
	return d == nil || d.Underlying().IsNull()
}

func (d *dom) IsUndefined() bool {
	return d.Underlying().IsUndefined()
}

func (d *dom) Document() Document {
	return &dom{d.Get("document")}
}

func (d *dom) QuerySelector(query string) HTMLElement {
	return &dom{d.Call("querySelector", query)}
}

func (d *dom) QuerySelectorAll(query string) []HTMLElement {
	raw := d.Call("querySelectorAll", query)
	list := make([]HTMLElement, 0, raw.Length())
	for i := 0; i < raw.Length(); i++ {
		list = append(list, &dom{raw.Index(i)})
	}
	return list
}

func (d *dom) GetElementByID(id string) HTMLElement {
	return &dom{d.Call("getElementById", id)}
}

func (d *dom) Body() HTMLElement {
	return &dom{d.Get("body")}
}

func (d *dom) DocumentElement() HTMLElement {
	return &dom{d.Get("documentElement")}
}

func (d *dom) CreateElement(tag string) HTMLElement {
	return &dom{d.Call("createElement", tag)}
}

func (d *dom) CreateElementNS(ns, tag string) HTMLElement {
	return &dom{d.Call("createElementNS", ns, tag)}
}

func (d *dom) ID() string {
	return d.Get("id").String()
}

func (d *dom) SetID(id string) {
	d.Set("id", id)
}

func (d *dom) TagName() string {
	return d.Get("tagName").String()
}

func (d *dom) AppendChild(el HTMLElement) {
	d.Call("appendChild", el.Underlying())
}

func (d *dom) RemoveChild(el HTMLElement) {
	d.Call("removeChild", el.Underlying())
}

func (d *dom) ChildElements() []HTMLElement {
	raw := d.Get("children")
	list := make([]HTMLElement, 0, raw.Length())
	for i := 0; i < raw.Length(); i++ {
		list = append(list, &dom{raw.Index(i)})
	}
	return list
}

func (d *dom) ParentElement() HTMLElement {
	raw := d.Get("parentElement")
	return &dom{raw}
}

func (d *dom) Class() Class {
	return &dom{d.Get("classList")}
}

func (d *dom) Add(classes ...string) {
	for _, cl := range classes {
		d.Call("add", cl)
	}
}

func (d *dom) Remove(classes ...string) {
	for _, cl := range classes {
		d.Call("remove", cl)
	}
}

func (d *dom) Contains(class string) bool {
	return d.Call("contains", class).Bool()
}

func (d *dom) Style() Style {
	return &dom{d.Get("style")}
}

func (d *dom) SetInnerHTML(html string) {
	d.Set("innerHTML", html)
}

func (d *dom) InnerHTML() string {
	return d.Get("innerHTML").String()
}

func (d *dom) SetAttribute(attr, value string) {
	d.Call("setAttribute", attr, value)
}

func (d *dom) SetData(name, value string) {
	d.Get("dataset").Set(name, value)
}
func (d *dom) Data(name string) js.Value {
	return d.Get("dataset").Get(name)
}

func (d *dom) SetProperty(prop, value, prio string) {
	d.Call("setProperty", prop, value, prio)
}

func (d *dom) Property(prop string) string {
	return d.Get(prop).String()
}

func (d *dom) RemoveProperty(prop string) {
	d.Call("removeProperty", prop)
}

func (d *dom) PreventDefault() {
	d.Call("preventDefault")
}

func (d *dom) StopPropagation() {
	d.Call("stopPropagation")
}

func (d *dom) DispatchEvent(ev Event) {
	d.Call("dispatchEvent", ev.Underlying())
}

func (d *dom) Underlying() js.Value {
	return d.Value
}

func Null() js.Value {
	return js.Null()
}

func Undefined() js.Value {
	return js.Undefined()
}

// EventSource is implemented by all elements that dispatch events
type EventSource interface {
	DispatchEvent(ev Event)
}

// EventTarget is implemented by all elements that receive events
type EventTarget interface {
	AddEventListener(event string, useCapture bool, function func(Event)) js.Func
	RemoveEventListener(event string, useCapture bool, cb js.Func)
}

func (d *dom) AddEventListener(event string, useCapture bool, function func(Event)) js.Func {
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		function(&dom{args[0]})
		return nil
	})
	d.Call("addEventListener", event, cb, useCapture)
	return cb
}

func (d *dom) RemoveEventListener(event string, useCapture bool, cb js.Func) {
	cb.Release()
}

// NewEvent returns a new event of type t
func NewEvent(t string) Event {
	return &dom{js.Global().Get("Event").New(t)}
}
