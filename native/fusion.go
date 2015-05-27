package native

/*
#cgo linux pkg-config: opencv
#cgo darwin pkg-config: opencv
#cgo windows LDFLAGS: -lopencv_core -lopencv_imgproc -lopencv_photo -lopencv_imgcodecs
*/
import "C"
