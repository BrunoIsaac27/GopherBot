package StreamAudio

/*
 The package named StreamAudio contains the functions needed to send
 audio bytes for discord. The operation is simple but extensive.
 Initially, we define in the function that will read the audio file,
 the variables that the project will use: buffer size, buffer,
 variable to store errors, bytes read in the operation, and
 finally the variable to determine the end of the stream.
*/
import (
	"io"
	"os"
)

func ReadFile(channel chan []byte) {
	var bufferSize = 4028
	var buffer = make([]byte, bufferSize)
	var err error
	var file *os.File
	var readBytes int
	var streamDone bool
	file, err = os.Open("./X2Download.com-Joji-Yeah-Right-_LEGENDADO_TRADUCAO_.ogg")
	defer file.Close()
	if err != nil {
		panic("unable to read the file")
	}
	for !streamDone {
		readBytes, err = file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				println("stream done.")
				streamDone = true
				close(channel)
				continue
			}
		}
		channel <- buffer[:readBytes]
		println(buffer[:readBytes])
	}
}
