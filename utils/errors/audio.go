package errors

var (
	ErrAudioInit        = New("audio: failed to initialize portaudio")
	ErrAudioGetDevices  = New("audio: failed to get devices")
	ErrAudioNoDevices   = New("audio: no audio devices available")
	ErrAudioOpenStream  = New("audio: failed to open stream")
	ErrAudioStartStream = New("audio: failed to start stream")
	ErrAudioTerminate   = New("audio: failed to terminate portaudio")
)
