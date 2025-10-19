# Raspi agent (poor man's Alexa)
**Raspi-Agent** is a modular, streaming voice assistant built for Raspberry Pi.  
It listens for a wake word, records your voice, processes the audio using OpenAI or a backend service,  
and responds with natural speech — all in real time.



![](https://github.com/ownerofglory/raspi-agent/actions/workflows/build.yaml/badge.svg)

*In development ...*

<img src="./docs/assets/Raspi-agent.png" width="480px" />

---

##  Features

- **Wake Word Detection** — powered by [Porcupine](https://picovoice.ai/platform/porcupine/)
- **Streaming Audio Playback** — real-time PCM or MP3 output via PortAudio
- **Natural Conversation** — integrates with OpenAI (STT, LLM, TTS)
- **Dual Architecture** — choose between:
    - **Onboard mode** — runs all AI calls directly from the Pi
    - **Offboard mode** — sends recordings to a backend for processing


## Architecture Overview
> Wake Word → Recorder → Voice Assistant → Player → User

Two orchestrator implementations:

| Mode | Description | Example                         |
|------|--------------|---------------------------------|
| **Onboard** | Runs everything locally via OpenAI APIs (STT, LLM, TTS) | `cmd/raspi-agent-onboard-local`       |
| **Offboard** | Streams recorded audio to a backend that processes it | `cmd/raspi-agent-onboard` |

---
