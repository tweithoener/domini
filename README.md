# domini -- Mini DOM Interface For Golang/WASM

This is an unstable alpha version. No guarantees whatsoever.

A minimalistiv package to access the web pages DOM from a go web app when using WebAssembly (wasm).
With this package the following cumbersome, verbose line of code

```golang
  js.Global().Call("getElementById", "msg").Get("style").Call("setProperty", "display", "block", "") 
```

becomes a lot more tidy:

```golang
  domini.GetWindow().GetElementById("msg").Style().SetProperty("display", "block", "")
```

As you can see, this is a little shorter, much easier to read, and hell of a lot faster to type as most of the typing can be done by the
source code completion of your editor. Also less error-prone as the compiler can check method names where there were string
arguments before.

Install this package using go get or require it as a as module.

