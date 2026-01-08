package audio

import (
	"sync"

	configTypes "github.com/chloyka/gorig/internal/config/types"
	"github.com/chloyka/gorig/internal/effects"
	"github.com/chloyka/gorig/internal/logger"
	"github.com/chloyka/gorig/internal/logger/keys"
	"github.com/chloyka/gorig/internal/onset"
	"github.com/chloyka/gorig/internal/rhythm"
	errs "github.com/chloyka/gorig/utils/errors"
	"github.com/gordonklaus/portaudio"
)

type Engine struct {
	mu            sync.Mutex
	logger        *logger.Logger
	cfg           *configTypes.AudioConfig
	stateConfig   *configTypes.StateConfig
	stream        *portaudio.Stream
	chain         *effects.Chain
	onsetDetector *onset.Detector
	rhythmEngine  *rhythm.Engine

	inputDevices  []*portaudio.DeviceInfo
	outputDevices []*portaudio.DeviceInfo
	inputIndex    int
	outputIndex   int
}

func newEngine(logger *logger.Logger, chain *effects.Chain, onsetDetector *onset.Detector, rhythmEngine *rhythm.Engine, cfg *configTypes.AudioConfig, stateConfig *configTypes.StateConfig) (*Engine, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, errs.Wrap(errs.ErrAudioInit, err)
	}

	e := &Engine{
		logger:        logger,
		cfg:           cfg,
		stateConfig:   stateConfig,
		chain:         chain,
		onsetDetector: onsetDetector,
		rhythmEngine:  rhythmEngine,
	}

	if err := e.loadDevices(); err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Engine) loadDevices() error {
	devices, err := portaudio.Devices()
	if err != nil {
		return errs.Wrap(errs.ErrAudioGetDevices, err)
	}

	defaultInput, _ := portaudio.DefaultInputDevice()
	defaultOutput, _ := portaudio.DefaultOutputDevice()

	for _, d := range devices {
		if d.MaxInputChannels > 0 {
			e.inputDevices = append(e.inputDevices, d)
			if defaultInput != nil && d.Name == defaultInput.Name {
				e.inputIndex = len(e.inputDevices) - 1
			}
		}
		if d.MaxOutputChannels > 0 {
			e.outputDevices = append(e.outputDevices, d)
			if defaultOutput != nil && d.Name == defaultOutput.Name {
				e.outputIndex = len(e.outputDevices) - 1
			}
		}
	}

	e.logger.Info("found audio devices",
		keys.DeviceInputCount(len(e.inputDevices)),
		keys.DeviceOutputCount(len(e.outputDevices)),
	)

	e.restoreSavedDevices()

	return nil
}

func (e *Engine) restoreSavedDevices() {

	if e.stateConfig.InputDevice != "" {
		idx := e.findDeviceByName(e.inputDevices, e.stateConfig.InputDevice)
		if idx >= 0 {
			e.inputIndex = idx
			e.logger.Info("restored saved input device", keys.DeviceName(e.stateConfig.InputDevice))
		} else {
			e.logger.Warn("saved input device not found, using default",
				keys.DeviceSavedName(e.stateConfig.InputDevice),
				keys.DeviceUsingName(e.inputDevices[e.inputIndex].Name),
			)
		}
	}

	if e.stateConfig.OutputDevice != "" {
		idx := e.findDeviceByName(e.outputDevices, e.stateConfig.OutputDevice)
		if idx >= 0 {
			e.outputIndex = idx
			e.logger.Info("restored saved output device", keys.DeviceName(e.stateConfig.OutputDevice))
		} else {
			e.logger.Warn("saved output device not found, using default",
				keys.DeviceSavedName(e.stateConfig.OutputDevice),
				keys.DeviceUsingName(e.outputDevices[e.outputIndex].Name),
			)
		}
	}
}

func (e *Engine) findDeviceByName(devices []*portaudio.DeviceInfo, name string) int {
	for i, d := range devices {
		if d.Name == name {
			return i
		}
	}
	return -1
}

func (e *Engine) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.startStream()
}

func (e *Engine) startStream() error {
	if len(e.inputDevices) == 0 || len(e.outputDevices) == 0 {
		return errs.ErrAudioNoDevices
	}

	inputDev := e.inputDevices[e.inputIndex]
	outputDev := e.outputDevices[e.outputIndex]

	e.logger.Info("using devices",
		keys.DeviceInputName(inputDev.Name),
		keys.DeviceOutputName(outputDev.Name),
	)

	streamParams := portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   inputDev,
			Channels: e.cfg.NumChannels,
			Latency:  e.cfg.TargetLatency,
		},
		Output: portaudio.StreamDeviceParameters{
			Device:   outputDev,
			Channels: e.cfg.NumChannels,
			Latency:  e.cfg.TargetLatency,
		},
		SampleRate:      float64(e.cfg.SampleRate),
		FramesPerBuffer: e.cfg.FramesPerBuffer,
	}

	onsetDet := e.onsetDetector
	rhythmEng := e.rhythmEngine

	stream, err := portaudio.OpenStream(streamParams, func(in, out []float32) {
		copy(out, in)

		if onsetDet != nil {
			onsetDet.Process(in)
		}

		if rhythmEng != nil {
			if q := rhythmEng.ProcessBuffer(len(in)); q != nil {
				effects.SetCurrentOnset(true, q.OriginalEvent.Energy, q.BeatPosition, q.SlotIndex)
			} else {
				effects.ClearCurrentOnset()
			}
		}

		e.chain.Process(out)
	})
	if err != nil {
		return errs.Wrap(errs.ErrAudioOpenStream, err)
	}

	e.stream = stream

	if err = stream.Start(); err != nil {
		return errs.Wrap(errs.ErrAudioStartStream, err)
	}

	info := stream.Info()

	log := e.logger.With(
		keys.AudioSampleRate(e.cfg.SampleRate),
		keys.AudioFramesPerBuffer(e.cfg.FramesPerBuffer),
	)

	if info != nil {
		log = log.With(
			keys.AudioInputLatencyMs(info.InputLatency.Seconds()*1000),
			keys.AudioOutputLatencyMs(info.OutputLatency.Seconds()*1000),
		)
	}

	log.Info("audio engine started")

	return nil
}

func (e *Engine) stopStream() {
	if e.stream != nil {
		if err := e.stream.Stop(); err != nil {
			e.logger.Warn("failed to stop stream", keys.Error(err))
		}

		if err := e.stream.Close(); err != nil {
			e.logger.Warn("failed to close stream", keys.Error(err))
		}

		e.stream = nil
	}
}

func (e *Engine) NextInputDevice() string {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.inputDevices) == 0 {
		return ""
	}

	e.stopStream()

	e.inputIndex = (e.inputIndex + 1) % len(e.inputDevices)
	name := e.inputDevices[e.inputIndex].Name

	e.logger.Info("switched input device", keys.DeviceName(name))

	if err := e.startStream(); err != nil {
		e.logger.Error("failed to restart stream", keys.Error(err))
	}

	e.stateConfig.SetInputDevice(name)

	return name
}

func (e *Engine) NextOutputDevice() string {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.outputDevices) == 0 {
		return ""
	}

	e.stopStream()

	e.outputIndex = (e.outputIndex + 1) % len(e.outputDevices)
	name := e.outputDevices[e.outputIndex].Name

	e.logger.Info("switched output device", keys.DeviceName(name))

	if err := e.startStream(); err != nil {
		e.logger.Error("failed to restart stream", keys.Error(err))
	}

	e.stateConfig.SetOutputDevice(name)

	return name
}

func (e *Engine) CurrentInputDevice() string {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.inputDevices) == 0 {
		return "none"
	}
	return e.inputDevices[e.inputIndex].Name
}

func (e *Engine) CurrentOutputDevice() string {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.outputDevices) == 0 {
		return "none"
	}
	return e.outputDevices[e.outputIndex].Name
}

func (e *Engine) OnsetDetector() *onset.Detector {
	return e.onsetDetector
}

func (e *Engine) RhythmEngine() *rhythm.Engine {
	return e.rhythmEngine
}

func (e *Engine) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.stopStream()

	if e.onsetDetector != nil {
		e.onsetDetector.Close()
	}

	if e.rhythmEngine != nil {
		e.rhythmEngine.Close()
	}

	if err := portaudio.Terminate(); err != nil {
		return errs.Wrap(errs.ErrAudioTerminate, err)
	}

	e.logger.Info("audio engine stopped")
	return nil
}
