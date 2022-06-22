module github.com/goplusjs/repl

go 1.16

require (
	github.com/goplus/igop v0.7.5-0.20220622134234-57b7b9638c95
	github.com/goplusjs/gopherjs v1.2.5
)

replace (
	github.com/goplus/gop => github.com/goplus/gop v1.1.2
	github.com/goplus/reflectx => github.com/goplusjs/reflectx v0.5.6
	github.com/petermattis/goid => github.com/visualfc/goid v0.1.0
	github.com/visualfc/xtype => github.com/visualfc/xtype_js v0.1.0
	golang.org/x/mod => golang.org/x/mod v0.5.1
	golang.org/x/tools => golang.org/x/tools v0.1.10
)
