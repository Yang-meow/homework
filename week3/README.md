##week3

errgroup: https://pkg.go.dev/golang.org/x/sync/errgroup   
使用WithContext的errgroup，可以在一个goroutine退出之后，让其他的goroutine也及时退出，达到整体优雅退出的效果。