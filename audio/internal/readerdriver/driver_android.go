// Copyright 2021 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package readerdriver

// TODO: Use AAudio and OpenSL. See https://github.com/google/oboe.

/*

#include <jni.h>
#include <stdlib.h>

static jclass android_media_AudioAttributes = NULL;
static jclass android_media_AudioAttributes_Builder = NULL;
static jclass android_media_AudioFormat = NULL;
static jclass android_media_AudioFormat_Builder = NULL;
static jclass android_media_AudioManager = NULL;
static jclass android_media_AudioTrack = NULL;

static char* initAudioTrack(uintptr_t java_vm, uintptr_t jni_env,
    int sampleRate, int channelNum, int bitDepthInBytes, jobject* audioTrack, int bufferSize) {
  JavaVM* vm = (JavaVM*)java_vm;
  JNIEnv* env = (JNIEnv*)jni_env;

  jclass android_os_Build_VERSION = (*env)->FindClass(env, "android/os/Build$VERSION");
  const jint availableSDK =
      (*env)->GetStaticIntField(
          env, android_os_Build_VERSION,
          (*env)->GetStaticFieldID(env, android_os_Build_VERSION, "SDK_INT", "I"));
  (*env)->DeleteLocalRef(env, android_os_Build_VERSION);

  if (!android_media_AudioFormat) {
    jclass local = (*env)->FindClass(env, "android/media/AudioFormat");
    android_media_AudioFormat = (*env)->NewGlobalRef(env, local);
    (*env)->DeleteLocalRef(env, local);
  }

  if (!android_media_AudioManager) {
    jclass local = (*env)->FindClass(env, "android/media/AudioManager");
    android_media_AudioManager = (*env)->NewGlobalRef(env, local);
    (*env)->DeleteLocalRef(env, local);
  }

  if (!android_media_AudioTrack) {
    jclass local = (*env)->FindClass(env, "android/media/AudioTrack");
    android_media_AudioTrack = (*env)->NewGlobalRef(env, local);
    (*env)->DeleteLocalRef(env, local);
  }

  const jint android_media_AudioManager_STREAM_MUSIC =
      (*env)->GetStaticIntField(
          env, android_media_AudioManager,
          (*env)->GetStaticFieldID(env, android_media_AudioManager, "STREAM_MUSIC", "I"));
  const jint android_media_AudioTrack_MODE_STREAM =
      (*env)->GetStaticIntField(
          env, android_media_AudioTrack,
          (*env)->GetStaticFieldID(env, android_media_AudioTrack, "MODE_STREAM", "I"));
  const jint android_media_AudioFormat_CHANNEL_OUT_MONO =
      (*env)->GetStaticIntField(
          env, android_media_AudioFormat,
          (*env)->GetStaticFieldID(env, android_media_AudioFormat, "CHANNEL_OUT_MONO", "I"));
  const jint android_media_AudioFormat_CHANNEL_OUT_STEREO =
      (*env)->GetStaticIntField(
          env, android_media_AudioFormat,
          (*env)->GetStaticFieldID(env, android_media_AudioFormat, "CHANNEL_OUT_STEREO", "I"));
  const jint android_media_AudioFormat_ENCODING_PCM_8BIT =
      (*env)->GetStaticIntField(
          env, android_media_AudioFormat,
          (*env)->GetStaticFieldID(env, android_media_AudioFormat, "ENCODING_PCM_8BIT", "I"));
  const jint android_media_AudioFormat_ENCODING_PCM_16BIT =
      (*env)->GetStaticIntField(
          env, android_media_AudioFormat,
          (*env)->GetStaticFieldID(env, android_media_AudioFormat, "ENCODING_PCM_16BIT", "I"));

  jint channel = android_media_AudioFormat_CHANNEL_OUT_MONO;
  switch (channelNum) {
  case 1:
    channel = android_media_AudioFormat_CHANNEL_OUT_MONO;
    break;
  case 2:
    channel = android_media_AudioFormat_CHANNEL_OUT_STEREO;
    break;
  default:
    return "invalid channel";
  }

  jint encoding = android_media_AudioFormat_ENCODING_PCM_8BIT;
  switch (bitDepthInBytes) {
  case 1:
    encoding = android_media_AudioFormat_ENCODING_PCM_8BIT;
    break;
  case 2:
    encoding = android_media_AudioFormat_ENCODING_PCM_16BIT;
    break;
  default:
    return "invalid bitDepthInBytes";
  }

  // If the available Android SDK is at least 24 (7.0 Nougat), the FLAG_LOW_LATENCY is available.
  // This requires a different constructor.
  if (availableSDK >= 24) {
    if (!android_media_AudioAttributes_Builder) {
      jclass local = (*env)->FindClass(env, "android/media/AudioAttributes$Builder");
      android_media_AudioAttributes_Builder = (*env)->NewGlobalRef(env, local);
      (*env)->DeleteLocalRef(env, local);
    }

    if (!android_media_AudioFormat_Builder) {
      jclass local = (*env)->FindClass(env, "android/media/AudioFormat$Builder");
      android_media_AudioFormat_Builder = (*env)->NewGlobalRef(env, local);
      (*env)->DeleteLocalRef(env, local);
    }

    if (!android_media_AudioAttributes) {
      jclass local = (*env)->FindClass(env, "android/media/AudioAttributes");
      android_media_AudioAttributes = (*env)->NewGlobalRef(env, local);
      (*env)->DeleteLocalRef(env, local);
    }

    jint android_media_AudioAttributes_USAGE_UNKNOWN =
        (*env)->GetStaticIntField(
            env, android_media_AudioAttributes,
            (*env)->GetStaticFieldID(env, android_media_AudioAttributes, "USAGE_UNKNOWN", "I"));
    jint android_media_AudioAttributes_CONTENT_TYPE_UNKNOWN =
        (*env)->GetStaticIntField(
            env, android_media_AudioAttributes,
            (*env)->GetStaticFieldID(env, android_media_AudioAttributes, "CONTENT_TYPE_UNKNOWN", "I"));
    jint android_media_AudioAttributes_FLAG_LOW_LATENCY =
        (*env)->GetStaticIntField(
            env, android_media_AudioAttributes,
            (*env)->GetStaticFieldID(env, android_media_AudioAttributes, "FLAG_LOW_LATENCY", "I"));

    const jobject aattrBld =
        (*env)->NewObject(
            env, android_media_AudioAttributes_Builder,
            (*env)->GetMethodID(env, android_media_AudioAttributes_Builder, "<init>", "()V"));

    (*env)->CallObjectMethod(
        env, aattrBld,
        (*env)->GetMethodID(env, android_media_AudioAttributes_Builder, "setUsage", "(I)Landroid/media/AudioAttributes$Builder;"),
        android_media_AudioAttributes_USAGE_UNKNOWN);
    (*env)->CallObjectMethod(
        env, aattrBld,
        (*env)->GetMethodID(env, android_media_AudioAttributes_Builder, "setContentType", "(I)Landroid/media/AudioAttributes$Builder;"),
        android_media_AudioAttributes_CONTENT_TYPE_UNKNOWN);
    (*env)->CallObjectMethod(
        env, aattrBld,
        (*env)->GetMethodID(env, android_media_AudioAttributes_Builder, "setFlags", "(I)Landroid/media/AudioAttributes$Builder;"),
        android_media_AudioAttributes_FLAG_LOW_LATENCY);
    const jobject aattr =
        (*env)->CallObjectMethod(
            env, aattrBld,
            (*env)->GetMethodID(env, android_media_AudioAttributes_Builder, "build", "()Landroid/media/AudioAttributes;"));
    (*env)->DeleteLocalRef(env, aattrBld);

    const jobject afmtBld =
        (*env)->NewObject(
            env, android_media_AudioFormat_Builder,
            (*env)->GetMethodID(env, android_media_AudioFormat_Builder, "<init>", "()V"));
    (*env)->CallObjectMethod(
        env, afmtBld,
        (*env)->GetMethodID(env, android_media_AudioFormat_Builder, "setSampleRate", "(I)Landroid/media/AudioFormat$Builder;"),
        sampleRate);
    (*env)->CallObjectMethod(
        env, afmtBld,
        (*env)->GetMethodID(env, android_media_AudioFormat_Builder, "setEncoding", "(I)Landroid/media/AudioFormat$Builder;"),
        encoding);
    (*env)->CallObjectMethod(
        env, afmtBld,
        (*env)->GetMethodID(env, android_media_AudioFormat_Builder, "setChannelMask", "(I)Landroid/media/AudioFormat$Builder;"),
        channel);
    const jobject afmt =
        (*env)->CallObjectMethod(
            env, afmtBld,
            (*env)->GetMethodID(env, android_media_AudioFormat_Builder, "build", "()Landroid/media/AudioFormat;"));
    (*env)->DeleteLocalRef(env, afmtBld);

    const jobject tmpAudioTrack =
        (*env)->NewObject(
            env, android_media_AudioTrack,
            (*env)->GetMethodID(env, android_media_AudioTrack, "<init>", "(Landroid/media/AudioAttributes;Landroid/media/AudioFormat;III)V"),
            aattr, afmt, bufferSize, android_media_AudioTrack_MODE_STREAM, 0);
    *audioTrack = (*env)->NewGlobalRef(env, tmpAudioTrack);
    (*env)->DeleteLocalRef(env, tmpAudioTrack);
    (*env)->DeleteLocalRef(env, aattr);
    (*env)->DeleteLocalRef(env, afmt);
  } else {
    const jobject tmpAudioTrack =
        (*env)->NewObject(
            env, android_media_AudioTrack,
            (*env)->GetMethodID(env, android_media_AudioTrack, "<init>", "(IIIIII)V"),
            android_media_AudioManager_STREAM_MUSIC,
            sampleRate, channel, encoding, bufferSize,
            android_media_AudioTrack_MODE_STREAM);
    *audioTrack = (*env)->NewGlobalRef(env, tmpAudioTrack);
    (*env)->DeleteLocalRef(env, tmpAudioTrack);
  }

  // TODO: Use an independent function?
  (*env)->CallVoidMethod(
      env, *audioTrack,
      (*env)->GetMethodID(env, android_media_AudioTrack, "play", "()V"));

  return NULL;
}

*/
import "C"

import (
	"io"

	"golang.org/x/mobile/app"
)

func IsAvailable() bool {
	return true
}

type contextImpl struct {
	sampleRate      int
	channelNum      int
	bitDepthInBytes int
}

func NewContext(sampleRate int, channelNum int, bitDepthInBytes int) (Context, error) {
	return &contextImpl{
		sampleRate:      sampleRate,
		channelNum:      channelNum,
		bitDepthInBytes: bitDepthInBytes,
	}, nil
}

func (c *contextImpl) NewPlayer(src io.Reader) Player {
	return &playerImpl{
		context: c,
		src:     src,
	}
}

func (c *contextImpl) Close() error {
	// TODO: Implement this
	return nil
}

type playerImpl struct {
	context *contextImpl
	src     io.Reader
	err     error
	state   playerState
	eof     bool

	audioTrack C.jobject
}

func (p *playerImpl) Pause() {
	if p.state != playerStatePlay {
		return
	}
	if p.audioTrack == nil {
		return
	}

	// TODO: call pause
}

func (p *playerImpl) Play() {
	if p.state != playerStatePaused {
		return
	}
	if p.eof {
		return
	}
	if p.audioTrack == nil {
		if err := app.RunOnJVM(func(vm, env, ctx uintptr) error {
			audioTrack := C.jobject(0)
			if msg := C.initAudioTrack(C.uintptr_t(vm), C.uintptr_t(env),
				C.int(sampleRate), C.int(channelNum), C.int(bitDepthInBytes),
				&audioTrack, C.int(p.context.maxBufferSize())); msg != nil {
				return errors.New("readerdriver: initAutioTrack failed: " + C.GoString(msg))
			}
			p.audioTrack = audioTrack
			return nil
		}); err != nil {
			p.err = err
			// TODO: mutex
			p.Reset()
			p.state = playerState
			return
		}
	}
	// TODO: call play
	go p.loop()
}

func (p *playerImpl) loop() {
	buf := make([]byte, p.context.oneBufferSize())
	for {
		n, err := p.src.Read(buf)
		if err != nil && err != io.EOF {
			p.err = err
			// TODO: Close?
			return
		}
		// TODO: write buf[:n]
		// TODO: lock
		if err == io.EOF {
			p.eof = true
			p.Pause() // ?
			return
		}
	}
}

func (p *playerImpl) IsPlaying() bool {
	return p.state == playerStatePlay
}

func (p *playerImpl) Reset() {
	if p.audioTrack == nil {
		return
	}

	// TODO: Call flush and stop.
	p.state = playerStatePaused
}

func (p *playerImpl) Volume() float64 {
	return 0
}

func (p *playerImpl) SetVolume(volume float64) {
}

func (p *playerImpl) UnplayedBufferSize() int64 {
	if p.audioTrack == nil {
		return 0
	}
	return 0
}

func (p *playerImpl) Err() error {
	return p.err
}

func (p *playerImpl) Close() error {
	p.state = playerStateClosed
	// TODO: Call release
	return p.err
}
